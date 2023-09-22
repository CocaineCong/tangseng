package recall

import (
	"context"

	"github.com/RoaringBitmap/roaring"

	"github.com/CocaineCong/tangseng/app/search_engine/analyzer"
	"github.com/CocaineCong/tangseng/app/search_engine/ranking"
	"github.com/CocaineCong/tangseng/app/search_engine/repository/db/dao"
	"github.com/CocaineCong/tangseng/app/search_engine/repository/storage"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/types"
)

// Recall 查询召回
type Recall struct {
}

func NewRecall() *Recall {
	return &Recall{}
}

// Search 入口
func (r *Recall) Search(ctx context.Context, query string) (res []*types.SearchItem, err error) {
	splitQuery, err := analyzer.GseCutForRecall(query)
	if err != nil {
		log.LogrusObj.Errorf("text2postingslists err: %v", err)
		return
	}

	res, err = r.searchDoc(ctx, splitQuery)

	return
}

// SearchQuery 入口
func (r *Recall) SearchQuery(query string) ([]*types.DictTireTree, error) {
	// return r.GetDict(query)
	return nil, nil
}

func (r *Recall) searchDoc(ctx context.Context, tokens []string) (recalls []*types.SearchItem, err error) {
	recalls = make([]*types.SearchItem, 0)
	// exist := make(map[int64]struct{}) // TODO redis 存放已经搜索过的 token
	allPostingsList := []*types.PostingsList{}
	for _, token := range tokens {
		postingsList, errx := fetchPostingsByToken(token)
		if errx != nil {
			log.LogrusObj.Errorln(errx)
			continue
		}
		allPostingsList = append(allPostingsList, postingsList...)
	}

	iDao := dao.NewInputDataDao(ctx)
	for _, p := range allPostingsList {
		if p.DocIds.IsEmpty() {
			continue
		}
		recallData, _ := iDao.ListInputDataByDocIds(p.DocIds.ToArray())
		searchItems := ranking.CalculateScoreBm25(p.Term, recallData)

		recalls = append(recalls, searchItems...)
	}

	log.LogrusObj.Infof("recalls size:%v", len(recalls))

	return
}

// 获取 token 所有seg的倒排表数据
func fetchPostingsByToken(token string) (postingsList []*types.PostingsList, err error) {
	// 遍历存储index的地方，token对应的doc Id 全部取出
	output := roaring.New()
	postingsList = make([]*types.PostingsList, 0, 1e6)
	for _, inverted := range storage.GobalInvertedDB {
		docIds, errx := inverted.GetInverted([]byte(token))
		if errx != nil {
			log.LogrusObj.Errorln(errx)
			continue
		}
		output = roaring.New()
		_ = output.UnmarshalBinary(docIds)
		// 存放到数组当中
		postings := &types.PostingsList{
			Term:   token,
			DocIds: output,
		}
		postingsList = append(postingsList, postings)
	}

	return
}
