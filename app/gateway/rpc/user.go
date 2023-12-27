package rpc

import (
	"context"
	"github.com/pkg/errors"

	"github.com/CocaineCong/tangseng/consts/e"
	userPb "github.com/CocaineCong/tangseng/idl/pb/user"
)

func UserLogin(ctx context.Context, req *userPb.UserLoginReq) (resp *userPb.UserDetailResponse, err error) {
	r, err := UserClient.UserLogin(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "UserClient.UserLogin error")
		return
	}

	if r.Code != e.SUCCESS {
		err = errors.Wrap(errors.New("登陆失败"), "r.Code is unsuccessful")
		return
	}

	return r, nil
}

func UserRegister(ctx context.Context, req *userPb.UserRegisterReq) (resp *userPb.UserCommonResponse, err error) {
	r, err := UserClient.UserRegister(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "UserClient.UserRegister error")
		return
	}

	if r.Code != e.SUCCESS {
		err = errors.Wrap(errors.New(r.Msg), "r.Code is unsuccessful")
		return
	}

	return
}
