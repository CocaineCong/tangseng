// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package dao

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	favoritePb "github.com/CocaineCong/tangseng/idl/pb/favorite"
	"github.com/CocaineCong/tangseng/repository/mysql/db"
	"github.com/CocaineCong/tangseng/repository/mysql/model"
)

type FavoriteDao struct {
	db *gorm.DB
}

func NewFavoriteDao(ctx context.Context) *FavoriteDao {
	return &FavoriteDao{db.NewDBClient(ctx)}
}

func (dao *FavoriteDao) ListFavorite(req *favoritePb.FavoriteListReq) (r []*model.Favorite, err error) {
	err = dao.db.Model(&model.Favorite{}).
		Where("user_id = ?", req.UserId).Find(&r).Error
	if err != nil {
		return r, errors.Wrapf(err, "failed to query favorite list, userId = %v ", req.UserId)
	}
	return
}

func (dao *FavoriteDao) CreateFavorite(req *favoritePb.FavoriteCreateReq) (err error) {
	favorite := model.Favorite{
		FavoriteName: req.FavoriteName,
		UserID:       req.UserId,
	}
	if err = dao.db.Create(&favorite).Error; err != nil {
		return errors.Wrapf(err, "failed to create favorite, userId = %v ", req.UserId)
	}

	return
}

func (dao *FavoriteDao) DeleteFavorite(req *favoritePb.FavoriteDeleteReq) (err error) {
	err = dao.db.Where("favorite_id = ?", req.FavoriteId).
		Delete(model.Favorite{}).Error
	if err != nil {
		return errors.Wrapf(err, "failed to delete favorite, favoriteId = %v", req.FavoriteId)
	}
	return
}

func (dao *FavoriteDao) UpdateFavorite(req *favoritePb.FavoriteUpdateReq) (err error) {
	fMap := make(map[string]interface{})
	fMap["favorite_name"] = req.FavoriteName
	err = dao.db.Where("favorite_id = ?", req.FavoriteId).Updates(&fMap).Error
	if err != nil {
		return errors.Wrapf(err, "failed to update favorite, favoriteId = %v ", req.FavoriteId)
	}

	return
}
