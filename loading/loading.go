package loading

import (
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/repository/mysql/db"
)

// Loading 全局loading
func Loading() {
	// es.InitEs()
	config.InitConfig()
	log.InitLog()
	db.InitDB()
}
