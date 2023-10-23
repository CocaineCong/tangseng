package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CocaineCong/tangseng/app/gateway/rpc"
	pb "github.com/CocaineCong/tangseng/idl/pb/search_vector"
	"github.com/CocaineCong/tangseng/pkg/ctl"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// SearchVector 向量搜索
func SearchVector(ctx *gin.Context) {
	var req pb.SearchVectorRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.LogrusObj.Errorf("SearchVector-ShouldBind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}

	r, err := rpc.SearchVector(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("SearchVector:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "SearchVector RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}
