package handler

import (
	"api-gateway/pkg/util"
	"errors"
)

// 包装错误
func PanicIfUserError(err error) {
	if err != nil {
		err = errors.New("userService--" + err.Error())
		util.LogrusObj.Info(err)
		panic(err)
	}
}

func PanicIfFavoriteError(err error) {
	if err != nil {
		err = errors.New("favoriteService--" + err.Error())
		util.LogrusObj.Info(err)
		panic(err)
	}
}

func PanicIfSearchEngineError(err error) {
	if err != nil {
		err = errors.New("searchEngineService--" + err.Error())
		util.LogrusObj.Info(err)
		panic(err)
	}
}