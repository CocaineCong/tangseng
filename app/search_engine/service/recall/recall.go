package recall

import (
	"context"

	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"

	"github.com/CocaineCong/tangseng/app/search_engine/analyzer"
	"github.com/CocaineCong/tangseng/app/search_engine/ranking"
	"github.com/CocaineCong/tangseng/app/search_engine/repository/db/dao"
	"github.com/CocaineCong/tangseng/app/search_engine/repository/storage"
	"github.com/CocaineCong/tangseng/app/search_engine/rpc"
	pb "github.com/CocaineCong/tangseng/idl/pb/search_vector"
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
func (r *Recall) Search(ctx context.Context, query string) (resp []*types.SearchItem, err error) {
	splitQuery, err := analyzer.GseCutForRecall(query)
	if err != nil {
		log.LogrusObj.Errorf("text2postingslists err: %v", err)
		return
	}

	// 倒排库搜索
	res, err := r.searchDoc(ctx, splitQuery)
	if err != nil {
		log.LogrusObj.Errorf("searchDoc err: %v", err)
		return
	}

	// 向量库搜索
	vRes, err := r.SearchVector(ctx, splitQuery)
	if err != nil {
		log.LogrusObj.Errorf("SearchVector err: %v", err)
		return
	}

	resp, _ = r.Multiplex(ctx, query, res, vRes)
	return
}

// Multiplex 多路融合排序
func (r *Recall) Multiplex(ctx context.Context, query string, iRes, vRes []int64) ([]*types.SearchItem, error) {
	// 融合去重
	iRes = append(iRes, vRes...)
	iRes = lo.Uniq(iRes)
	recallData, _ := dao.NewInputDataDao(ctx).ListInputDataByDocIds(iRes)
	// 排序
	searchItems := ranking.CalculateScoreBm25(query, recallData)

	return searchItems, nil
}

// SearchVector 搜索向量
func (r *Recall) SearchVector(ctx context.Context, queries []string) (docIds []int64, err error) {
	// rpc 调用python接口 获取
	req := &pb.SearchVectorRequest{Query: queries}
	vectorResp, err := rpc.SearchVector(ctx, req)
	if err != nil {
		log.LogrusObj.Errorln(err)
		return
	}
	docIds = make([]int64, len(vectorResp.DocIds))
	for i, v := range vectorResp.DocIds {
		docIds[i] = cast.ToInt64(v)
	}

	// 去重
	// vDocIds := lo.Uniq(vectorResp.DocIds)

	// 查询正排库
	// docIds := make([]uint32, len(vectorResp.DocIds))
	// for _, v := range vDocIds {
	// 	docIds = append(docIds, cast.ToUint32(v))
	// }
	// vList, err := dao.NewInputDataDao(ctx).ListInputDataByDocIds(docIds)
	// if err != nil {
	// 	log.LogrusObj.Errorln(err)
	// 	return
	// }

	// for _, v := range vList {
	// 	res = append(res, &types.SearchItem{
	// 		DocId:        v.DocId,
	// 		Content:      v.Content,
	// 		Title:        v.Title,
	// 		Score:        0,
	// 		DocCount:     0,
	// 		ContentScore: 0,
	// 	})
	// }

	return
}

// SearchQueryWord 入口词语联想
func (r *Recall) SearchQueryWord(query string) (resp []string, err error) {
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

func (r *Recall) searchDoc(ctx context.Context, tokens []string) (recalls []int64, err error) {
	recalls = make([]int64, 0)
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
		// TODO: term的position，后面再更新
		for _, v := range docIds.ToArray() {
			recalls = append(recalls, cast.ToInt64(v))
		}
	}

	// 排序打分
	// iDao := dao.NewInputDataDao(ctx)
	// for _, p := range allPostingsList {
	// 	if p == nil || p.DocIds == nil || p.DocIds.IsEmpty() {
	// 		continue
	// 	}
	// 	recallData, _ := iDao.ListInputDataByDocIds(p.DocIds.ToArray())
	// 	searchItems := ranking.CalculateScoreBm25(p.Term, recallData)
	// 	recalls = append(recalls, searchItems...)
	// }

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
