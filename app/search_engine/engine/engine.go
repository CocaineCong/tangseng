package engine

import (
	"sync"
	"sync/atomic"

	"github.com/CocaineCong/tangseng/app/search_engine/segment"
	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/pkg/trie"
)

// Engine 写入引擎
type Engine struct {
	meta            *Meta                              // 元数据
	Scheduler       *MergeScheduler                    // 合并调度器
	BufCount        int64                              // 倒排索引 缓冲区的文档数
	BufSize         int64                              // 设定的缓冲区大小
	PostingsHashBuf segment.InvertedIndexHash          // 倒排索引缓冲区
	TrieTree        *trie.Trie                         // 词典前缀树
	CurrSegId       segment.SegId                      // 当前engine关联的segId查询
	Seg             map[segment.SegId]*segment.Segment // 当前engine关联的segment
	// TODO 更换并发安全的map，需要写入性能好的
}

var TangSengEngineIns *Engine
var TangSengEngineOnce sync.Once

// NewTangSengEngine 每次初始化的时候调整meta数据
func NewTangSengEngine(meta *Meta, engineMode segment.Mode) *Engine {
	TangSengEngineOnce.Do(func() {
		segId, seg := segment.NewSegments(meta.SegMeta, engineMode)
		TangSengEngineIns = &Engine{
			meta:            meta,
			Scheduler:       NewScheduler(meta),
			BufSize:         consts.EngineBufSize,
			PostingsHashBuf: make(segment.InvertedIndexHash),
			TrieTree:        trie.NewTrie(),
			CurrSegId:       segId,
			Seg:             seg,
		}
	})

	return TangSengEngineIns
}

// indexCount index 计数
func (e *Engine) indexCount() {
	atomic.AddInt64(&e.meta.IndexCount, 1)
}

// Close --
func (e *Engine) Close() {
	for _, seg := range e.Seg {
		seg.Close()
	}

	e.Scheduler.Close()
}
