package engine

import (
	"sync"
	"sync/atomic"

	"github.com/CocaineCong/tangseng/app/search-engine/internal/query"
	"github.com/CocaineCong/tangseng/app/search-engine/internal/segment"
	"github.com/CocaineCong/tangseng/app/search-engine/internal/storage"
	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
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

// Close --
func (e *Engine) Close() {
	for _, seg := range e.Seg {
		seg.Close()
	}

	e.Scheduler.Close()
}

// indexCount index 计数
func (e *Engine) indexCount() {
	atomic.AddInt64(&e.meta.IndexCount, 1)
}

// AddDoc 添加正排
func (e *Engine) AddDoc(doc *storage.Document) error {
	return e.Seg[e.CurrSegId].AddForwardByDoc(doc)
}

// Text2PostingsLists --
func (e *Engine) Text2PostingsLists(text string, docId int64) (err error) {
	tokens, err := query.GseCut(text)
	if err != nil {
		log.LogrusObj.Errorf("text2PostingsLists err:%v", err)
		return
	}

	bufInvertedHash := make(segment.InvertedIndexHash)
	for _, token := range tokens {
		err = segment.Token2PostingsLists(bufInvertedHash, token, docId)
		if err != nil {
			log.LogrusObj.Errorf("Token2PostingsLists err:%v", err)
			return
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

	// 达到阈值，刷新存储
	if len(e.PostingsHashBuf) > 0 && (e.BufCount >= e.BufSize) {
		log.LogrusObj.Infof("text2PostingsLists need flush")
		err = e.Flush()
		if err != nil {
			log.LogrusObj.Errorf("Flush err:%v", err)
			return
		}
	}

	e.indexCount()
	return nil
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

// Flush 落盘操作
func (e *Engine) Flush(isEnd ...bool) (err error) {
	err = e.Seg[e.CurrSegId].Flush(e.PostingsHashBuf)
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
		return nil
	}

	segId, seg := segment.NewSegments(e.meta.SegMeta, segment.IndexMode)

	e.BufCount = 0
	e.PostingsHashBuf = make(segment.InvertedIndexHash)
	e.CurrSegId = segId
	e.Seg = seg

	return nil
}
