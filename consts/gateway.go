package consts

import (
	"time"
)

const (
	AccessTokenHeader          = "access_token"
	RefreshTokenHeader         = "refresh_token"
	HeaderForwardedProto       = "X-Forwarded-Proto"
	MaxAge                     = 3600 * 24
	AccessTokenExpireDuration  = 24 * time.Hour
	RefreshTokenExpireDuration = 10 * 24 * time.Hour
)

const UserInfoKey = "user_info_key"
