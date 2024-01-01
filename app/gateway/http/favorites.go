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

package http

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CocaineCong/tangseng/app/gateway/rpc"
	pb "github.com/CocaineCong/tangseng/idl/pb/favorite"
	"github.com/CocaineCong/tangseng/pkg/ctl"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func ListFavorite(ctx *gin.Context) {
	var req pb.FavoriteListReq
	if err := ctx.Bind(&req); err != nil {
		log.LogrusObj.Errorf("Bind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		log.LogrusObj.Errorf("ctl.GetUserInfo failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteList(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("rpc.FavoriteList failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "FavoriteList RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

func CreateFavorite(ctx *gin.Context) {
	var req pb.FavoriteCreateReq
	if err := ctx.ShouldBind(&req); err != nil {
		log.LogrusObj.Errorf("ShouldBind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		log.LogrusObj.Errorf("ctl.GetUserInfo failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteCreate(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("rpc.FavoriteCreate failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "FavoriteCreateReq RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

func UpdateFavorite(ctx *gin.Context) {
	var req pb.FavoriteCreateReq
	if err := ctx.Bind(&req); err != nil {
		log.LogrusObj.Errorf("Bind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		log.LogrusObj.Errorf("ctl.GetUserInfo failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteCreate(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("rpc.FavoriteCreate failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "UpdateFavorite RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

func DeleteFavorite(ctx *gin.Context) {
	var req pb.FavoriteDeleteReq
	if err := ctx.Bind(&req); err != nil {
		log.LogrusObj.Errorf("req:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		log.LogrusObj.Errorf("ctl.GetUserInfo failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteDelete(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("rpc.FavoriteDelete failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "DeleteFavorite RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

func ListFavoriteDetail(ctx *gin.Context) {
	var req pb.FavoriteDetailListReq
	if err := ctx.Bind(&req); err != nil {
		log.LogrusObj.Errorf("Bind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		log.LogrusObj.Errorf("ctl.GetUserInfo failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteDetailList(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("rpc.FavoriteDetailList failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "FavoriteDetailList RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

func CreateFavoriteDetail(ctx *gin.Context) {
	var req pb.FavoriteDetailCreateReq
	if err := ctx.Bind(&req); err != nil {
		log.LogrusObj.Errorf("Bind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		log.LogrusObj.Errorf("ctl.GetUserInfo failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteDetailCreate(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("rpc.FavoriteDetailCreate failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "FavoriteDetailCreate RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

func DeleteFavoriteDetail(ctx *gin.Context) {
	var req pb.FavoriteDetailDeleteReq
	if err := ctx.Bind(&req); err != nil {
		log.LogrusObj.Errorf("Bind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		log.LogrusObj.Errorf("ctl.GetUserInfo failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteDetailDelete(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("rpc.FavoriteDetailDelete failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "FavoriteDetailDelete RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}
