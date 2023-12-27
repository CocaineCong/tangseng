package rpc

import (
	"context"
	"github.com/pkg/errors"

	pb "github.com/CocaineCong/tangseng/idl/pb/index_platform"
)

// BuildIndex 建立索引的RPC调用
func BuildIndex(ctx context.Context, req *pb.BuildIndexReq) (resp *pb.BuildIndexResp, err error) {
	resp, err = IndexPlatformClient.BuildIndexService(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "IndexPlatformClient.BuildIndexService err")
		return
	}

	return
}
