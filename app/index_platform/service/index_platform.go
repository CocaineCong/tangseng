package service

import (
	"context"
	"sync"

	"github.com/CocaineCong/tangseng/app/index_platform/mapreduce/input_data_mr"

	"github.com/CocaineCong/tangseng/app/index_platform/consts"
	"github.com/CocaineCong/tangseng/app/index_platform/service/woker"
	"github.com/CocaineCong/tangseng/consts/e"
	pb "github.com/CocaineCong/tangseng/idl/pb/index_platform"
)

type IndexPlatformSrv struct {
	pb.UnimplementedIndexPlatformServiceServer
}

var (
	IndexPlatIns  *IndexPlatformSrv
	IndexPlatOnce sync.Once
)

func GetIndexPlatformSrv() *IndexPlatformSrv {
	IndexPlatOnce.Do(func() {
		IndexPlatIns = new(IndexPlatformSrv)
	})
	return IndexPlatIns
}

// BuildIndexService 构建索引
func (s *IndexPlatformSrv) BuildIndexService(ctx context.Context, req *pb.BuildIndexReq) (resp *pb.BuildIndexResp, err error) {
	// 时间估计
	resp = new(pb.BuildIndexResp)
	resp.Code = e.SUCCESS
	resp.Message = e.GetMsg(e.SUCCESS)
	// TODO: 最后改成异步
	_ = NewMaster(req.FilePath, consts.ReduceDefaultNum)
	woker.Worker(ctx, input_data_mr.Map, input_data_mr.Reduce)

	return
}
