package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/CocaineCong/tangseng/app/gateway/http"
	"github.com/CocaineCong/tangseng/app/gateway/middleware"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors(), middleware.ErrorMiddleware())
	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(sessions.Sessions("mysession", store))
	v1 := r.Group("/api/v1")
	{
		v1.GET("ping", func(context *gin.Context) {
			context.JSON(200, "success")
		})
		// 用户服务
		v1.POST("/user/register", http.UserRegister)
		v1.POST("/user/login", http.UserLogin)
		// 索引平台
		IndexPlatformRegisterHandlers(v1)
		// 搜索平台
		SearchRegisterHandlers(v1)
		// 需要登录保护
		authed := v1.Group("/")
		authed.Use(middleware.AuthMiddleware())
		{
			// 收藏夹模块
			FavoriteRegisterHandlers(authed)
		}
	}
	return r
}
