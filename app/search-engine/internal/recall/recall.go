package recall

import (
	"fmt"
	"sort"

	"github.com/CocaineCong/Go-SearchEngine/app/search-engine/internal/engine"
	"github.com/CocaineCong/Go-SearchEngine/app/search-engine/internal/segment"
	log "github.com/CocaineCong/Go-SearchEngine/pkg/logger"
)

// Recall 查询召回
type Recall struct {
	*engine.Engine
	docCount     int64 // 文档总数 ，用于计算相关性
	enablePhrase bool
}

// 用于实现排序的map
type queryTokenHash struct {
	token         string
	invertedIndex *segment.InvertedIndexValue
	fetchPostings *segment.PostingsList
}

// SearchItem 查询结果
type SearchItem struct {
	DocId int64
	Score float64
}

// Recalls 召回结果
type Recalls []*SearchItem

// token游标 标识当前位置
type searchCursor struct {
	doc     *segment.PostingsList // 文档编号的序列
	current *segment.PostingsList // 当前文档编号
}

// 短语游标
type phraseCursor struct {
	positions []int64 // 位置信息
	base      int64   // 词元在查询中的位置
	current   *int64  // 当前的位置信息
	index     int     // 当前位置index
}

// Search 入口
func (r *Recall) Search(query string) (Recalls, error) {
	err := r.splitQuery2Tokens(query)
	if err != nil {
		log.LogrusObj.Errorf("splitQuery2Tokens err:%v", err)
		return nil, err
	}

	recall, err := r.searchDoc()
	if err != nil {
		log.LogrusObj.Errorf("searchDoc err:%v", err)
		return nil, err
	}

	return recall, nil
}

func (r *Recall) splitQuery2Tokens(query string) error {
	err := r.Text2PostingsLists(query, 0)
	if err != nil {
		return fmt.Errorf("text2postingslists err: %v", err)
	}
	return nil
}

func (r *Recall) searchDoc() (Recalls, error) {
	recalls := make(Recalls, 0)
	tokens := make([]*queryTokenHash, 0)

	// 为每个token初始化游标
	for token, post := range r.PostingsHashBuf {
		// 正常不会出现
		if token == "" {
			return nil, fmt.Errorf("token is nil1")
		}
		postings, count, err := r.fetchPostingsBySegs(token)
		if err != nil {
			return nil, err
		}
		if postings == nil {
			return nil, err
		}
		log.LogrusObj.Infof("token:%s,incvertedIndex:%s", token, postings.DocId)
		post.DocCount = count
		t := &queryTokenHash{
			token:         token,
			invertedIndex: post,
			fetchPostings: postings,
		}
		tokens = append(tokens, t)
	}

	tokens = r.s
}

// token 根据 doc count 升序排序
func (r *Recall) sortToken(tokens []*queryTokenHash) []*queryTokenHash {
	for _, t := range tokens {
		log.LogrusObj.Infof("token:%v,docCount:%v", t.token, t.invertedIndex.DocCount)
	}
	sort.Sort()
}

// 获取 token 所有seg的倒排表数据
func (r *Recall) fetchPostingsBySegs(token string) (*segment.PostingsList, int64, error) {
	postings := &segment.PostingsList{}
	postings = nil
	docCount := int64(0)
	for i, seg := range r.Engine.Seg {
		p, c, err := seg.FetchPostings(token)
		if err != nil {
			log.LogrusObj.Errorf("seg.FetchPostings index:%v", i)
			return nil, 0, err
		}
		log.LogrusObj.Infof("post:%v", p)
		postings = segment.MergePostings(postings, p)
		log.LogrusObj.Infof("pos next:%v", postings.Next)
		docCount += c
	}
	log.LogrusObj.Infof("token:%v,pos:%v,doc:%v", token, postings, docCount)

	return postings, docCount, nil
}

// 排序
type docCountSort []*queryTokenHash

func (q docCountSort) Less(i, j int) bool {
	return q[i].invertedIndex.DocCount < q[i].invertedIndex.DocCount
}

func (q docCountSort) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q docCountSort) Len() int {
	return len(q)
}
