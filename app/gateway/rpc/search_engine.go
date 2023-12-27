package rpc

import (
	"context"
	"github.com/pkg/errors"

	"github.com/CocaineCong/tangseng/consts/e"
	pb "github.com/CocaineCong/tangseng/idl/pb/search_engine"
)

func SearchEngineSearch(ctx context.Context, req *pb.SearchEngineRequest) (r *pb.SearchEngineResponse, err error) {
	r, err = SearchEngineClient.SearchEngineSearch(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "SearchEngineClient.SearchEngineSearch error")
		return
	}

	if r.Code != e.SUCCESS {
		err = errors.Wrap(errors.New(r.Msg), "r.Code is unsuccessful")
		return
	}

	return
}

func WordAssociation(ctx context.Context, req *pb.SearchEngineRequest) (r *pb.WordAssociationResponse, err error) {
	r, err = SearchEngineClient.WordAssociation(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "SearchEngineClient.WordAssociation error")
		return
	}

	if r.Code != e.SUCCESS {
		err = errors.Wrap(errors.New(r.Msg), "r.Code is unsuccessful")
		return
	}

	return
}
