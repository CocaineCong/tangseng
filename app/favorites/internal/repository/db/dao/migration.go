package dao

import (
	"os"

	"github.com/CocaineCong/Go-SearchEngine/app/favorites/internal/repository/db/model"
	logging "github.com/CocaineCong/Go-SearchEngine/pkg/util/logger"
)

func migration() {
	// 自动迁移模式
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.Favorites{},
			&model.FavoritesDetails{},
		)
	if err != nil {
		logging.LogrusObj.Infoln("register table fail")
		os.Exit(0)
	}
	logging.LogrusObj.Infoln("register table success")
}
