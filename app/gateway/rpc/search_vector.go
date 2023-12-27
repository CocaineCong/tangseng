package rpc

import (
	"context"
	"github.com/pkg/errors"

	"github.com/CocaineCong/tangseng/consts/e"
	pb "github.com/CocaineCong/tangseng/idl/pb/search_vector"
)

func SearchVector(ctx context.Context, req *pb.SearchVectorRequest) (resp *pb.SearchVectorResponse, err error) {
	resp, err = SearchVectorClient.SearchVector(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "SearchEngineClient.SearchVector error")
		return
	}

	if resp.Code != e.SUCCESS {
		err = errors.Wrap(errors.New(resp.Msg), "resp.Code is unsuccessful")
		return
	}

	return
}
