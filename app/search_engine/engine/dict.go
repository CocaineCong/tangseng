package engine

import (
	"github.com/CocaineCong/tangseng/app/search_engine/types"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// GetDict 获取dict
func (e *Engine) GetDict(query string) (res []*types.DictTireTree, err error) {
	trieTree, err := e.Seg[e.CurrSegId+1].GetTrieTreeDict()
	if err != nil {
		return
	}
	res = make([]*types.DictTireTree, 0)
	r := trieTree.FindAllByPrefix(query)
	for i := range r {
		// 计算相关性
		res = append(res, &types.DictTireTree{Value: r[i]})
	}

	return
}

// FlushDict 刷新dict
func (e *Engine) FlushDict(isEnd ...bool) (err error) {

	err = e.Seg[e.CurrSegId].FlushTokenDict(e.TrieTree)
	if err != nil {
		log.LogrusObj.Errorln("Flush", err)
		return
	}

	if len(isEnd) > 0 && isEnd[0] {
		return
	}

	return
}
