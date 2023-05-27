package model

type FavoriteDetail struct {
	FavoriteDetailID int64       `gorm:"primarykey"`
	UserID           int64       // 用户id
	UrlID            int64       // url的id
	Url              string      // url地址
	Desc             string      // url的描述
	Favorite         []*Favorite `gorm:"many2many:f_to_fd;"`
}
