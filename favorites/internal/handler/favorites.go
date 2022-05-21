package handler

import (
	"context"
	"favorites/internal/repository"
	"favorites/internal/service"
	"favorites/pkg/e"
)

type FavoriteService struct {

}

func NewFavoriteService() *FavoriteService {
	return &FavoriteService{}
}

func (*FavoriteService) FavoriteCreate(ctx context.Context,req *service.FavoritesRequest) (resp *service.CommonResponse,err error) {
	var favorite repository.Favorites
	resp = new(service.CommonResponse)
	resp.Code = e.SUCCESS
	err = favorite.Create(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = e.GetMsg(e.ERROR)
		resp.Data = err.Error()
		return resp, err
	}
	resp.Msg = e.GetMsg(uint(resp.Code))
	return resp,nil
}
