package engine

import (
	"sync"
	"sync/atomic"

	"github.com/spf13/cast"

	"github.com/CocaineCong/tangseng/app/search_engine/query"
	"github.com/CocaineCong/tangseng/app/search_engine/segment"
	"github.com/CocaineCong/tangseng/app/search_engine/types"
	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/trie"
)

// Engine 写入引擎
type Engine struct {
	meta            *Meta                              // 元数据
	Scheduler       *MergeScheduler                    // 合并调度器
	BufCount        int64                              // 倒排索引 缓冲区的文档数
	BufSize         int64                              // 设定的缓冲区大小
	PostingsHashBuf segment.InvertedIndexHash          // 倒排索引缓冲区
	CurrSegId       segment.SegId                      // 当前engine关联的segId查询
	Seg             map[segment.SegId]*segment.Segment // 当前engine关联的segment
	// TODO 更换并发安全的map，需要写入性能好的
}

var EngineIns *Engine
var EngineOnce sync.Once

// NewEngine 每次初始化的时候调整meta数据
func NewEngine(meta *Meta, engineMode segment.Mode) *Engine {
	EngineOnce.Do(func() {
		segId, seg := segment.NewSegments(meta.SegMeta, engineMode)
		EngineIns = &Engine{
			meta:            meta,
			Scheduler:       NewScheduler(meta),
			BufSize:         consts.EngineBufSize,
			PostingsHashBuf: make(segment.InvertedIndexHash),
			CurrSegId:       segId,
			Seg:             seg,
		}
	})

	return EngineIns
}

// indexCount index 计数
func (e *Engine) indexCount() {
	atomic.AddInt64(&e.meta.IndexCount, 1)
}

// AddForwardIndex 落库正排索引
func (e *Engine) AddForwardIndex(doc *types.Document) error {
	return e.Seg[e.CurrSegId].AddForwardByDoc(doc)
}

// Text2PostingsLists 文本 转成 倒排索引记录表
func (e *Engine) Text2PostingsLists(text string, docId int64) (err error) {
	tokens, err := query.GseCut(text)
	if err != nil {
		log.LogrusObj.Errorf("text2PostingsLists err:%v", err)
		return
	}

	bufInvertedHash := make(segment.InvertedIndexHash)
	trieTree := new(trie.Trie)
	for _, token := range tokens {
		err = segment.Token2PostingsLists(bufInvertedHash, token, docId)
		if err != nil {
			log.LogrusObj.Errorf("Token2PostingsLists err:%v", err)
			return
		}
		trieTree.Insert(token.Token)
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

	// 达到阈值，刷新存储
	if len(e.PostingsHashBuf) > 0 && (e.BufCount >= e.BufSize) {
		log.LogrusObj.Infof("text2PostingsLists need flush")

		err = e.FlushDict(trieTree)
		if err != nil {
			log.LogrusObj.Errorf("Flush err:%v", err)
			return
		}

		err = e.FlushInvertedIndex()
		if err != nil {
			log.LogrusObj.Errorf("Flush err:%v", err)
			return
		}
	}

	e.indexCount()
	return
}

func (e *Engine) UpdateCount(num int64) (err error) {
	seg := e.Seg[e.CurrSegId]
	count, err := seg.ForwardCount()
	if err != nil {
		log.LogrusObj.Errorf("updateCount err:%v", err)
		return
	}
	count += num
	return seg.UpdateForwardCount(count)
}

// FlushInvertedIndex 倒排索引落盘操作
func (e *Engine) FlushInvertedIndex(isEnd ...bool) (err error) {
	err = e.Seg[e.CurrSegId].FlushInvertedIndex(e.PostingsHashBuf)
	if err != nil {
		log.LogrusObj.Errorln("Flush", err)
		return
	}

	// 更新 meta info
	err = e.meta.UpdateSegMeta(e.CurrSegId, e.BufCount)
	if err != nil {
		log.LogrusObj.Errorln("UpdateSegMeta", err)
		return
	}

	err = e.UpdateCount(e.meta.IndexCount)
	if err != nil {
		log.LogrusObj.Errorln("UpdateCount", err)
		return
	}
	e.Seg[e.CurrSegId].Close()
	delete(e.Seg, e.CurrSegId)

	if len(e.meta.SegMeta.SegInfo) > 1 {
		e.Scheduler.MayMerge()
	}

	if len(isEnd) > 0 && isEnd[0] {
		return
	}

	segId, seg := segment.NewSegments(e.meta.SegMeta, segment.IndexMode)

	e.BufCount = 0
	e.PostingsHashBuf = make(segment.InvertedIndexHash)
	e.CurrSegId = segId
	e.Seg = seg

	return
}

// FlushDict 刷新dict
func (e *Engine) FlushDict(trieTree *trie.Trie, isEnd ...bool) (err error) {
	currSegId := cast.ToInt64(e.CurrSegId)
	err = e.Seg[e.CurrSegId].FlushTokenDict(currSegId, trieTree)
	if err != nil {
		log.LogrusObj.Errorln("Flush", err)
		return
	}

	if len(isEnd) > 0 && isEnd[0] {
		return
	}

	return
}

// Close --
func (e *Engine) Close() {
	for _, seg := range e.Seg {
		seg.Close()
	}

	e.Scheduler.Close()
}
