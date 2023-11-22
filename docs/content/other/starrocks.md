# starrocks OLAP

其实 tangseng 现阶段来说mysql这种OLTP已经足够了，但是我们还是要学一下OLAP，因为当mysql到了一定量级之后，比亿级数据的时候，查询就会变的很慢，这时候索引也就不活这张表，或者这个库，到这个时候，早就分库分表了。

当然除了分库分表，我们也可以通过dts(data transmission service),比如flink cdc，canal等等...将增量数据同步到OLAP集群。从OLAP读取，OLAP读取非常的快，单机 starrocks 的单表支持亿级数据毫秒级查询。

当然像starrocks这种OLAP只是查询快而已，插入或者更新相对OLTP来说比较慢，但其实在tangseng中我们只关心读取而已，所有的建立索引的写操作都是脚本启动。

另外starrocks是完美兼容mysql的，其实很多最近出现的数据库都会选择兼容mysql，比如tidb。

starrocks文档：https://docs.starrocks.io/en-us/latest/introduction/StarRocks_intro

## sr简单介绍

## sr的CURD操作

### sr 启动!

docker all in one启动starrocks

```shell
sudo docker run -p 9030:9030 -p 8030:8030 -p 8040:8040 -itd starrocks.docker.scarf.sh/starrocks/allin1-ubuntu
```

讲一下几个端口的用途

- 8030:网页端连接，localhost:8083,用户名:root,密码为空
- 8040:BE的连接节点
- 9030:客户端连接starrocks，类似mysql的3306

连接starrocks，其实和mysql一样的连接，只是接口改成了9030

### 读操作

由于是完美兼容mysql的，所以我们可以用gorm来连接starrocks，

### 写操作

starrocks没有提供go语言的sdk，所以我们只能根据http请求来进行写操作，在sr中，我们尽量采用批量写入的方式。tangseng中是使用`双缓冲数据通道`来处理的

定义上传对象

```go
type DirectUpload struct {
    ctx     context.Context
    doneCtx context.Context

    data   []*types.Data2Starrocks // 数据
    upData []*types.Data2Starrocks // 上传的数据
    wLock  *sync.Mutex      // 写锁
    upLock *sync.RWMutex    // 上传加锁
    task   *types.Task

    done func()
}
```

我们可以监听上游的数据请求，来一条数据，我们就把数据放到`d.data`中。这里加锁是因为最开始的时候`d.data`其实是map类型，后面才改成数组类型，所以这个`d.wLock`其实是可以去掉的，毕竟频率加解锁也会影响性能。

```go
func (d *DirectUpload) Push(records *types.Data2Starrocks) int {
    d.wLock.Lock()
    defer d.wLock.Unlock()
    d.data = append(d.data, records)
    log.LogrusObj.Infof("direct_upload push bi_table:%s", d.task.BiTable)

    return len(d.data)
}
```

五分钟消费进行一次消费,如果脚本重启或者停止了，将会把内存中所存在的进行消费读入。

```go
func (d *DirectUpload) consume() {
    gapTime := 5 * time.Minute
    for {
        select {
        case <-time.After(gapTime):
            log.LogrusObj.Infof("direct upload")
            _, _ = d.StreamUpload()
        case <-d.doneCtx.Done():
            _, _ = d.StreamUpload()
        }
    }
}
```

具体的消费上传逻辑如下：

首先我们把数据从`d.data`转移到上传数据中`d.upData`这个数组中,接着我们将`d.data`的数据清空,再解开写锁`wLock`,便完成了转移，之后的`d.data`将继续进行消费上游请求操作。

```go
d.wLock.Lock()
if len(d.data) == 0 {
    d.upData = d.data
} else {
    d.upData = append(d.upData, d.data...)
}
d.data = make([]*types.Data2Starrocks, 0)
count = len(d.upData)
d.wLock.Unlock()
```

开始上报数据,上报数据的过程我们加上一个上报数据的锁`upLock`。

```go
d.upLock.Lock()
defer d.upLock.Unlock()

if len(d.upData) == 0 {
    log.LogrusObj.Infof("finish upload")
}
```

tangseng中的上报sr是采用构建csv的方式上传,我们构建一个csv格式的文件流进行传输

```go
// 构建csv
rowDelimiter := "@##@" // 分割线，自定义，后面构建文件流传入即可
sb := &bytes.Buffer{}
write := bufio.NewWriter(sb)
for i := 0; i < count; i++ {
    line := strings.Join([]string{
        cast.ToString(d.upData[i].DocId),
        d.upData[i].Title,
        d.upData[i].Desc,
        d.upData[i].Url,
        cast.ToString(d.upData[i].Score),
    }, ",")
    _, err = write.WriteString(line + rowDelimiter)
    if err != nil {
        log.LogrusObj.Errorf("WriteString Error")
    }
}
```

基于http连接sr客户端

```go
starrocksClient := &http.Client{
    CheckRedirect: func(req *http.Request, via []*http.Request) (err error) {
        v := via[0]
        req.Header = v.Header
        req.Body, err = v.GetBody()
        if err != nil {
            log.LogrusObj.Errorf("starrock woker")
        }
        return err
    },
    Timeout: time.Minute,
}
```

将csv文本构建在请求头中，其中
- host: starrocks的host
- db: 所上传sr的db名
- table: 上传的sr的表
- sb.Bytes(): csv文件流

```go
srConfig := config.Conf.StarRocks
hp, err := cli.SetDebug(true).R().SetContext(d.ctx).
    SetBasicAuth(srConfig.UserName, srConfig.Password).
    SetPathParams(map[string]string{
        "host":  srConfig.LoadUrl,
        "db":    srConfig.Database,
        "table": d.task.BiTable,
    }).SetBody(sb.Bytes()).SetContentLength(true).
    Put("https://{host}/api/{db}/{table}/_stream_load")
if err != nil {
    log.LogrusObj.Errorf("stream load error :%+v", err)
}
```

完成上传操作之后，我们进行重置`upData`。

```go
// 重置 updata
d.wLock.Lock()
d.upData = make([]*types.Data2Starrocks, 0)
d.wLock.Unlock()
```

这就是基于双缓存通道的上传方式上传到sr，其实我们还可以进行优化，比如将http请求改成ws长链接，减少网络连接消耗。
