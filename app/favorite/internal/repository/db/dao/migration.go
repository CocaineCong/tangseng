package dao

import (
	"os"

	logging "github.com/CocaineCong/Go-SearchEngine/pkg/util/logger"

	"github.com/CocaineCong/Go-SearchEngine/app/favorite/internal/repository/db/model"
)

func migration() {
	// 自动迁移模式
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.Favorite{},
			&model.FavoriteDetail{},
		)
	if err != nil {
		logging.LogrusObj.Infoln("register table fail")
		os.Exit(0)
	}
	logging.LogrusObj.Infoln("register table success")
}
