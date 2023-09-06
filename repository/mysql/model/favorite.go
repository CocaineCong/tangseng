package model

type Favorite struct {
	FavoriteID     int64             `gorm:"primarykey"` // 收藏夹id
	UserID         int64             `gorm:"index"`      // 用户id
	FavoriteName   string            `gorm:"unique"`     // 收藏夹名字
	FavoriteDetail []*FavoriteDetail `gorm:"many2many:f_to_fd;"`
}
