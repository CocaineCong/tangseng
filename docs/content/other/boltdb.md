# boltdb kv store

具体代码在`app/index_platform/repository/storage/`、`app/search_engine/repository/storage/`下

> boltdb 是一款键值对存储介质，非常简单的单机架构。简单到有一把全局锁，如果有程序对某个boltdb文件进行读取或写入，其他程序想要对这个boltdb进行操作，将会被阻塞！

## boltdb 的CURD操作

定义全局倒排索引存储的boltdb

```go
var GlobalInvertedDB []*InvertedDB

type InvertedDB struct {
    file   *os.File
    db     *bolt.DB
    offset int64
}
```

连接boltdb

```go
// InitInvertedDB 初始化倒排索引库
func InitInvertedDB(ctx context.Context) []*InvertedDB {
    dbs := make([]*InvertedDB, 0)
    // 从redis中读取所有的boltdb路径
    filePath, _ := redis.ListInvertedPath(ctx, redis.InvertedIndexDbPathKeys)
    for _, file := range filePath {
        f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
        if err != nil {
            log.LogrusObj.Error(err)
        }
        stat, err := f.Stat()
        if err != nil {
            log.LogrusObj.Error(err)
        }
        db, err := bolt.Open(file, 0600, nil)
        if err != nil {
            log.LogrusObj.Error(err)
        }
        dbs = append(dbs, &InvertedDB{f, db, stat.Size()})
    }
    if len(filePath) == 0 {
        return nil
    }
    GlobalInvertedDB = dbs
    return nil
}
```

写入kv

```go
// Put 通过bolt写入数据
func Put(db *bolt.DB, bucket string, key []byte, value []byte) error {
    return db.Update(func(tx *bolt.Tx) error {
        b, err := tx.CreateBucketIfNotExists([]byte(bucket))
        if err != nil {
            return err
        }
        return b.Put(key, value)
    })
}
```

获取数据

```go
// Get 通过bolt获取数据
func Get(db *bolt.DB, bucket string, key []byte) (r []byte, err error) {
    err = db.View(func(tx *bolt.Tx) (err error) {
        b := tx.Bucket([]byte(bucket))
        if b == nil {
            b, _ = tx.CreateBucketIfNotExists([]byte(bucket))
        }
        r = b.Get(key)
        if r == nil { // 如果是空的话，直接创建这个key，然后返回这个key的初始值，也就是0
            r = []byte("0")
            return
        }
        return
    })

    return
}
```
