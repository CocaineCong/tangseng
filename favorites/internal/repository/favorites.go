package repository

import (
	"errors"
	"favorites/internal/service"
	"favorites/pkg/util"
)


type Favorites struct {
	FavoriteID   uint   `gorm:"primarykey"` // 收藏夹id
	UserID       uint   `gorm:"index"`      // 用户id
	FavoriteName string `gorm:"unique"`     // 收藏夹名字
	FavoritesDetail     []FavoritesDetails `gorm:"many2many:f_to_fd;"`
}

func (favorite *Favorites) CheckExist(req *service.FavoritesRequest) bool {
	var count int64
	DB.Where("favorite_name=?", req.FavoriteName).Count(&count)
	if count != 0 {
		return true
	}
	return false
}

func (favorite *Favorites) Show (req *service.FavoritesRequest)(favoriteList []Favorites,err error) {
	err = DB.Model(Favorites{}).Where("user_id=?", req.UserID).Find(&favoriteList).Error
	if err != nil {
		return favoriteList, err
	}
	return favoriteList, nil
}

func (*Favorites) Create (req *service.FavoritesRequest) error {
	var favorite Favorites
	if exist := favorite.CheckExist(req); exist {
		return errors.New("UserName Not Exist")
	}
	favorite = Favorites{
		FavoriteName: req.FavoriteName,
		UserID:       uint(req.UserID),
	}
	if err := DB.Create(&favorite).Error; err != nil {
		util.LogrusObj.Error("Insert Favorite Error:" + err.Error())
		return err
	}
	return nil
}

func (*Favorites) Delete(req *service.FavoritesRequest) error {
	err := DB.Where("favorite_id=?", req.FavoriteID).Delete(Favorites{}).Error
	return err
}

func (*Favorites) Update(req *service.FavoritesRequest) error {
	f := Favorites{}
	if exist := f.CheckExist(req); exist {
		return errors.New("Favorite Name Exist!")
	}
	err := DB.Where("favorite_id=?", req.FavoriteID).First(&f).Error
	if err != nil {
		return err
	}
	f.FavoriteName = req.FavoriteName
	err = DB.Save(&f).Error
	return err
}

func BuildFavorites(item []Favorites) (fList []*service.FavoritesModel) {
	for _,v := range item{
		f := BuildFavorite(v)
		fList = append(fList,f)
	}
	return fList
}

func BuildFavorite(item Favorites) *service.FavoritesModel {
	urlInfo := []*service.UrlModel{}
	urlInfo = BuildUrlInfos(item.FavoritesDetail)
	return &service.FavoritesModel{
		FavoriteID:uint32(item.FavoriteID),
		FavoriteName:item.FavoriteName,
		UserID:uint32(item.UserID),
		UrlInfo: urlInfo,
	}
}

func BuildUrlInfo(item FavoritesDetails) *service.UrlModel {
	return &service.UrlModel{
		UrlID: uint32(item.UrlID),
		Url:   item.Url,
		Desc:  item.Desc,
	}
}

func BuildUrlInfos(item []FavoritesDetails) (urlList []*service.UrlModel) {
	for _,v:=range item{
		u := BuildUrlInfo(v)
		urlList = append(urlList,u)
	}
	return urlList
}