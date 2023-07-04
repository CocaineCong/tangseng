package dao

import (
	"os"

	"github.com/CocaineCong/tangseng/app/favorite/internal/repository/db/model"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func migration() {
	// 自动迁移模式
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.Favorite{},
			&model.FavoriteDetail{},
		)
	if err != nil {
		log.LogrusObj.Infoln("register table fail")
		os.Exit(0)
	}
	log.LogrusObj.Infoln("register table success")
}
