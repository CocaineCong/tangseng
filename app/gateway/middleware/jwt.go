package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/CocaineCong/tangseng/consts"
	e2 "github.com/CocaineCong/tangseng/consts/e"
	"github.com/CocaineCong/tangseng/pkg/ctl"
	"github.com/CocaineCong/tangseng/pkg/jwt"
)

// AuthMiddleware token验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = e2.SUCCESS
		accessToken := c.GetHeader("access_token")
		refreshToken := c.GetHeader("refresh_token")
		if accessToken == "" {
			code = e2.InvalidParams
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e2.GetMsg(code),
				"data":   "Token不能为空",
			})
			c.Abort()
			return
		}
		newAccessToken, newRefreshToken, err := jwt.ParseRefreshToken(accessToken, refreshToken)
		if err != nil {
			code = e2.ErrorAuthCheckTokenFail
		}
		if code != e2.SUCCESS {
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e2.GetMsg(code),
				"data":   "鉴权失败",
				"error":  err.Error(),
			})
			c.Abort()
			return
		}
		claims, err := jwt.ParseToken(newAccessToken)
		if err != nil {
			code = e2.ErrorAuthCheckTokenFail
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e2.GetMsg(code),
				"data":   err.Error(),
			})
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(), &ctl.UserInfo{Id: claims.ID, UserName: claims.Username}))
		SetToken(c, newAccessToken, newRefreshToken)
		ctl.InitUserInfo(c.Request.Context())
		c.Next()
	}
}

// SetToken 设置token
func SetToken(c *gin.Context, accessToken, refreshToken string) {
	secure := IsHttps(c)
	c.Header(consts.AccessTokenHeader, accessToken)
	c.Header(consts.RefreshTokenHeader, refreshToken)
	c.SetCookie(consts.AccessTokenHeader, accessToken, consts.MaxAge, "/", "", secure, true)
	c.SetCookie(consts.RefreshTokenHeader, refreshToken, consts.MaxAge, "/", "", secure, true)
}

// IsHttps 判断是否https
func IsHttps(c *gin.Context) bool {
	if c.GetHeader(consts.HeaderForwardedProto) == "https" || c.Request.TLS != nil {
		return true
	}
	return false
}
