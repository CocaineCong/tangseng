package handler

import (
	"api-gateway/internal/service"
	"api-gateway/pkg/e"
	"api-gateway/pkg/res"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Add(ginCtx *gin.Context) {
	var seReq service.SearchEngineRequest
	PanicIfSearchEngineError(ginCtx.ShouldBind(&seReq))
	// 从gin.Key中取出服务实例
	searchEngineService := ginCtx.Keys["se"].(service.SearchEngineServiceClient)
	searchEngineResp, err := searchEngineService.SearchEngineAdd(context.Background(), &seReq)
	PanicIfSearchEngineError(err)
	r := res.Response{
		Data:   searchEngineResp,
		Status: uint(searchEngineResp.Code),
		Msg:    e.GetMsg(uint(searchEngineResp.Code)),
	}
	ginCtx.JSON(http.StatusOK, r)
}

func Search(ginCtx *gin.Context) {
	var seReq service.SearchEngineRequest
	PanicIfSearchEngineError(ginCtx.ShouldBind(&seReq))
	// 从gin.Key中取出服务实例
	searchEngineService := ginCtx.Keys["se"].(service.SearchEngineServiceClient)
	searchEngineResp, err := searchEngineService.SearchEngineSearch(context.Background(), &seReq)
	PanicIfSearchEngineError(err)
	r := res.Response{
		Data:   searchEngineResp,
		Status: uint(searchEngineResp.Code),
		Msg:    e.GetMsg(uint(searchEngineResp.Code)),
	}
	ginCtx.JSON(http.StatusOK, r)
}