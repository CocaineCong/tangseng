# Image Retrieval

基于ResNet50的以图搜图

## 项目结构

```shell
search_vector/
├── cirtorch
│   ├── datasets
│   ├── examples
│   ├── layers
│   ├── networks
│   └── utils
├── config
├── consts
├── ctl
├── etcd_operate
├── index
├── kafka_operate
├── lshash
├── milvus
├── service
├── utils
└── weights
```

其中，`cirtorch`部分来自于[CNN Image Retrieval in PyTorch](https://github.com/filipradenovic/cnnimageretrieval-pytorch)，使用了该项目的网络架构和预训练模型进行特征编码。

`ImageRetrieval\jpg`文件夹下存放用于进行查找的图像库，在本处，选用从[悟空数据集](https://wukong-dataset.github.io/wukong-dataset/index.html)的`Wukong100m`中爬取的20000张图片。

`index`文件夹下存放从图片库中抽取的特征信息以及LSH索引信息，LSH索引地址为[dataset_index_wukong.pkl](https://pan.baidu.com/s/1t_BXCGVEO0U_9tVCHnY5pw?pwd=e1fa)。

`lshash`部分使用[LSHash](https://github.com/kayzhu/LSHash)的代码，使用局部敏感哈希以加快检索速度。

`utils\retrieval_feature.py`部分为通过预训练的模型进行特征抽取，并使用LSH计算索引，并将特征数据和索引数据保存到本地。

`weights`目录下保存所使用的预训练模型，本项目中采用的是[CNN Image Retrieval in PyTorch](https://github.com/filipradenovic/cnnimageretrieval-pytorch)中使用ResNet50，Pooling层使用GeM，在`google-landmarks-2018 (gl18)`数据集上进行预训练的模型，
模型地址为[gl18-tl-resnet50-gem-w](http://cmp.felk.cvut.cz/cnnimageretrieval/data/networks/gl18/gl18-tl-resnet50-gem-w-83fdc30.pth)。 这里记得看目录 app/search_vector 下的README.md文件，下载模型放到当前的这个 weights 目录

`config.yaml`中保存着项目的一些基本配置，包括Flask的地址和端口，以及使用的预训练模型名称。

`interface.py`是项目Web端的代码。

## 项目主要依赖

可查看根目录下的`requirement.txt`文件

## 权重下载

夸克云盘权重下载链接：<https://pan.quark.cn/s/8526d0c23356>
下载完放到 `weights/` 目录下
