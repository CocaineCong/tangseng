package loading

import (
	"github.com/CocaineCong/tangseng/pkg/es"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// Loading 全局loading
func Loading() {
	es.InitEs()
	log.InitLog()
}
