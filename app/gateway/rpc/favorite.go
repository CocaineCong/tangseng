package rpc

import (
	"context"
	"github.com/pkg/errors"

	"github.com/CocaineCong/tangseng/consts/e"
	favoritePb "github.com/CocaineCong/tangseng/idl/pb/favorite"
)

func FavoriteCreate(ctx context.Context, req *favoritePb.FavoriteCreateReq) (resp *favoritePb.FavoriteCommonResponse, err error) {
	resp, err = FavoriteClient.FavoriteCreate(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "FavoriteClient.FavoriteCreate error")
		return
	}
	if resp.Code != e.SUCCESS {
		err = errors.Wrap(errors.New(resp.Error), "resp.Code is not success")
		return
	}

	return
}

func FavoriteUpdate(ctx context.Context, req *favoritePb.FavoriteUpdateReq) (resp *favoritePb.FavoriteCommonResponse, err error) {
	resp, err = FavoriteClient.FavoriteUpdate(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "FavoriteClient.FavoriteUpdate error")
		return
	}

	if resp.Code != e.SUCCESS {
		err = errors.Wrap(errors.New(resp.Error), "resp.Code is not success")
		return
	}

	return
}

func FavoriteList(ctx context.Context, req *favoritePb.FavoriteListReq) (resp *favoritePb.FavoriteListResponse, err error) {
	resp, err = FavoriteClient.FavoriteList(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "FavoriteClient.FavoriteList error")
		return
	}
	if resp.Code != e.SUCCESS {
		err = errors.Wrap(errors.New("FavoriteList 出现错误"), "resp.Code is not success") // TODO 整个错误 proto
		return
	}

	return
}

func FavoriteDelete(ctx context.Context, req *favoritePb.FavoriteDeleteReq) (resp *favoritePb.FavoriteCommonResponse, err error) {
	resp, err = FavoriteClient.FavoriteDelete(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "FavoriteClient.FavoriteDelete error")
		return
	}
	if resp.Code != e.SUCCESS {
		err = errors.Wrap(errors.New(resp.Error), "resp.Code is not success")
		return
	}

	return
}

func FavoriteDetailList(ctx context.Context, req *favoritePb.FavoriteDetailListReq) (resp *favoritePb.FavoriteDetailListResponse, err error) {
	resp, err = FavoriteClient.FavoriteDetailList(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "FavoriteClient.FavoriteDetailList error")
		return
	}
	if resp.Code != e.SUCCESS {
		err = errors.Wrap(errors.New("出现错误"), "resp.Code is not success")
		return
	}

	return
}

func FavoriteDetailDelete(ctx context.Context, req *favoritePb.FavoriteDetailDeleteReq) (resp *favoritePb.FavoriteCommonResponse, err error) {
	resp, err = FavoriteClient.FavoriteDetailDelete(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "FavoriteClient.FavoriteDetailDelete error")
		return
	}
	if resp.Code != e.SUCCESS {
		err = errors.Wrap(errors.New(resp.Error), "resp.Code is not success")
		return
	}

	return
}

func FavoriteDetailCreate(ctx context.Context, req *favoritePb.FavoriteDetailCreateReq) (resp *favoritePb.FavoriteCommonResponse, err error) {
	resp, err = FavoriteClient.FavoriteDetailCreate(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "FavoriteClient.FavoriteDetailCreate error")
		return
	}
	if resp.Code != e.SUCCESS {
		err = errors.Wrap(errors.New(resp.Error), "resp.Code is not success")
		return
	}

	return
}
