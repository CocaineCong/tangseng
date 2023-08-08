package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/CocaineCong/tangseng/app/gateway/internal/handler"
)

func SearchRegisterHandlers(rg *gin.RouterGroup) {
	favoriteGroup := rg.Group("/search_engine")
	{
		favoriteGroup.POST("/search", handler.SearchEngineSearch)
	}
}
