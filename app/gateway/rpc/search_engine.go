package rpc

import (
	"context"
	"errors"

	"github.com/CocaineCong/tangseng/consts/e"
	pb "github.com/CocaineCong/tangseng/idl/pb/search_engine"
)

func SearchEngineSearch(ctx context.Context, req *pb.SearchEngineRequest) (r *pb.SearchEngineResponse, err error) {
	r, err = SearchEngineClient.SearchEngineSearch(ctx, req)
	if err != nil {
		return
	}

	if r.Code != e.SUCCESS {
		err = errors.New(r.Msg)
		return
	}

	return
}
