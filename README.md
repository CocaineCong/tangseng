# Tangseng 基于Go语言的搜索引擎

**[项目详细内容地址点击这里](https://cocainecong.github.io/tangseng/#/)** 

## 项目大体框架&功能

1. gin作为http框架，grpc作为rpc框架，etcd作为服务发现。
2. 总体服务分成`用户模块`、`收藏夹模块`、`索引平台`、`搜索引擎(文字模块)`、`搜索引擎(图片模块)`。
3. 分布式爬虫爬取数据，并发送到kafka集群中，再落库消费。 (虽然爬虫还没写，但不妨碍我画饼...) 
4. 搜索引擎模块的文本搜索单独设立使用boltdb存储index。
5. 使用trie tree实现词条联想。 
6. 图片搜索使用ResNet50来进行向量化查询 + Milvus or Faiss 向量数据库的查询 (开始做了... DeepLearning也太难了...)。
7. 支持多路召回，go中进行倒排索引召回，python进行向量召回。通过grpc调用连接，进行融合。
8. 支持TF-IDF，BM25等等算法排序。

![项目大体框架](docs/images/tangseng.png)

## 🧑🏻‍💻 前端地址

前端用的是 react, but still coding

[react-tangseng](https://github.com/CocaineCong/react-tangseng)

## 未来规划
### 架构相关

- [ ] 引入降级熔断
- [ ] 引入jaeger进行全链路追踪(go追踪到python)
- [ ] 引入skywalking or prometheus进行监控
- [ ] 抽离dao的init，用key来获取相关数据库实例

### 功能相关

- [x] 构建索引的时候太慢了.后面加上并发，建立索引的地方加上并发
- [ ] 索引压缩，inverted index，也就是倒排索引表，后续改成存offset,用mmap
- [x] 相关性的计算要考虑一下，TFIDF，bm25
- [x] 使用前缀树存储联想信息
- [ ] 哈夫曼编码压缩前缀树
- [ ] inverted 和 trie tree 的存储支持一致性hash分片存储
- [ ] 词向量，pagerank
- [ ] 分离 trie tree 的 build 和 recall 过程
- [x] 分词加入ik分词器
- [x] 构建索引平台，计算存储分离，构建索引与召回分开
- [ ] 并且差运算
- [ ] 分页，排序
- [ ] 纠正输入的query,比如“陆加嘴”-->“陆家嘴”
- [x] 输入进行词条可以进行联想，比如 “东方明” 提示--> “东方明珠”
- [x] 目前是基于块的索引方法，后续看看能不能改成分布式mapreduce来构建索引 (6.824 lab1)
- [ ] 在上一条的基础上再加上动态索引（还不知道上一条能不能实现...）
- [x] 改造倒排索引，使用 roaring bitmap 存储docid (好难)
- [ ] 实现TF类
- [ ] 所有的输入数据都收口到starrocks，从starrocks读取来构建索引
- [x] 搜索完一个接着搜索，没有清除缓存导致结果是和上一个产生并集
- [x] 排序器优化

![文本搜索](docs/images/text2text.jpg)

# 快速开始
## Python
1. 确保电脑已经安装了python

2. 安装venv环境

```shell
python -m venv venv
```

3. 激活 venv python 环境

macos:

```shell
source venv/bin/activate
```

windows:

等我清完C盘再兼容一下...还没在win上跑过...

## Golang

下载第三方依赖包

```shell
go mod tidy
```

目录下执行
```shell
make run-xxx(user,favortie ...)
```