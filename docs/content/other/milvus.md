# milvus 向量数据库

> 多路召回中，我们需要召回各种各样的数据结构类型的索引，包括但不限于`倒排索引`,`kvv索引`,`向量索引`等等。而在tangseng中的多路召回支持了`倒排索引`和`向量索引`两种模式的召回。

tangseng的向量数据库的选型是 milvus，一款基于go语言的向量数据库，非常强大。但其实向量数据库的操作都在tangseng的python部分，主要是因为文本向量化的操作基本是python垄断...

所以我们在go中将query传给python，python进行向量化，然后查询milvus，再返回给go。

## 具体操作

### 向量化操作

我们直接加载模型就好了，很简单！然后对query进行向量化操作

```python
from sentence_transformers import SentenceTransformer

TRANSFORMER_MODEL = SentenceTransformer(TRANSFORMER_MODEL_NAME)
query_feat = TRANSFORMER_MODEL.encode(query)
```

### milvus的CURD操作

代码在`app/search_vector/milvus/`下

定义一个`Milvus`类，初始化时传入`MILVUS_HOST`,`MILVUS_PORT` 并在这个类上实现操作一系列的CURD操作。

```python
class Milvus:

    def __init__(self, host=MILVUS_HOST, port=MILVUS_PORT):
        try:
            self.collection = None
            connections.connect(host=host, port=port)
        except Exception as e:
            LOGGER.error(f"init milvus {e}")
            sys.exit(1)
```

传入`collection_name`创建一个`collection`，**milvus的collection概念有点类似mysql的table** , 在创建collection的时候，我们需要制定字段，也就是FieldSchema，就类似mysql的表的字段，**但是milvus中必须要有一个字段是向量类型的。**

```python
def create_collection(self, collection_name):
    try:
        if not self.has_collection(collection_name):
            doc_id = FieldSchema(name='doc_id', dtype=DataType.INT64,
                                    description='doc_id', max_length=500,
                                    is_primary=True, auto_id=False)
            # One and only vector can be created.
            # tow or zero will be wrong.
            body_vec = FieldSchema(name='embedding', dtype=DataType.FLOAT_VECTOR,
                                    description='body embedding vectors',
                                    dim=VECTOR_DIMENSION, is_primary=False)
            schema = CollectionSchema(fields=[doc_id, body_vec], description='vec_info')
            self.collection = Collection(name=collection_name, schema=schema)
            self.create_index(collection_name)
            LOGGER.info(f"create milvus collection:{collection_name}")
        else:
            self.set_collection(collection_name)
        return "ok"
    except Exception as e:
        LOGGER.error(f"failed to create collection {e}")
        sys.exit(1)
```

索引创建, 在创建的时候需要传入几个参数，`index_type`,`metric_type` 和 `nprobe`。

解释一下这两个参数的意思

- index_type: 索引类型,设置索引类型为"IVF_FLAT",这是一种基于倒排文件（IVF）的索引类型，它通过扁平扫描（FLAT）来实现精确的距离计算。这种索引类型适用于中等大小的数据集。
- metric_type: 设置距离度量方式。一般使用常用的度量向量相似性的方法。
- nprobe:设置索引的参数。这里我们设置参数 "nlist" 为 16，它表示倒排文件中创建 16 个倒排列表（inverted lists）。较大的 nlist 值有助于提高搜索速度，但会增加索引的内存消耗。


```python
def create_index(self, collection_name):
    try:
        self.set_collection(collection_name)
        default_index = {"index_type": "IVF_FLAT", "metric_type": METRIC_TYPE, "params": {"nlist": 16384}}
        print(default_index)
        status = self.collection.create_index(field_name="embedding", index_params=default_index)
        if not status.code:
            LOGGER.info(f"successfully create index in collection :{collection_name} with param:{default_index}")
            return status
        else:
            raise Exception(status.message)
    except Exception as e:
        print(e)
        LOGGER.error(f"failed to create index:{e}")
        sys.exit(1)
```

创建完index之后，我们可以直接插入数据,指定`collection_name`,以及doc id和需要插入的向量文本`body_vec`

```python
def insert(self, collection_name, doc_id, body_vec):
    try:
        data = [doc_id, body_vec]
        self.set_collection(collection_name)
        mr = self.collection.insert(data)
        ids = mr.primary_keys
        self.collection.load()
        LOGGER.info(
            f"insert vector to milvus in collection: {collection_name} with body {len(body_vec)}")
        return ids
    except Exception as e:
        LOGGER.error(f"insert to load data to milvus {e}")
        sys.exit(1)
```

搜索索引我们同样需要输入`metric_type`和`nprobe`这两个参数

```python
def search_vectors(self, collection_name, vectors, top_k):
    try:
        self.set_collection(collection_name)
        search_params = {"metric_type": METRIC_TYPE, "params": {"nprobe": 16}}
        res = self.collection.search(vectors, anns_field="embedding", param=search_params, limit=top_k)
        LOGGER.info(f"successfully search in collection:{res}")
        return res
    except Exception as e:
        LOGGER.error(f"failed to search vectors in milvus:{e}")
        sys.exit(1)
```
