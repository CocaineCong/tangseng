package repository

import (
	"favorites/internal/service"
)

type FavoritesDetails struct {
	FavoritesDetailID uint `gorm:"primarykey"`
	UserID            uint   // 用户id
	UrlID             uint   // url的id
	Url               string // url地址
	Desc              string // url的描述
	Favorite          []Favorites `gorm:"many2many:f_to_fd;"`
}

type FavoritesDetailsResp struct {
	FavoriteID        uint
	UserID            uint
	FavoritesDetailID uint
	FavoritesName     string
	Url               string
	Desc              string
}

// 收藏夹可以重复收藏
func (*FavoritesDetails) Create(req *service.FavoritesRequest)(err error) {
	var f []Favorites
	DB.Where("favorite_id=?", req.FavoriteID).Find(&f)
	fd := FavoritesDetails{
		UserID:   uint(req.UserID),
		UrlID:    uint(req.UrlInfo.UrlID),
		Url:      req.UrlInfo.Url,
		Desc:     req.UrlInfo.Desc,
		Favorite: f,
	}
	err = DB.Model(FavoritesDetails{}).Create(&fd).Error
	return err
}

func (*FavoritesDetails) Show (req *service.FavoritesRequest) ([]Favorites,error) {
	f := []Favorites{}
	fRespList := []Favorites{}
	DB.Where("user_id=?", req.UserID).Find(&f)
	for _,v := range f {
		_ = DB.Model(&v).Association("FavoritesDetail").Find(&v.FavoritesDetail)
		fRespList = append(fRespList, v)
	}
	//wg := &sync.WaitGroup{}
	//for _, v := range f {
	//	wg.Add(1)
	//	go func(w *sync.WaitGroup) {
	//		wg.Add(1)
	//		for _, vt := range fd {
	//			if vt.Favorite[0].FavoriteID == v.FavoriteID {
	//				fResp := FavoritesDetailsResp{
	//					FavoriteID:        v.FavoriteID,
	//					UserID:            v.UserID,
	//					FavoritesDetailID: vt.FavoritesDetailID,
	//					FavoritesName:     v.FavoriteName,
	//					Url:               vt.Url,
	//					Desc:              vt.Desc,
	//				}
	//				fRespList = append(fRespList, fResp)
	//			}
	//		}
	//		wg.Done()
	//	}(wg)
	//}
	return fRespList, nil
}

func (*FavoritesDetails) Delete(req *service.FavoritesRequest) error {
	var f Favorites
	var fd FavoritesDetails
	DB.Where("favorite_id=?", req.FavoriteID).First(&f)
	DB.Where("favorites_detail_id=?", req.FavoriteDetailID).First(&fd)
	err := DB.Model(&f).Association("FavoritesDetail").Delete(&fd)
	return err
}