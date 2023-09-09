package engine

import (
	"github.com/CocaineCong/tangseng/app/search_engine/analyzer"
	"github.com/CocaineCong/tangseng/app/search_engine/segment"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/trie"
)

// FlushInvertedIndex 倒排索引落盘操作
func (e *Engine) FlushInvertedIndex(isEnd ...bool) (err error) {
	err = e.Seg[e.CurrSegId].FlushInvertedIndex(e.PostingsHashBuf)
	if err != nil {
		log.LogrusObj.Errorln("FlushInvertedIndex-FlushInvertedIndex", err)
		return
	}

	// 更新 meta info
	err = e.meta.UpdateSegMeta(e.CurrSegId, e.BufCount)
	if err != nil {
		log.LogrusObj.Errorln("FlushInvertedIndex-UpdateSegMeta", err)
		return
	}

	err = e.UpdateCount(e.meta.IndexCount)
	if err != nil {
		log.LogrusObj.Errorln("FlushInvertedIndex-UpdateCount", err)
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
	e.TrieTree = trie.NewTrie()
	e.CurrSegId = segId
	e.Seg = seg

	return
}

// Text2PostingsLists 建立索引专用 文本 转成 倒排索引记录表
func (e *Engine) Text2PostingsLists(text string, docId int64) (err error) {
	tokens, err := analyzer.GseCutForBuildIndex(text)
	if err != nil {
		log.LogrusObj.Errorf("Text2PostingsLists-GseCut err:%v", err)
		return
	}

	bufInvertedHash := make(segment.InvertedIndexHash)
	trieTree := trie.NewTrie()
	for _, token := range tokens {
		err = segment.Token2PostingsLists(bufInvertedHash, token, docId)
		if err != nil {
			log.LogrusObj.Errorf("Text2PostingsLists-Token2PostingsLists err:%v", err)
			return
		}
		trieTree.Insert(token.Token)
	}

	e.TrieTree.Merge(trieTree)
	log.LogrusObj.Infof("InvertedHash: %v", bufInvertedHash)

	if e.PostingsHashBuf != nil && len(e.PostingsHashBuf) > 0 {
		// 合并命中相同的token的不同doc
		segment.MergeInvertedIndex(e.PostingsHashBuf, bufInvertedHash)
	} else {
		e.PostingsHashBuf = bufInvertedHash // 已经初始化过了
	}
	e.BufCount++
	// 达到阈值，刷新存储
	if len(e.PostingsHashBuf) > 0 && (e.BufCount >= e.BufSize) {
		log.LogrusObj.Infof("text2PostingsLists need flush")
		err = e.FlushDict()
		if err != nil {
			log.LogrusObj.Errorf("FlushDict err:%v", err)
			return
		}

		err = e.FlushInvertedIndex()
		if err != nil {
			log.LogrusObj.Errorf("FlushInvertedIndex err:%v", err)
			return
		}
	}

	e.indexCount()
	return
}

// Text2PostingsListsForRecall 召回专用 文本 转成 倒排索引记录表
func (e *Engine) Text2PostingsListsForRecall(text string, docId int64) (err error) {
	// TODO 后面加上推荐词的发现 for example: query:陆家嘴 推荐词:东方明珠, 上海迪士尼 or 南京东路
	tokens, err := analyzer.GseCutForRecall(text)
	if err != nil {
		log.LogrusObj.Errorf("Text2PostingsLists-GseCut err:%v", err)
		return
	}

	bufInvertedHash := make(segment.InvertedIndexHash)
	trieTree := trie.NewTrie()
	for _, token := range tokens {
		err = segment.Token2PostingsLists(bufInvertedHash, token, docId)
		if err != nil {
			log.LogrusObj.Errorf("Text2PostingsLists-Token2PostingsLists err:%v", err)
			return
		}
		trieTree.Insert(token.Token)
	}

	e.TrieTree.Merge(trieTree)
	log.LogrusObj.Infof("InvertedHash: %v", bufInvertedHash)

	if e.PostingsHashBuf != nil && len(e.PostingsHashBuf) > 0 {
		// 合并命中相同的token的不同doc
		segment.MergeInvertedIndex(e.PostingsHashBuf, bufInvertedHash)
	} else {
		e.PostingsHashBuf = bufInvertedHash // 已经初始化过了
	}
	e.BufCount++
	// 达到阈值，刷新存储
	if len(e.PostingsHashBuf) > 0 && (e.BufCount >= e.BufSize) {
		log.LogrusObj.Infof("text2PostingsLists need flush")
		err = e.FlushDict()
		if err != nil {
			log.LogrusObj.Errorf("FlushDict err:%v", err)
			return
		}

		err = e.FlushInvertedIndex()
		if err != nil {
			log.LogrusObj.Errorf("FlushInvertedIndex err:%v", err)
			return
		}
	}

	e.indexCount()
	return
}
