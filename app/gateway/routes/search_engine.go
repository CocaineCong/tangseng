package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/CocaineCong/tangseng/app/gateway/http"
)

func SearchRegisterHandlers(rg *gin.RouterGroup) {
	searchEngineGroup := rg.Group("/search_engine")
	{
		searchEngineGroup.GET("/search", http.SearchEngineSearch)
		searchEngineGroup.GET("/query", http.WordAssociation)
	}

	searchVectorGroup := rg.Group("/search_vector")
	{
		searchVectorGroup.GET("/vector", http.SearchVector)
	}
}
