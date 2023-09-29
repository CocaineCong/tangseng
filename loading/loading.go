package loading

import (
	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/pkg/kfk"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/repository/mysql/db"
)

// Loading 全局loading
func Loading() {
	// es.InitEs()
	config.InitConfig()
	log.InitLog()
	config.InitConfig()

	db.InitDB()
	kfk.InitKafka()
	// dao.InitMysqlDirectUpload(ctx)
}
