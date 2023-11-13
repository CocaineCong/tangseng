# 索引平台

> index_platform 处理 & 存储倒排索引，正排索引结构.

## 项目结构

```shell
index_platform
├── analyzer            // 解析器，分词作用，与搜索平台的解析不一样
├── cmd                 // 项目启动入口
│   ├── job             // 脚本注册
│   └── kfk_register    // kafka消费注册
├── consts              // 存放该模块下的常量
├── crawl               // 爬虫(大饼🫓)
├── input_data          // 输入的数据
├── repository          // 存储介质
│   ├── db              // OLTP:mysql
│   │   └── dao
│   ├── starrock        // OLAP:starrocks 
│   │   ├── bi_dao
│   │   └── initdb.d
│   └── storage         // KV DB: boltdb
├── service             // 服务层
└── trie                // 存储trie相关，暂时无用.
```

## 服务层详解

> 主要业务逻辑在这个文件 `app/index_platform/service/index_platform.go`

### 1. 倒排索引

定义倒排索引结构

```go
invertedIndex := cmap.New[*roaring.Bitmap]() // 倒排索引
```

考虑到支持并发，所以这里用的是第三方的map结构，`github.com/orcaman/concurrent-map/v2`,相比较于官方sync包下的map结构更适合，适合**并发读写。**

存储索引id使用的数据结构是`roaring bitmap`, `github.com/RoaringBitmap/roaring` 是一种bitmap的扩展，更能压缩空间。

> 后面补充一下roaring bitmap的数据结构


### 2. mapreduce构建索引

第一版的mapreduce是使用的是grpc进行处理worker节点。第二版采用chan和goroutine来处理worker，以减少rpc的调用。

```go
_, _ = mapreduce.MapReduce(func(source chan<- []byte) {
    // 输入的文件
    for _, path := range req.FilePath {
        content, _ := os.ReadFile(path)
        source <- content
    }
}, func(item []byte, writer mapreduce.Writer[[]*types.KeyValue], cancel func(error)) {
    // 控制并发
    var wg sync.WaitGroup
    ch := make(chan struct{}, 3)

    keyValueList := make([]*types.KeyValue, 0, 1e3)
    lines := strings.Split(string(item), "\r\n")
    for _, line := range lines[1:] {
        ch <- struct{}{}
        wg.Add(1)
        // 输入的line结构
        docStruct, _ := input_data.Doc2Struct(line) // line 转 docs struct
        if docStruct.DocId == 0 {
            continue
        }

        // 分词
        tokens, _ := analyzer.GseCutForBuildIndex(docStruct.DocId, docStruct.Body)
        for _, v := range tokens {
            if v.Token == "" || v.Token == " " {
                continue
            }
            keyValueList = append(keyValueList, &types.KeyValue{
                    Key: v.Token, 
                    Value: cast.ToString(v.DocId)
                })
            // 前缀树的插入
            dictTrie.Insert(v.Token)
        }

        // 建立正排索引
        go func(docStruct *types.Document) {
            err = input_data.DocData2Kfk(docStruct)
            if err != nil {
                logs.LogrusObj.Error(err)
            }
            defer wg.Done()
            <-ch
        }(docStruct)
    }
    wg.Wait()
    // 排序shuffle操作
    sort.Sort(types.ByKey(keyValueList))
    writer.Write(keyValueList)
}, func(pipe <-chan []*types.KeyValue, writer mapreduce.Writer[string], cancel func(error)) {
    for values := range pipe {
        for _, v := range values { // 构建倒排索引
            if value, ok := invertedIndex.Get(v.Key); ok {
                value.AddInt(cast.ToInt(v.Value))
                invertedIndex.Set(v.Key, value)
            } else {
                docIds := roaring.NewBitmap()
                docIds.AddInt(cast.ToInt(v.Value))
                invertedIndex.Set(v.Key, docIds)
            }
        }
    }
})
```

map操作中，我们拆分出所有的词以及对应的id，如以下数据结构

eg: 

1:那里湖面总是澄清
2:那里空气充满宁静

