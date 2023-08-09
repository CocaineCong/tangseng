package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/CocaineCong/tangseng/app/gateway/http"
)

func FavoriteRegisterHandlers(rg *gin.RouterGroup) {
	favoriteGroup := rg.Group("/favorite")
	{
		favoriteGroup.POST("/create", http.CreateFavorite)
		favoriteGroup.GET("/list", http.ListFavorite)
		favoriteGroup.POST("/update", http.UpdateFavorite)
		favoriteGroup.POST("/delete", http.DeleteFavorite)
	}

	favoriteDetailGroup := rg.Group("/favorite_detail")
	{
		favoriteDetailGroup.POST("/create", http.CreateFavoriteDetail)
		favoriteDetailGroup.GET("/list", http.ListFavoriteDetail)
		favoriteDetailGroup.POST("/delete", http.DeleteFavoriteDetail)
	}

}
