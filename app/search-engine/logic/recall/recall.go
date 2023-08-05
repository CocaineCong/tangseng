package recall

import (
	"errors"
	"sort"

	"github.com/CocaineCong/tangseng/app/search-engine/logic/engine"
	"github.com/CocaineCong/tangseng/app/search-engine/logic/segment"
	"github.com/CocaineCong/tangseng/app/search-engine/logic/types"
	log "github.com/CocaineCong/tangseng/pkg/logger"
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
	invertedIndex *types.InvertedIndexValue
	fetchPostings *types.PostingsList
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
	doc     *types.PostingsList // 文档编号的序列
	current *types.PostingsList // 当前文档编号
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

	return r.searchDoc()
}

func (r *Recall) splitQuery2Tokens(query string) (err error) {
	err = r.Text2PostingsLists(query, 0)
	if err != nil {
		log.LogrusObj.Errorf("text2postingslists err: %v", err)
		return
	}
	return
}

func (r *Recall) searchDoc() (recall Recalls, err error) {
	recalls := make(Recalls, 0)
	tokens := make([]*queryTokenHash, 0)

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
		t := &queryTokenHash{
			token:         token,
			invertedIndex: post,
			fetchPostings: postings,
		}
		tokens = append(tokens, t)
	}

	tokens = r.sortToken(tokens)

	tokenCount := len(tokens)
	if tokenCount == 0 {
		log.LogrusObj.Infof("searchDoc-tokenCount is 0")
		return
	}

	cursors := make([]searchCursor, tokenCount)
	for i, t := range tokens {
		cursors[i].doc = t.fetchPostings
		cursors[i].current = t.fetchPostings
	}

	// 整个遍历token来匹配doc
	for cursors[0].current != nil {
		var docId, nextDocId int64
		// 拥有文档最少的token作为标尺
		docId = cursors[0].current.DocId

		// 匹配其他token的doc
		for i := 1; i < tokenCount; i++ {
			cur := &cursors[i]
			for cur.current != nil && cur.current.DocId < docId {
				cur.current = cur.current.Next
			}

			// 存在token关联的docid都小雨cursors[0]的docid,则跳出
			if cur.current == nil {
				log.LogrusObj.Infof("cur.current is nil\n")
				break
			}

			// 对于除词元A以外的词元，如果其他document_id不等于词元A的document_id,那么就将这个document_id设定为next_doc_id
			if cur.current.DocId != docId {
				nextDocId = cur.current.DocId
				break
			}
		}

		log.LogrusObj.Infof("当前doc id：%v，next doc id:%v", docId, nextDocId)
		if nextDocId > 0 {
			// 不断获取A的下一个document_id，直到其当前的document_id不小于next_doc_id为止
			for cursors[0].current != nil && cursors[0].current.DocId < nextDocId {
				cursors[0].current = cursors[0].current.Next
			}
		} else {
			// 有匹配的docid
			phraseCount := int64(-1)
			if r.enablePhrase {
				phraseCount = r.searchPhrase(tokens, cursors)
			}
			score := 0.0
			if phraseCount > 0 {
				r.calculateScore(cursors, int64(tokenCount)) // TODO:计算相关性
			}
			cursors[0].current = cursors[0].current.Next
			log.LogrusObj.Infof("匹配召回docID:%v,nextDocID:%v,phrase:%d", docId, nextDocId, phraseCount)
			recalls = append(recalls, &SearchItem{DocId: docId, Score: score})
		}
	}
	log.LogrusObj.Infof("recalls size:%v", len(recalls))

	return recalls, nil
}

// calculateScore 计算相关性
func (r *Recall) calculateScore(cursor []searchCursor, tokenCount int64) float64 {
	return 0.0
}

// searchPhrase 返回检索出的短语数 查询query的倒排索引 tokenCursors是fetched文档的倒排索引
func (r *Recall) searchPhrase(queryToken []*queryTokenHash, tokenCursors []searchCursor) int64 {
	// 获取遍历查询query分词之后的词元总数
	positionsSum := int64(0)
	for _, t := range queryToken {
		positionsSum += t.invertedIndex.PositionCount
	}
	cursors := make([]phraseCursor, positionsSum)
	phraseCount := int64(0)
	// 初始化游标 获取token关联的第一篇doc的pos相关数据
	n := 0
	for i, t := range queryToken {
		for _, pos := range t.invertedIndex.PostingsList.Positions {
			cursors[n].base = pos                                    // 记录查询中出现的位置
			cursors[n].positions = tokenCursors[i].current.Positions // 获取token关联的文件中token
			cursors[n].current = &cursors[i].positions[0]            // 获取文档中出现的位置
			cursors[n].index = 0                                     // 获取文档中出现的索引位置
			log.LogrusObj.Infof("token:%s,pos:%v cur:%v,positions:%v",
				t.token, pos, *cursors[n].current, cursors[n].positions)
			n++
		}
	}

	for cursors[0].current != nil {
		var relPos, nextRelPos int64
		relPos = *cursors[0].current - cursors[0].base
		nextRelPos = relPos
		// 对于除词元A以外的词元，不断地向后读取其出现位置，直到其偏移量不小于词元A的偏移量为止
		for i := 1; i < len(cursors); i++ {
			cur := &cursors[i]
			for cur.current != nil && *cur.current-cur.base < relPos {
				cur.index++
				if cur.index >= len(cur.positions) {
					log.LogrusObj.Infof("cur.index >= len(cur.positions)\n")
					cur.current = nil
					break
				}
				cur.current = &cur.positions[cur.index]
			}
			if cur.current == nil {
				break
			}
			if *cur.current-cur.base != relPos {
				nextRelPos = *cur.current - cur.base
				break
			}
		}

		if nextRelPos > relPos {
			// 不断向后读取，直到词元A的偏移量不小于next rel position为止
			for cursors[0].current != nil && *cursors[0].current-cursors[0].base < nextRelPos {
				cursors[0].index++
				if cursors[0].index >= len(cursors[0].positions) {
					log.LogrusObj.Infof("cursors[0].index >= len(cursors[0].positions)\n")
					cursors[0].current = nil
					break
				}
				cursors[0].current = &cursors[0].positions[cursors[0].index]
			}
		} else {
			// 找到短语
			phraseCount++
			cursors[0].index++
			// 判断是否有下一个命中的短语
			if cursors[0].index >= len(cursors[0].positions) {
				log.LogrusObj.Infof("cursors[0].index:%d>= len(cursors[0].positions):%d",
					cursors[0].index, len(cursors[0].positions))
				cursors[0].current = nil
			} else {
				cursors[0].current = &cursors[0].positions[cursors[0].index]
			}
		}
	}

	return phraseCount
}

// token 根据 doc count 升序排序，回去之后还要再进行一次按照score的排序
func (r *Recall) sortToken(tokens []*queryTokenHash) []*queryTokenHash {
	// 检验是否排序成功
	for _, t := range tokens {
		log.LogrusObj.Infof("token:%v,docCount:%v", t.token, t.invertedIndex.DocCount)
	}
	sort.Sort(docCountSort(tokens))
	for _, t := range tokens {
		log.LogrusObj.Infof("token:%v,docCount:%v", t.token, t.invertedIndex.DocCount)
	}

	return tokens
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

// 排序
type docCountSort []*queryTokenHash

func (q docCountSort) Less(i, j int) bool {
	return q[i].invertedIndex.DocCount < q[j].invertedIndex.DocCount
}

func (q docCountSort) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q docCountSort) Len() int {
	return len(q)
}
