package recall

import (
	"errors"

	"github.com/CocaineCong/tangseng/app/search_engine/analyzer"
	"github.com/CocaineCong/tangseng/app/search_engine/engine"
	"github.com/CocaineCong/tangseng/app/search_engine/ranking"
	"github.com/CocaineCong/tangseng/app/search_engine/segment"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/types"
)

// Recall 查询召回
type Recall struct {
	*engine.Engine
	docCount     int64 // 文档总数 ，用于计算相关性
	enablePhrase bool
}

// NewRecall --
func NewRecall(meta *engine.Meta) *Recall {
	e := engine.NewTangSengEngine(meta, segment.SearchMode)
	var docCount int64 = 0
	for _, seg := range e.Seg {
		num, err := seg.ForwardCount()
		if err != nil {
			log.LogrusObj.Errorf("error:%v", err)
		}
		docCount += num
	}

	return &Recall{e, docCount, true}
}

// Search 入口
func (r *Recall) Search(query string) (res []*types.SearchItem, err error) {
	splitQuery, err := r.splitQuery2Tokens(query)
	if err != nil {
		log.LogrusObj.Errorf("splitQuery2Tokens err:%v", err)
		return nil, err
	}

	res, err = r.searchDoc(splitQuery)

	return
}

// SearchQuery 入口
func (r *Recall) SearchQuery(query string) ([]*types.DictTireTree, error) {
	return r.GetDict(query)
}

func (r *Recall) splitQuery2Tokens(query string) (res []string, err error) {
	res, err = analyzer.GseCutForRecall(query)
	if err != nil {
		log.LogrusObj.Errorf("text2postingslists err: %v", err)
		return
	}

	return
}

func (r *Recall) searchDoc(tokens []string) (recalls []*types.SearchItem, err error) {
	recalls = make([]*types.SearchItem, 0)
	exist := make(map[int64]struct{})
	for _, token := range tokens { // TODO: 重复的值结果缓存起来
		if token == "" {
			err = errors.New("token is nil1")
			return
		}
		postings, count, errx := fetchPostingsBySegs(token)
		if errx != nil {
			err = errx
			return
		}
		if postings == nil {
			return
		}
		log.LogrusObj.Infof("token:%s,incvertedIndex:%d", token, postings.DocId)
		post.DocCount = count
		for postings != nil {
			docId := postings.DocId
			if postings.DocId == 0 {
				postings = postings.Next
				continue
			}
			sItem := &types.SearchItem{
				DocId:    docId,
				Content:  "",
				Score:    0.0,
				DocCount: post.DocCount,
			}
			sItem, err = r.getContentByDocId(sItem)
			if err != nil {
				log.LogrusObj.Errorf("getContentByDocId:%d, err:%v", docId, err)
				return
			}
			if _, ok := exist[sItem.DocId]; !ok {
				recalls = append(recalls, sItem)
			}
			postings = postings.Next
		}

		recalls = ranking.CalculateScoreBm25(token, recalls)
	}

	log.LogrusObj.Infof("recalls size:%v", len(recalls))

	return
}

// 获取 token 所有seg的倒排表数据
func fetchPostingsBySegs(token string) (postings *types.PostingsList, docCount int64, err error) {
	// 遍历存储index的地方，将所有的token都存储一次
	return
}

// 获取 token 所有seg的倒排表数据
func (r *Recall) fetchPostingsBySegs(token string) (postings *types.PostingsList, docCount int64, err error) {
	postings = new(types.PostingsList)
	for i, seg := range r.Engine.Seg {
		p, errx := seg.FetchPostings(token)
		if errx != nil {
			err = errx
			log.LogrusObj.Errorf("seg.FetchPostings index:%v", i)
			return
		}
		log.LogrusObj.Infof("post:%v", p)
		postings = segment.MergePostings(postings, p.PostingsList)
		log.LogrusObj.Infof("pos next:%v", postings.Next)
		docCount += p.DocCount
	}
	log.LogrusObj.Infof("token:%v,pos:%v,doc:%v", token, postings, docCount)

	return
}

func (r *Recall) getContentByDocId(s *types.SearchItem) (item *types.SearchItem, err error) {
	for i, seg := range r.Engine.Seg {
		p, errx := seg.GetForward(s.DocId)
		if errx != nil {
			err = errx
			log.LogrusObj.Errorf("seg.FetchPostings index:%v", i)
			return
		}
		s.Content = string(p)
	}
	item = new(types.SearchItem)
	item = s

	return
}
