package recall

import (
	"context"

	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/types"
)

// SearchRecall 词条回归
func SearchRecall(ctx context.Context, query string) (res []*types.SearchItem, err error) {
	recallService := NewRecall()
	res, err = recallService.Search(ctx, query)
	if err != nil {
		log.LogrusObj.Errorf("SearchRecall-NewRecallServ:%+v", err)
		return
	}

	return
}

// SearchQuery 词条联想
func SearchQuery(query string) (res []string, err error) {
	recallService := NewRecall()
	res, err = recallService.SearchQuery(query)
	if err != nil {
		log.LogrusObj.Errorf("SearchRecall-NewRecallServ:%+v", err)
		return
	}

	return
}
