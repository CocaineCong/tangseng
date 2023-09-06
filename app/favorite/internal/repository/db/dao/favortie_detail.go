package dao

import (
	"context"

	"gorm.io/gorm"

	favoritePb "github.com/CocaineCong/tangseng/idl/pb/favorite"
	"github.com/CocaineCong/tangseng/repository/mysql/db"
	"github.com/CocaineCong/tangseng/repository/mysql/model"
)

type FavoriteDetailDao struct {
	*gorm.DB
}

func NewFavoriteDetailDao(ctx context.Context) *FavoriteDetailDao {
	return &FavoriteDetailDao{db.NewDBClient(ctx)}
}

// CreateFavoriteDetail 收藏夹可以重复收藏
func (dao *FavoriteDetailDao) CreateFavoriteDetail(req *favoritePb.FavoriteDetailCreateReq) (err error) {
	var f []*model.Favorite
	err = dao.DB.Where("favorite_id = ?", req.FavoriteId).Find(&f).Error
	if err != nil {
		return
	}

	fd := model.FavoriteDetail{
		UserID:   req.UserId,
		UrlID:    req.UrlId,
		Url:      req.Url,
		Desc:     req.Desc,
		Favorite: f,
	}
	err = dao.DB.Model(&model.FavoriteDetail{}).Create(&fd).Error

	return
}

func (dao *FavoriteDetailDao) ListFavoriteDetail(req *favoritePb.FavoriteDetailListReq) (r []*model.Favorite, err error) {
	var f []*model.Favorite
	dao.DB.Where("user_id = ?", req.UserId).Find(&f)
	for _, v := range f {
		_ = dao.DB.Model(&v).Association("FavoriteDetail").Find(&v.FavoriteDetail)
		r = append(r, v)
	}

	return
}

func (dao *FavoriteDetailDao) DeleteFavoriteDetail(req *favoritePb.FavoriteDetailDeleteReq) error {
	var f model.Favorite
	var fd model.FavoriteDetail
	dao.DB.Where("favorite_id = ?", req.FavoriteId).First(&f)
	dao.DB.Where("favorite_detail_id = ?", req.FavoriteDetailId).First(&fd)
	err := dao.DB.Model(&f).Association("FavoriteDetail").Delete(&fd)
	return err
}
