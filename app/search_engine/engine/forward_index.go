package engine

import (
	"github.com/CocaineCong/tangseng/app/search_engine/types"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// AddForwardIndex 落库正排索引
func (e *Engine) AddForwardIndex(doc *types.Document) error {
	return e.Seg[e.CurrSegId].AddForwardByDoc(doc)
}

func (e *Engine) UpdateCount(num int64) (err error) {
	seg := e.Seg[e.CurrSegId]
	count, err := seg.ForwardCount()
	if err != nil {
		log.LogrusObj.Errorf("UpdateCount err:%v", err)
		return
	}
	count += num

	return seg.UpdateForwardCount(count)
}
