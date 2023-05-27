package loading

import (
	"github.com/CocaineCong/Go-SearchEngine/pkg/es"
	log "github.com/CocaineCong/Go-SearchEngine/pkg/logger"
)

// Loading 全局loading
func Loading() {
	es.InitEs()
	log.InitLog()
}
