package handler

import (
	"api-gateway/internal/service"
	"api-gateway/pkg/e"
	"api-gateway/pkg/res"
	"api-gateway/pkg/util"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetFavoriteList(ginCtx *gin.Context) {
	var fReq service.FavoritesRequest
	PanicIfFavoriteError(ginCtx.Bind(&fReq))
	claim, _ := util.ParseToken(ginCtx.GetHeader("Authorization"))
	fReq.UserID = uint32(claim.UserID)
	favoriteService := ginCtx.Keys["favorites"].(service.FavoritesServiceClient)
	fmt.Println(fReq)
	favoriteResp, err := favoriteService.FavoriteShow(context.Background(), &fReq)
	PanicIfFavoriteError(err)
	r := res.Response{
		Data:   favoriteResp,
		Status: uint(favoriteResp.Code),
		Msg:    e.GetMsg(uint(favoriteResp.Code)),
	}
	ginCtx.JSON(http.StatusOK, r)
}

func CreateFavorite(ginCtx *gin.Context) {
	var fReq service.FavoritesRequest
	PanicIfFavoriteError(ginCtx.Bind(&fReq))
	claim, _ := util.ParseToken(ginCtx.GetHeader("Authorization"))
	fReq.UserID = uint32(claim.UserID)
	favoriteService := ginCtx.Keys["favorites"].(service.FavoritesServiceClient)
	favoriteResp, err := favoriteService.FavoriteCreate(context.Background(), &fReq)
	PanicIfFavoriteError(err)
	r := res.Response{
		Data:   favoriteResp,
		Status: uint(favoriteResp.Code),
		Msg:    e.GetMsg(uint(favoriteResp.Code)),
	}
	ginCtx.JSON(http.StatusOK, r)
}

func UpdateFavorite(ginCtx *gin.Context) {
	var fReq service.FavoritesRequest
	PanicIfFavoriteError(ginCtx.Bind(&fReq))
	claim, _ := util.ParseToken(ginCtx.GetHeader("Authorization"))
	fReq.UserID = uint32(claim.UserID)
	favoriteService := ginCtx.Keys["favorites"].(service.FavoritesServiceClient)
	favoriteResp, err := favoriteService.FavoriteUpdate(context.Background(), &fReq)
	PanicIfFavoriteError(err)
	r := res.Response{
		Data:   favoriteResp,
		Status: uint(favoriteResp.Code),
		Msg:    e.GetMsg(uint(favoriteResp.Code)),
	}
	ginCtx.JSON(http.StatusOK, r)
}


func DeleteFavorite(ginCtx *gin.Context) {
	var fReq service.FavoritesRequest
	PanicIfFavoriteError(ginCtx.Bind(&fReq))
	claim, _ := util.ParseToken(ginCtx.GetHeader("Authorization"))
	fReq.UserID = uint32(claim.UserID)
	favoriteService := ginCtx.Keys["favorites"].(service.FavoritesServiceClient)
	favoriteResp, err := favoriteService.FavoriteDelete(context.Background(), &fReq)
	PanicIfFavoriteError(err)
	r := res.Response{
		Data:   favoriteResp,
		Status: uint(favoriteResp.Code),
		Msg:    e.GetMsg(uint(favoriteResp.Code)),
	}
	ginCtx.JSON(http.StatusOK, r)
}

func GetFavoriteDetail(ginCtx *gin.Context) {
	var fReq service.FavoritesRequest
	PanicIfFavoriteError(ginCtx.Bind(&fReq))
	claim, _ := util.ParseToken(ginCtx.GetHeader("Authorization"))
	fReq.UserID = uint32(claim.UserID)
	favoriteService := ginCtx.Keys["favorites"].(service.FavoritesServiceClient)
	favoriteResp, err := favoriteService.FavoriteDetailsShow(context.Background(), &fReq)
	PanicIfFavoriteError(err)
	r := res.Response{
		Data:   favoriteResp,
		Status: uint(favoriteResp.Code),
		Msg:    e.GetMsg(uint(favoriteResp.Code)),
	}
	ginCtx.JSON(http.StatusOK, r)
}

func CreateFavoriteDetail(ginCtx *gin.Context) {
	var fReq service.FavoritesRequest
	PanicIfFavoriteError(ginCtx.Bind(&fReq))
	claim, _ := util.ParseToken(ginCtx.GetHeader("Authorization"))
	fReq.UserID = uint32(claim.UserID)
	favoriteService := ginCtx.Keys["favorites"].(service.FavoritesServiceClient)
	favoriteResp, err := favoriteService.FavoriteDetailsCreate(context.Background(), &fReq)
	PanicIfFavoriteError(err)
	r := res.Response{
		Data:   favoriteResp,
		Status: uint(favoriteResp.Code),
		Msg:    e.GetMsg(uint(favoriteResp.Code)),
	}
	ginCtx.JSON(http.StatusOK, r)
}

func DeleteFavoriteDetail(ginCtx *gin.Context) {
	var fReq service.FavoritesRequest
	PanicIfFavoriteError(ginCtx.Bind(&fReq))
	claim, _ := util.ParseToken(ginCtx.GetHeader("Authorization"))
	fReq.UserID = uint32(claim.UserID)
	favoriteService := ginCtx.Keys["favorites"].(service.FavoritesServiceClient)
	favoriteResp, err := favoriteService.FavoriteDetailsDelete(context.Background(), &fReq)
	PanicIfFavoriteError(err)
	r := res.Response{
		Data:   favoriteResp,
		Status: uint(favoriteResp.Code),
		Msg:    e.GetMsg(uint(favoriteResp.Code)),
	}
	ginCtx.JSON(http.StatusOK, r)
}