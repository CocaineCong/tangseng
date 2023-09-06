package db

import (
	"os"

	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/repository/mysql/model"
)

func migration() {
	// 自动迁移模式
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.User{},
			&model.InputData{},
			&model.Favorite{},
			&model.FavoriteDetail{},
		)
	if err != nil {
		log.LogrusObj.Infoln("register table fail")
		os.Exit(0)
	}
	log.LogrusObj.Infoln("register table success")
}
