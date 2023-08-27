package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CocaineCong/tangseng/app/gateway/rpc"
	pb "github.com/CocaineCong/tangseng/idl/pb/search_engine"
	"github.com/CocaineCong/tangseng/pkg/ctl"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// SearchEngineSearch 搜索
func SearchEngineSearch(ctx *gin.Context) {
	var req pb.SearchEngineRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.LogrusObj.Errorf("SearchEngineSearch-ShouldBind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}

	r, err := rpc.SearchEngineSearch(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("SearchEngineSearch:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "SearchEngineSearch RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

// WordAssociation 词条联想
func WordAssociation(ctx *gin.Context) {
	var req *pb.SearchEngineRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.LogrusObj.Errorf("WordAssociation-ShouldBind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}

	r, err := rpc.WordAssociation(ctx, req)
	if err != nil {
		log.LogrusObj.Errorf("WordAssociation:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "WordAssociation RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}
