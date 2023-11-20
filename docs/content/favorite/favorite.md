# 收藏模块

## 项目结构

```shell
/favorite
├── cmd               // 启动器
└── internal    
    ├── repository    // 存储仓库
    │   └── db        // 存储db操作
    │       └── dao
    └── service       // 具体实现的微服务 
```

## 表结构

其实收藏模块比较简单，只是有一个多对多的模块需要了解一下

收藏夹模块

```go
type Favorite struct {
    FavoriteID     int64             `gorm:"primarykey"` // 收藏夹id
    UserID         int64             `gorm:"index"`      // 用户id
    FavoriteName   string            `gorm:"unique"`     // 收藏夹名字
    FavoriteDetail []*FavoriteDetail `gorm:"many2many:f_to_fd;"`
}
```

具体的收藏细节模块

```go
type FavoriteDetail struct {
    FavoriteDetailID int64       `gorm:"primarykey"`
    UserID           int64       // 用户id
    UrlID            int64       // url的id
    Url              string      // url地址
    Desc             string      // url的描述
    Favorite         []*Favorite `gorm:"many2many:f_to_fd;"`
}
```

当我们执行`migrate`操作之后，就会在 `Favorite` 和 `FavoriteDetail` 之间建立一张`f_to_fd`多对多的关系表，这张表中，就记录了Favorite和FavoriteDetail两者的关联关系。

## 多对多的CURD操作

具体代码在`app/favorite/internal/repository/db/dao/`下

> 当然在ToC的应用中，是不允许有任何的外键关联的！虽然我们生产上不用，但也是要掌握的！
