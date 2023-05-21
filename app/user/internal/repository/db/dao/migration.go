package dao

import (
	"os"

	"github.com/CocaineCong/Go-SearchEngine/app/user/internal/repository/db/model"
	"github.com/CocaineCong/Go-SearchEngine/pkg/util/logger"
)

// 自动迁移模式
func migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.User{},
		)
	if err != nil {
		logger.LogrusObj.Infoln("register table fail")
		os.Exit(0)
	}
	logger.LogrusObj.Infoln("register table success")
}
