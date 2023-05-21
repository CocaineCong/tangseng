package rpc

import (
	"context"
	"errors"

	userPb "github.com/CocaineCong/Go-SearchEngine/idl/pb/user"
	"github.com/CocaineCong/Go-SearchEngine/pkg/e"
)

func UserLogin(ctx context.Context, req *userPb.UserLoginReq) (resp *userPb.UserDetailResponse, err error) {
	r, err := UserClient.UserLogin(ctx, req)
	if err != nil {
		return
	}

	if r.Code != e.SUCCESS {
		err = errors.New("登陆失败")
		return
	}

	return r, nil
}

func UserRegister(ctx context.Context, req *userPb.UserRegisterReq) (resp *userPb.UserCommonResponse, err error) {
	r, err := UserClient.UserRegister(ctx, req)
	if err != nil {
		return
	}

	if r.Code != e.SUCCESS {
		err = errors.New(r.Msg)
		return
	}

	return
}
