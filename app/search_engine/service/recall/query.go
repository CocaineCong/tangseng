package recall

import (
	"context"

	"github.com/pkg/errors"

	"github.com/CocaineCong/tangseng/types"
)

// SearchRecall 词条回归
func SearchRecall(ctx context.Context, query string) (res []*types.SearchItem, err error) {
	recallService := NewRecall()
	res, err = recallService.Search(ctx, query)
	if err != nil {
		err = errors.WithMessage(err, "SearchRecall-NewRecallServ error")
		return
	}

	return
}

// SearchQuery 词条联想
func SearchQuery(query string) (res []string, err error) {
	recallService := NewRecall()
	res, err = recallService.SearchQueryWord(query)
	if err != nil {
		err = errors.WithMessage(err, "SearchRecall-NewRecallServ error")
		return
	}

	return
}
