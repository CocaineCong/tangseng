package rpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/CocaineCong/tangseng/consts/e"
	pb "github.com/CocaineCong/tangseng/idl/pb/search_vector"
)

func SearchVector(ctx context.Context, req *pb.SearchVectorRequest) (resp *pb.SearchVectorResponse, err error) {
	fmt.Println("req", req)
	resp, err = SearchVectorClient.SearchVector(ctx, req)
	if err != nil {
		return
	}

	if resp.Code != e.SUCCESS {
		err = errors.New(resp.Msg)
		return
	}

	return
}
