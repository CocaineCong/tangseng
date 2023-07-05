package service

import (
	"context"
	"sync"

	"github.com/CocaineCong/tangseng/app/user/internal/repository/db/dao"
	e2 "github.com/CocaineCong/tangseng/consts/e"
	pb "github.com/CocaineCong/tangseng/idl/pb/user"
)

var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

type UserSrv struct {
	pb.UnimplementedUserServiceServer
}

func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

func (u *UserSrv) UserLogin(ctx context.Context, req *pb.UserLoginReq) (resp *pb.UserDetailResponse, err error) {
	resp = new(pb.UserDetailResponse)
	resp.Code = e2.SUCCESS
	r, err := dao.NewUserDao(ctx).GetUserInfo(req)
	if err != nil {
		resp.Code = e2.ERROR
		return
	}
	resp.UserDetail = &pb.UserResp{
		UserId:   r.UserID,
		UserName: r.UserName,
		NickName: r.NickName,
	}
	return
}

func (u *UserSrv) UserRegister(ctx context.Context, req *pb.UserRegisterReq) (resp *pb.UserCommonResponse, err error) {
	resp = new(pb.UserCommonResponse)
	resp.Code = e2.SUCCESS
	err = dao.NewUserDao(ctx).CreateUser(req)
	if err != nil {
		resp.Code = e2.ERROR
		return
	}
	resp.Data = e2.GetMsg(int(resp.Code))
	return
}
