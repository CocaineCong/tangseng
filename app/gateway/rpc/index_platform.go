package rpc

import (
	"context"

	pb "github.com/CocaineCong/tangseng/idl/pb/index_platform"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// BuildIndex 建立索引的RPC调用
func BuildIndex(ctx context.Context, req *pb.BuildIndexReq) (resp *pb.BuildIndexResp, err error) {
	resp, err = IndexPlatformClient.BuildIndexService(ctx, req)
	if err != nil {
		log.LogrusObj.Error("BuildIndex-BuildIndexService", err)
		return
	}

	return
}
