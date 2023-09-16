package segment

import (
	"github.com/CocaineCong/tangseng/app/search_engine/analyzer"
	"github.com/CocaineCong/tangseng/app/search_engine/repository/storage"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/trie"
	"github.com/CocaineCong/tangseng/pkg/util/codec"
	"github.com/CocaineCong/tangseng/types"
)

type Segment struct {
	*storage.ForwardDB  // 正排索引库
	*storage.InvertedDB // 倒排索引库
	*storage.DictDB     // 存储trie树
}

// Token2PostingsLists 词条 转化成 倒排索引表
func Token2PostingsLists(bufInvertHash InvertedIndexHash, token analyzer.Tokenization, docId int64) (err error) {
	bufInvert := new(types.InvertedIndexValue)
	if len(bufInvertHash) > 0 {
		if item, ok := bufInvertHash[token.Token]; ok {
			bufInvert = item
		}
	}

	pl := new(types.PostingsList)
	if bufInvert != nil && bufInvert.PostingsList != nil {
		pl = bufInvert.PostingsList
		pl.PositionCount++
	} else {
		var docCount int64 = 1
		bufInvert = CreateNewInvertedIndex(token, docCount)
		bufInvertHash[token.Token] = bufInvert
		pl = CreateNewPostingsList(docId)
		bufInvert.PostingsList = pl
	}

	pl.Positions = append(pl.Positions, token.Position) // 存储位置信息
	bufInvert.PositionCount++                           // 统计该token关联的所有的doc的position的个数

	return
}

// getTokenCount 通过token获取doc数量 insert 标识是写入还是查询 写入时不为空
func (e *Segment) getTokenCount(token string) (termInfo *types.TermValue, err error) {
	termInfo, err = e.InvertedDB.GetTermInfo(token)
	if err != nil || termInfo == nil {
		log.LogrusObj.Errorf("getTokenCount GetTermInfo err:%v", err)
		return
	}

	return
}

// FetchPostings 通过 token 读取倒排表数据，返回倒排索引
func (e *Segment) FetchPostings(token string) (p *types.InvertedIndexValue, err error) {
	p, err = e.InvertedDB.GetInvertedInfo(token)
	if err != nil {
		log.LogrusObj.Errorf("FetchPostings GetInvertedDoc err: %v", err)
		return
	}
	return
}

// FlushInvertedIndex 落盘操作
func (e *Segment) FlushInvertedIndex(PostingsHashBuf InvertedIndexHash) (err error) {
	if len(PostingsHashBuf) == 0 {
		log.LogrusObj.Infof("Flush err: %v", "in.PostingsHashBuf is empty")
		return
	}
	for token, invertedIndex := range PostingsHashBuf {
		log.LogrusObj.Infof("token:%s,invertedIndex:%v \n", token, invertedIndex)
		err = e.storagePostings(invertedIndex)
		if err != nil {
			log.LogrusObj.Errorf("Flush-storagePostings err: %v", err)
			return
		}
	}
	return
}

// FlushTokenDict 刷新写入 token dict
func (e *Segment) FlushTokenDict(trieTree *trie.Trie) (err error) {
	err = e.StorageDict(trieTree)

	return
}

// storagePostings 落盘
func (e *Segment) storagePostings(p *types.InvertedIndexValue) (err error) {
	if p == nil {
		log.LogrusObj.Errorf("updatePostings p is nil")
		return
	}

	// 编码
	buf, err := codec.EncodePostings(p)
	if err != nil {
		log.LogrusObj.Errorf("updatePostings encodePostings err: %v", err)
		return
	}

	// 开始写入数据库
	return e.InvertedDB.StoragePostings(p.Token, buf)
}

// Close --
func (e *Segment) Close() {
	e.InvertedDB.Close()
	e.ForwardDB.Close()
}

// NewSegments 创建新的segments 更新next seg
func NewSegments(meta *SegMeta, mode Mode) (SegId, map[SegId]*Segment) {
	segs := make(map[SegId]*Segment, 0)
	if mode == MergeMode || mode == IndexMode {
		segId := meta.NextSeg
		err := meta.NewSegmentItem()
		if err != nil {
			return 0, nil
		}
		seg := NewSegment(segId)
		segs[segId] = seg
		return segId, segs
	}
	log.LogrusObj.Infof("meta:%v", meta)
	for segId := range meta.SegInfo {
		seg := NewSegment(segId)
		log.LogrusObj.Infof("db init segId:%v,next:%v", segId, meta.NextSeg)
		segs[segId] = seg
	}

	return -1, segs
}

func NewSegment(segId SegId) *Segment {
	inDb, forDb, dictDb, _ := InitSegmentDb(segId)
	return &Segment{
		InvertedDB: inDb,
		ForwardDB:  forDb,
		DictDB:     dictDb,
	}
}
