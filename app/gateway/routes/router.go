package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/CocaineCong/Go-SearchEngine/app/gateway/internal/handler"
	"github.com/CocaineCong/Go-SearchEngine/app/gateway/middleware"
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
		v1.POST("/user/register", handler.UserRegister)
		v1.POST("/user/login", handler.UserLogin)

		// v1.POST("/add", handler.Add)
		// // 搜索引擎
		// v1.GET("/search", handler.Search)
		// v1.GET("/allindex", handler.AllIndex)
		// v1.GET("/allindexcount", handler.AllIndexCount)
		// v1.GET("/search-word", handler.SearchWord)

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
