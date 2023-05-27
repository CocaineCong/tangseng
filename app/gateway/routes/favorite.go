package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/CocaineCong/Go-SearchEngine/app/gateway/internal/handler"
)

func FavoriteRegisterHandlers(rg *gin.RouterGroup) {
	favoriteGroup := rg.Group("/favorite")
	{
		favoriteGroup.POST("/create", handler.CreateFavorite)
		favoriteGroup.GET("/list", handler.ListFavorite)
		favoriteGroup.POST("/update", handler.UpdateFavorite)
		favoriteGroup.POST("/delete", handler.DeleteFavorite)
	}

	favoriteDetailGroup := rg.Group("/favorite_detail")
	{
		favoriteDetailGroup.POST("/create", handler.ListFavoriteDetail)
		favoriteDetailGroup.GET("/list", handler.CreateFavoriteDetail)
		favoriteDetailGroup.POST("/update", handler.DeleteFavoriteDetail)
	}

}
