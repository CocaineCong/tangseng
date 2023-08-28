package recall

import (
	"errors"
	"sort"

	"github.com/CocaineCong/tangseng/app/search_engine/engine"
	"github.com/CocaineCong/tangseng/app/search_engine/segment"
	"github.com/CocaineCong/tangseng/app/search_engine/types"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/util/relevant"
)

// Recall 查询召回
type Recall struct {
	*engine.Engine
	docCount     int64 // 文档总数 ，用于计算相关性
	enablePhrase bool
}

// NewRecall --
func NewRecall(meta *engine.Meta) *Recall {
	e := engine.NewEngine(meta, segment.SearchMode)
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
func (r *Recall) Search(query string) ([]*types.SearchItem, error) {
	err := r.splitQuery2Tokens(query)
	if err != nil {
		log.LogrusObj.Errorf("splitQuery2Tokens err:%v", err)
		return nil, err
	}

	return r.searchDoc()
}

// SearchQuery 入口
func (r *Recall) SearchQuery(query string) ([]*types.DictTireTree, error) {
	return r.GetDict(query)
}

func (r *Recall) splitQuery2Tokens(query string) (err error) {
	err = r.Text2PostingsLists(query, 0)
	if err != nil {
		log.LogrusObj.Errorf("text2postingslists err: %v", err)
		return
	}

	return
}

func (r *Recall) searchDoc() (recalls []*types.SearchItem, err error) {
	recalls = make([]*types.SearchItem, 0)

	// 为每个token初始化游标
	for token, post := range r.PostingsHashBuf {
		// 正常不会出现
		if token == "" {
			err = errors.New("token is nil1")
			return
		}
		postings, count, errx := r.fetchPostingsBySegs(token)
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
			recalls = append(recalls, sItem)
			postings = postings.Next
		}

		recalls = r.calculateScore(token, recalls)
	}

	log.LogrusObj.Infof("recalls size:%v", len(recalls))

	return
}

// calculateScore 计算相关性
func (r *Recall) calculateScore(token string, searchItem []*types.SearchItem) (resp []*types.SearchItem) {
	recallToken := make([]string, 0)

	for i := range searchItem {
		recallToken = append(recallToken, searchItem[i].Content)
	}
	corpus, _ := relevant.MakeCorpus(recallToken)
	docs := relevant.MakeDocuments(recallToken, corpus)
	tf := relevant.New()

	for _, doc := range docs {
		tf.Add(doc)
	}
	tf.CalculateIDF()
	tokenRecall := relevant.Doc{corpus[token]}
	bm25Scores := relevant.BM25(tf, tokenRecall, docs, 1.5, 0.75)
	sort.Sort(sort.Reverse(bm25Scores))

	for i := range bm25Scores {
		searchItem[bm25Scores[i].ID].Score = bm25Scores[i].Score
	}
	sort.Slice(searchItem, func(i, j int) bool {
		return searchItem[i].Score > searchItem[j].Score
	})
	resp = make([]*types.SearchItem, 0)
	resp = searchItem

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
