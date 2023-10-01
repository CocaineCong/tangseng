package recall

import (
	"context"

	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"

	"github.com/CocaineCong/tangseng/app/search_engine/analyzer"
	"github.com/CocaineCong/tangseng/app/search_engine/ranking"
	"github.com/CocaineCong/tangseng/app/search_engine/repository/db/dao"
	"github.com/CocaineCong/tangseng/app/search_engine/repository/storage"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/repository/redis"
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
func (r *Recall) SearchQuery(query string) (resp []string, err error) {
	dictTreeList := make([]string, 0, 1e3)
	for _, trieDb := range storage.GlobalTrieDB {
		trie, errx := trieDb.GetTrieTreeDict()
		if errx != nil {
			log.LogrusObj.Errorln(errx)
			continue
		}
		queryTrie := trie.FindAllByPrefixForRecall(query)
		dictTreeList = append(dictTreeList, queryTrie...)
	}

	resp = lo.Uniq(dictTreeList)
	return
}

func (r *Recall) searchDoc(ctx context.Context, tokens []string) (recalls []*types.SearchItem, err error) {
	recalls = make([]*types.SearchItem, 0)
	allPostingsList := []*types.PostingsList{}
	for _, token := range tokens {
		docIds, errx := redis.GetInvertedIndexTokenDocIds(ctx, token)
		var postingsList []*types.PostingsList
		if errx != nil || docIds == nil {
			// 如果缓存不存在，就去索引表里面读取
			postingsList, err = fetchPostingsByToken(token)
			if err != nil {
				log.LogrusObj.Errorln(err)
				continue
			} else {
				// 如果缓存存在，就直接读缓存，不用担心实时性问题，缓存10分钟清空一次，这延迟是能接受到
				postingsList = append(postingsList, &types.PostingsList{
					Term:   token,
					DocIds: docIds,
				})
			}
		}
		allPostingsList = append(allPostingsList, postingsList...)
	}

	// 排序打分
	iDao := dao.NewInputDataDao(ctx)
	for _, p := range allPostingsList {
		if p == nil || p.DocIds == nil || p.DocIds.IsEmpty() {
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
	postingsList = make([]*types.PostingsList, 0, 1e6)
	for _, inverted := range storage.GlobalInvertedDB {
		docIds, errx := inverted.GetInverted([]byte(token))
		if errx != nil {
			log.LogrusObj.Errorln(errx)
			continue
		}
		output := roaring.New()
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
