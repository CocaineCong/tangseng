package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/CocaineCong/Go-SearchEngine/app/gateway/internal/handler"
	"github.com/CocaineCong/Go-SearchEngine/app/gateway/middleware"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(middleware.Cors(), middleware.ErrorMiddleware())
	store := cookie.NewStore([]byte("something-very-secret"))
	ginRouter.Use(sessions.Sessions("mysession", store))
	v1 := ginRouter.Group("/api/v1")
	{
		v1.GET("ping", func(context *gin.Context) {
			context.JSON(200, "success")
		})
		// 用户服务
		v1.POST("/user/register", handler.UserRegister)
		v1.POST("/user/login", handler.UserLogin)

		// 需要登录保护
		v1.POST("/add", handler.Add)
		// 搜索引擎
		v1.GET("/search", handler.Search)
		v1.GET("/allindex", handler.AllIndex)
		v1.GET("/allindexcount", handler.AllIndexCount)
		v1.GET("/search-word", handler.SearchWord)

		authed := v1.Group("/")
		authed.Use(middleware.AuthMiddleware())
		{
			// 收藏夹模块
			authed.GET("favorites", handler.GetFavoriteList)
			authed.POST("favorites", handler.CreateFavorite)
			authed.PUT("favorites", handler.UpdateFavorite)
			authed.DELETE("favorites", handler.DeleteFavorite)

			// 收藏夹详情模块
			authed.GET("favorites-detail", handler.GetFavoriteDetail)
			authed.POST("favorites-detail", handler.CreateFavoriteDetail)
			authed.DELETE("favorites-detail", handler.DeleteFavoriteDetail)
		}
	}
	return ginRouter
}
