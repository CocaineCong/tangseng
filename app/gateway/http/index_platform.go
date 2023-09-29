package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CocaineCong/tangseng/app/gateway/rpc"
	pb "github.com/CocaineCong/tangseng/idl/pb/index_platform"
	"github.com/CocaineCong/tangseng/pkg/ctl"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func BuildIndexByFiles(ctx *gin.Context) {
	var req pb.BuildIndexReq
	if err := ctx.Bind(&req); err != nil {
		log.LogrusObj.Errorf("Bind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}

	r, err := rpc.BuildIndex(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("BuildIndexByFiles:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "BuildIndexByFiles RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}
