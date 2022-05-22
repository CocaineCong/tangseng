package routes

import (
	"api-gateway/internal/handler"
	"api-gateway/middleware"
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

func NewRouter(service ...interface{}) *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(middleware.Cors(), middleware.InitMiddleware(service), middleware.ErrorMiddleware())
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
		authed := v1.Group("/")
		authed.Use(middleware.JWT())
		{
			authed.GET("favorites", handler.GetFavoriteList)
			authed.POST("favorites", handler.CreateFavorite)
			authed.PUT("favorites", handler.UpdateFavorite)
			authed.DELETE("favorites", handler.DeleteFavorite)

			authed.GET("favorites-detail", handler.GetFavoriteDetail)
			authed.POST("favorites-detail", handler.CreateFavoriteDetail)
			authed.DELETE("favorites-detail", handler.DeleteFavoriteDetail)
		}
	}
	return ginRouter
}

