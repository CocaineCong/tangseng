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

package service

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/CocaineCong/tangseng/app/user/internal/repository/db/dao"
	e2 "github.com/CocaineCong/tangseng/consts/e"
	pb "github.com/CocaineCong/tangseng/idl/pb/user"
	"github.com/CocaineCong/tangseng/repository/mysql/model"
)

var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

type UserSrv struct {
	pb.UnimplementedUserServiceServer
}

func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

func (u *UserSrv) UserLogin(ctx context.Context, req *pb.UserLoginReq) (resp *pb.UserDetailResponse, err error) {
	resp = new(pb.UserDetailResponse)
	resp.Code = e2.SUCCESS
	r, err := dao.NewUserDao(ctx).GetUserInfo(req)
	if err != nil {
		resp.Code = e2.ERROR
		err = errors.WithMessage(err, "getUserInfo error")
		return
	}
	resp.UserDetail = &pb.UserResp{
		UserId:   r.UserID,
		UserName: r.UserName,
		NickName: r.NickName,
	}
	return
}

func (u *UserSrv) UserRegister(ctx context.Context, req *pb.UserRegisterReq) (resp *pb.UserCommonResponse, err error) {
	resp = new(pb.UserCommonResponse)
	resp.Code = e2.SUCCESS
	user := &model.User{
		UserName: req.UserName,
		NickName: req.NickName,
	}
	_ = user.SetPassword(req.Password)
	err = dao.NewUserDao(ctx).CreateUser(user)
	if err != nil {
		resp.Code = e2.ERROR
		err = errors.WithMessage(err, "createUser error")
		return
	}
	resp.Data = e2.GetMsg(int(resp.Code))
	return
}
