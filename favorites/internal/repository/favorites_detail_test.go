package repository

import (
	"favorites/internal/service"
	"fmt"
	"testing"
)

func TestFavoritesDetails_Create(t *testing.T) {
	InitDB()
	f := new(FavoritesDetails)
	req := new(service.FavoritesRequest)
	reqUrl := new(service.UrlModel)
	req.FavoriteID = 1
	req.UserID = 4
	reqUrl.UrlID = 3
	reqUrl.Url = "https://bilibili.com"
	reqUrl.Desc = "love your love"
	req.UrlInfo = reqUrl
	err := f.Create(req)
	fmt.Println(err)
}

func TestFavoritesDetails_Show(t *testing.T) {
	InitDB()
	f := new(FavoritesDetails)
	req := new(service.FavoritesRequest)
	req.UserID = 4
	res, _ := f.Show(req)
	fmt.Println(res)
}

func TestFavoritesDetails_Delete(t *testing.T) {
	InitDB()
	f := new(FavoritesDetails)
	req := new(service.FavoritesRequest)
	req.FavoriteDetailID=4
	req.FavoriteID=1
	err := f.Delete(req)
	fmt.Println(err)
}