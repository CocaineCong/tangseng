package loading

import (
	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/pkg/kfk"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/repository/mysql/db"
	"github.com/CocaineCong/tangseng/repository/redis"
)

// Loading 全局loading
func Loading() {
	// es.InitEs()
	config.InitConfig()
	log.InitLog()

	db.InitDB()
	redis.InitRedis()
	kfk.InitKafka()
	// dao.InitMysqlDirectUpload(ctx)
}
