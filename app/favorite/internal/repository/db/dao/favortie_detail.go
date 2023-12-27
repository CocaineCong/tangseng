package dao

import (
	"context"
	"github.com/pkg/errors"

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
		return errors.Wrapf(err, "failed to query favorite, favoriteId = %v ", req.FavoriteId)
	}

	fd := model.FavoriteDetail{
		UserID:   req.UserId,
		UrlID:    req.UrlId,
		Url:      req.Url,
		Desc:     req.Desc,
		Favorite: f,
	}
	err = dao.DB.Model(&model.FavoriteDetail{}).Create(&fd).Error
	if err != nil {
		return errors.Wrapf(err, "failed to create favoriteDetail，userID = %v,urlId = %v", req.UserId, req.UrlId)
	}
	return
}

func (dao *FavoriteDetailDao) ListFavoriteDetail(req *favoritePb.FavoriteDetailListReq) (r []*model.Favorite, err error) {
	var f []*model.Favorite
	err = dao.DB.Where("user_id = ?", req.UserId).Find(&f).Error
	if err != nil {
		return r, errors.Wrapf(err, "failed to query favorite, userId = %v ", req.UserId)
	}
	for _, v := range f {
		err = dao.DB.Model(&v).Association("FavoriteDetail").Find(&v.FavoriteDetail)
		r = append(r, v)
	}
	if err != nil {
		err = errors.Wrapf(err, "failed to query favoriteDetail")
	}
	return
}

func (dao *FavoriteDetailDao) DeleteFavoriteDetail(req *favoritePb.FavoriteDetailDeleteReq) (err error) {
	var f model.Favorite
	var fd model.FavoriteDetail
	err = dao.DB.Where("favorite_id = ?", req.FavoriteId).First(&f).Error
	if err != nil {
		return errors.Wrapf(err, "failed to query favorite, favoriteId = %v ", req.FavoriteId)
	}
	err = dao.DB.Where("favorite_detail_id = ?", req.FavoriteDetailId).First(&fd).Error
	if err != nil {
		return errors.Wrapf(err, "failed to query favoriteDetail, favoriteDetailId = %v ", req.FavoriteDetailId)
	}
	err = dao.DB.Model(&f).Association("FavoriteDetail").Delete(&fd)
	if err != nil {
		return errors.Wrapf(err, "failed to delete favoriteDetail, favoriteDetailId = %v ", req.FavoriteDetailId)
	}
	return
}
