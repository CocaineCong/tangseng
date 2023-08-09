package http

import (
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
		log.LogrusObj.Errorf("GetUserInfo:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteList(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("FavoriteList:%v", err)
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
		log.LogrusObj.Errorf("GetUserInfo:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteCreate(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("FavoriteCreate:%v", err)
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
		log.LogrusObj.Errorf("GetUserInfo:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteCreate(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("FavoriteCreate:%v", err)
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
		log.LogrusObj.Errorf("GetUserInfo:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteDelete(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("FavoriteDelete:%v", err)
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
		log.LogrusObj.Errorf("GetUserInfo:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteDetailList(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("FavoriteDetailList:%v", err)
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
		log.LogrusObj.Errorf("GetUserInfo:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteDetailCreate(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("FavoriteDetailCreate:%v", err)
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
		log.LogrusObj.Errorf("GetUserInfo:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteDetailDelete(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("FavoriteDetailDelete:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "FavoriteDetailDelete RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}
