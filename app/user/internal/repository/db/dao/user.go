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

	userPb "github.com/CocaineCong/tangseng/idl/pb/user"
	"github.com/CocaineCong/tangseng/repository/mysql/db"
	"github.com/CocaineCong/tangseng/repository/mysql/model"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{db.NewDBClient(ctx)}
}

// GetUserInfo 获取用户信息
func (dao *UserDao) GetUserInfo(req *userPb.UserLoginReq) (r *model.User, err error) {
	err = dao.db.Model(&model.User{}).Where("user_name = ?", req.UserName).
		First(&r).Error
	if err != nil {
		err = errors.Wrapf(err, "failed to get user info, userName = %v", req.UserName)
	}
	return
}

// CreateUser 用户创建
func (dao *UserDao) CreateUser(in *model.User) (err error) {
	var count int64
	dao.db.Model(&model.User{}).Where("user_name = ?", in.UserName).Count(&count)
	if count != 0 {
		return errors.Wrapf(errors.New("UserName Exist"), "failed to create user, userName = %v", in.UserName)
	}
	if err = dao.db.Model(&model.User{}).Create(&in).Error; err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return
}
