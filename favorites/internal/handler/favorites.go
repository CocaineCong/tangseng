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

func (*FavoriteService) FavoriteShow(ctx context.Context,req *service.FavoritesRequest) (resp *service.FavoritesDetailResponse,err error) {
	var f repository.Favorites
	resp = new(service.FavoritesDetailResponse)
	fRep, err := f.Show(req)
	resp.Code = e.SUCCESS
	if err != nil {
		resp.Code=e.ERROR
		return resp, err
	}
	resp.FavoritesDetail = repository.BuildFavorites(fRep)
	return resp,nil
}

func (*FavoriteService) FavoriteUpdate(ctx context.Context,req *service.FavoritesRequest) (resp *service.CommonResponse,err error) {
	var favorite repository.Favorites
	resp = new(service.CommonResponse)
	resp.Code = e.SUCCESS
	err = favorite.Update(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = e.GetMsg(e.ERROR)
		resp.Data = err.Error()
		return resp, err
	}
	resp.Msg = e.GetMsg(uint(resp.Code))
	return resp,nil
}

func (*FavoriteService) FavoriteDelete(ctx context.Context,req *service.FavoritesRequest) (resp *service.CommonResponse,err error) {
	var favorite repository.Favorites
	resp = new(service.CommonResponse)
	resp.Code = e.SUCCESS
	err = favorite.Delete(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = e.GetMsg(e.ERROR)
		resp.Data = err.Error()
		return resp, err
	}
	resp.Msg = e.GetMsg(uint(resp.Code))
	return resp,nil
}

func (*FavoriteService) FavoriteDetailsCreate(ctx context.Context,req *service.FavoritesRequest) (resp *service.CommonResponse,err error) {
	var favorite repository.FavoritesDetails
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

func (*FavoriteService) FavoriteDetailsDelete(ctx context.Context,req *service.FavoritesRequest) (resp *service.CommonResponse,err error) {
	var favorite repository.FavoritesDetails
	resp = new(service.CommonResponse)
	resp.Code = e.SUCCESS
	err = favorite.Delete(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = e.GetMsg(e.ERROR)
		resp.Data = err.Error()
		return resp, err
	}
	resp.Msg = e.GetMsg(uint(resp.Code))
	return resp, nil
}

func (*FavoriteService) FavoriteDetailsShow (ctx context.Context,req *service.FavoritesRequest) (resp *service.FavoritesDetailResponse,err error) {
	var favorite repository.FavoritesDetails
	resp = new(service.FavoritesDetailResponse)
	resp.Code = e.SUCCESS
	fdResp ,err := favorite.Show(req)
	if err != nil {
		resp.Code = e.ERROR
		return resp, err
	}
	resp.FavoritesDetail = repository.BuildFavorites(fdResp)
	return resp, nil
}