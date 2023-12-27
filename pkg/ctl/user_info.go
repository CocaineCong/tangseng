package ctl

import (
	"context"
	"github.com/pkg/errors"

	"github.com/CocaineCong/tangseng/consts"
)

type UserInfo struct {
	Id       int64  `json:"id"`
	UserName string `json:"user_name"`
}

func GetUserInfo(ctx context.Context) (*UserInfo, error) {
	user, ok := FromContext(ctx)
	if !ok {
		return nil, errors.Wrap(errors.New("获取用户信息错误"), "FromContext error")
	}
	return user, nil
}

func NewContext(ctx context.Context, u *UserInfo) context.Context {
	return context.WithValue(ctx, consts.UserInfoKey, u) // nolint:golint,staticcheck
}

func FromContext(ctx context.Context) (*UserInfo, bool) {
	u, ok := ctx.Value(consts.UserInfoKey).(*UserInfo)
	return u, ok
}

func InitUserInfo(ctx context.Context) {
	// TOOD 放缓存，之后的用户信息，走缓存
}
