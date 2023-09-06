package dao

import (
	"context"
	"errors"

	"gorm.io/gorm"

	userPb "github.com/CocaineCong/tangseng/idl/pb/user"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/repository/mysql/db"
	"github.com/CocaineCong/tangseng/repository/mysql/model"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{db.NewDBClient(ctx)}
}

// GetUserInfo 获取用户信息
func (dao *UserDao) GetUserInfo(req *userPb.UserLoginReq) (r *model.User, err error) {
	err = dao.Model(&model.User{}).Where("user_name=?", req.UserName).
		First(&r).Error

	return
}

// CreateUser 用户创建
func (dao *UserDao) CreateUser(req *userPb.UserRegisterReq) (err error) {
	var user model.User
	var count int64
	dao.Model(&model.User{}).Where("user_name = ?", req.UserName).Count(&count)
	if count != 0 {
		return errors.New("UserName Exist")
	}

	user = model.User{
		UserName: req.UserName,
		NickName: req.NickName,
	}
	_ = user.SetPassword(req.Password)
	if err = dao.Model(&model.User{}).Create(&user).Error; err != nil {
		log.LogrusObj.Error("Insert User Error:" + err.Error())
		return
	}

	return
}