```
那里:1
湖面:1
总是:1
澄清:1
那里:2
空气:2
充满:2
宁静:2
```

reduce操作中,我们聚合所有相同value的结构,构造出倒排索引的机构

```
那里:1,2
湖面:1
总是:1
澄清:1
空气:2
充满:2
宁静:2
```

另外在构造的过程中，也将消息也生产给kafka进行消费，下游进行kafka的监听，并生成mysql的正排索引，一级milvus的向量索引。

```go
go func(docStruct *types.Document) {
    err = input_data.DocData2Kfk(docStruct)
    if err != nil {
        logs.LogrusObj.Error(err)
    }
    defer wg.Done()
    <-ch
}(docStruct)
```

### 3. 存储结构

在存储索引的过程中，我们也采用了异步处理，但注意这里我们需要克隆一个context，否则容易造成`context cancel`,因为这个ctx是主进程的，主进程结束了自然就cancel了，所以我们需要clone一个新的，来进行传递。具体克隆的代码在 `pkg/clone/context.go` 这部分会另外开一篇详细来讲.

```go
go func() {
    newCtx := clone.NewContextWithoutDeadline()
    newCtx.Clone(ctx)
    err = storeInvertedIndexByHash(newCtx, invertedIndex)
    if err != nil {
        logs.LogrusObj.Error("storeInvertedIndexByHash error ", err)
    }
}()
```

boltdb存储

```go
func storeInvertedIndexByHash(ctx context.Context, invertedIndex cmap.ConcurrentMap[string, *roaring.Bitmap]) (err error) {
	dir, _ := os.Getwd()
	outName := fmt.Sprintf("%s/%s.%s", dir, timeutils.GetNowTime(), cconsts.InvertedBucket)
	invertedDB := storage.NewInvertedDB(outName)
	// 对所有的key进行存储
	for k, val := range invertedIndex.Items() {
		outByte, errx := val.MarshalBinary()
		if errx != nil {
			logs.LogrusObj.Error("storeInvertedIndexByHash-MarshalBinary", errx)
			continue
		}
		err = invertedDB.StoragePostings(k, outByte)
		if err != nil {
			logs.LogrusObj.Error("storeInvertedIndexByHash-StoragePostings", err)
			continue
		}
	}
	invertedDB.Close()
	err = redis.PushInvertedPath(ctx, redis.InvertedIndexDbPathDayKey, []string{outName})
	if err != nil {
		logs.LogrusObj.Error(err)
		return
	}
	return
}
```

### 4. 前缀树

> 前缀树相关的在 `pkg/trie/trie.go` 文件中

定义前缀树结构

```go
dictTrie := trie.NewTrie()  // 前缀树
```

定义前缀树的node结构, 由于 `cmap.ConcurrentMap[string, *TrieNode]` 没法反序列化，可以序列化json成string格式进行存储，但无法从string反序列化json进行读取，于是引入`ChildrenRecall`来处理召回请求。这个召回章节再讲解。

```go
type TrieNode struct {
	IsEnd          bool                                  `json:"is_end"`   // 标记该节点是否为一个单词的末尾
	Children       cmap.ConcurrentMap[string, *TrieNode] `json:"children"` // 存储子节点的指针
	ChildrenRecall map[string]*TrieNode                  `json:"children_recall"`
}
```

插入前缀树

```go
func (trie *Trie) Insert(word string) {
	words := []rune(word)
	node := trie.Root
	for i := 0; i < len(words); i++ {
		c := string(words[i])
		if _, ok := node.Children.Get(c); !ok {
			node.Children.Set(c, NewTrieNode())
		}
		node, _ = node.Children.Get(c)
	}
	node.IsEnd = true
}
```

查询前缀树

```go
func (trie *Trie) FindAllByPrefix(prefix string) []string {
	prefixs := []rune(prefix)
	node := trie.Root
	for i := 0; i < len(prefixs); i++ {
		c := string(prefixs[i])
		if _, ok := node.Children.Get(c); !ok {
			return nil
		}
		node, _ = node.Children.Get(c)
	}
	words := make([]string, 0)
	trie.dfs(node, prefix, &words)
	return words
}
```