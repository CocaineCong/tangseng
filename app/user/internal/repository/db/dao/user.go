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
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{db.NewDBClient(ctx)}
}

// GetUserInfo 获取用户信息
func (dao *UserDao) GetUserInfo(req *userPb.UserLoginReq) (r *model.User, err error) {
	err = dao.Model(&model.User{}).Where("user_name=?", req.UserName).
		First(&r).Error
	if err != nil {
		err = errors.Wrapf(err, "failed to get user info, userName = %v", req.UserName)
	}
	return
}

// CreateUser 用户创建
func (dao *UserDao) CreateUser(req *userPb.UserRegisterReq) (err error) {
	var user model.User
	var count int64
	dao.Model(&model.User{}).Where("user_name = ?", req.UserName).Count(&count)
	if count != 0 {
		return errors.Wrapf(errors.New("UserName Exist"), "failed to create user, userName = %v", req.UserName)
	}

	user = model.User{
		UserName: req.UserName,
		NickName: req.NickName,
	}
	_ = user.SetPassword(req.Password)
	if err = dao.Model(&model.User{}).Create(&user).Error; err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return
}
