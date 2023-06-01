package engine

import (
	"fmt"

	"github.com/CocaineCong/Go-SearchEngine/app/search-engine/internal/query"
	"github.com/CocaineCong/Go-SearchEngine/app/search-engine/internal/segment"
	"github.com/CocaineCong/Go-SearchEngine/app/search-engine/internal/storage"
	log "github.com/CocaineCong/Go-SearchEngine/pkg/logger"
)

// ErrCountKeyNotFound 计数key不存在
var ErrCountKeyNotFound = "get token:forwardCount err:key not found"

// Engine 写入引擎
type Engine struct {
	meta      *Meta // 元数据
	Scheduler *MergeScheduler

	BufCount        int64                              // 倒排索引 缓冲区的文档数
	BufSize         int64                              // 设定的缓冲区大小
	PostingsHashBuf segment.InvertedIndexHash          // 倒排索引缓冲区
	CurrSegId       segment.SegId                      // 当前engine关联的segId查询
	Seg             map[segment.SegId]*segment.Segment // 当前engine关联的segment

	N int64 // ngram
}

// NewEngine 每次初始化的时候调整meta数据
func NewEngine(meta *Meta, engineMode segment.Mode) *Engine {
	sche := NewScheduler(meta)
	segId, seg := segment.NewSegments(meta.SegMeta, engineMode)
	return &Engine{
		meta:            meta,
		Scheduler:       sche,
		BufSize:         5,
		PostingsHashBuf: make(segment.InvertedIndexHash),
		CurrSegId:       segId,
		Seg:             seg,
		N:               2,
	}
}

// Close --
func (e *Engine) Close() {
	for _, seg := range e.Seg {
		seg.Close()
	}

	e.Scheduler.Close()
}

// indexCount index 计数
func (e *Engine) indexCount() {
	e.meta.Lock()
	e.meta.IndexCount++
	e.meta.Unlock()
}

// AddDoc 添加正排
func (e *Engine) AddDoc(doc *storage.Document) error {
	return e.Seg[e.CurrSegId].AddForwardByDoc(doc)
}

// Text2PostingsLists --
func (e *Engine) Text2PostingsLists(text string, docId int64) error {
	tokens, err := query.Ngram(text, docId)
	if err != nil {
		return fmt.Errorf("text2PostingsLists Ngram err:%v", err)
	}

	bufInvertedHash := make(segment.InvertedIndexHash)
	for _, token := range tokens {
		err := segment.Token2PostingsLists(bufInvertedHash, token.Token, token.Position, docId)
		if err != nil {
			return err
		}
	}

	log.LogrusObj.Infof("buf InvertedHash :%v", bufInvertedHash)

	if e.PostingsHashBuf != nil && len(e.PostingsHashBuf) > 0 {
		// 合并命中相同的token的不同doc
		segment.MergeInvertedIndex(e.PostingsHashBuf, bufInvertedHash)
	} else {
		// 已经初始化过了
		e.PostingsHashBuf = bufInvertedHash
	}

	e.BufCount++

	// 达到阈值
	if len(e.PostingsHashBuf) > 0 && (e.BufCount >= e.BufSize) {
		log.LogrusObj.Infof("text2PostingsLists need flush")
		e.Flush()
	}

	e.indexCount()
	return nil
}

func (e *Engine) UpdateCount(num int64) error {
	seg := e.Seg[e.CurrSegId]
	count, err := seg.ForwardCount()
	if err != nil {
		if err.Error() == ErrCountKeyNotFound {
			count = 0
		} else {
			return fmt.Errorf("updateCount err:%v", err)
		}
	}
	count += num
	return seg.UpdateForwardCount(count)
}

func (e *Engine) Flush(isEnd ...bool) error {
	e.Seg[e.CurrSegId].Flush(e.PostingsHashBuf)

	// 更新 meta info
	err := e.meta.UpdateSegMeta(e.CurrSegId, e.BufCount)
	if err != nil {
		return err
	}

	e.UpdateCount(e.meta.IndexCount)
	e.Seg[e.CurrSegId].Close()
	delete(e.Seg, e.CurrSegId)

	if len(e.meta.SegMeta.SegInfo) > 1 {
		e.Scheduler.MayMerge()
	}

	// new
	if len(isEnd) > 0 && isEnd[0] {
		return nil
	}

	segId, seg := segment.NewSegments(e.meta.SegMeta, segment.IndexMode)

	e.BufCount = 0
	e.PostingsHashBuf = make(segment.InvertedIndexHash)
	e.CurrSegId = segId
	e.Seg = seg

	return nil
}
