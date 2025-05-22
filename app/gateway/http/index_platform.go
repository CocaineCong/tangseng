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
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/CocaineCong/tangseng/app/gateway/rpc"
	"github.com/CocaineCong/tangseng/consts/e"
	pb "github.com/CocaineCong/tangseng/idl/pb/index_platform"
	"github.com/CocaineCong/tangseng/pkg/ctl"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func BuildIndexByFiles(ctx *gin.Context) {
	var req pb.BuildIndexReq
	if err := ctx.ShouldBind(&req); err != nil {
		log.LogrusObj.Errorf("Bind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}

	r, err := rpc.BuildIndex(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("rpc.BuildIndex failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "BuildIndexByFiles RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

func UploadIndexByFiles(ctx *gin.Context) {
	var req pb.BuildIndexReq
	if err := ctx.ShouldBind(&req); err != nil {
		log.LogrusObj.Errorf("Bind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	file, fileHeader, _ := ctx.Request.FormFile("file")
	if fileHeader == nil {
		err := errors.New(e.GetMsg(e.ErrorUploadFile))
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "上传错误"))
		log.LogrusObj.Error(err)
		return
	}
	r, err := rpc.UploadByStream(ctx, &req, file, fileHeader.Size)
	if err != nil {
		log.LogrusObj.Errorf("rpc.BuildIndex failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "UploadIndexByFiles RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}
