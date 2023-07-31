package segment

import (
	"bytes"
	"fmt"

	"github.com/CocaineCong/tangseng/app/search-engine/internal/query"
	"github.com/CocaineCong/tangseng/app/search-engine/internal/storage"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

type Segment struct {
	*storage.ForwardDB  // 正排索引库
	*storage.InvertedDB // 倒排索引库
}

// Token2PostingsLists 词条 转化成 倒排索引表
func Token2PostingsLists(bufInvertHash InvertedIndexHash, token query.Tokenization, docId int64) (err error) {
	bufInvert := new(InvertedIndexValue)
	if len(bufInvertHash) > 0 {
		if item, ok := bufInvertHash[token.Token]; ok {
			bufInvert = item
		}
	}

	pl := new(PostingsList)
	if bufInvert != nil && bufInvert.PostingsList != nil {
		pl = bufInvert.PostingsList
		pl.PositionCount++
	} else {
		var docCount int64 = 1
		bufInvert = CreateNewInvertedIndex(token.Token, docCount)
		bufInvertHash[token.Token] = bufInvert
		pl = CreateNewPostingsList(docId)
		bufInvert.PostingsList = pl
	}

	pl.Positions = append(pl.Positions, token.Position) // 存储位置信息
	bufInvert.PositionCount++                           // 统计该token关联的所有的doc的position的个数

	return
}

// getTokenCount 通过token获取doc数量 insert 标识是写入还是查询 写入时不为空
func (e *Segment) getTokenCount(token string) (termInfo *storage.TermValue, err error) {
	termInfo, err = e.InvertedDB.GetTermInfo(token)
	if err != nil || termInfo == nil {
		log.LogrusObj.Errorf("getTokenCount GetTermInfo err:%v", err)
		return
	}

	return
}

// FetchPostings 通过 token 读取倒排表数据，返回倒排表，长度 和 err
func (e *Segment) FetchPostings(token string) (*PostingsList, int64, error) {
	term, err := e.InvertedDB.GetTermInfo(token)
	if err != nil {
		return nil, 0, fmt.Errorf("FetchPostings getForwardAddr err: %v", err)
	}

	c, err := e.InvertedDB.GetInvertedDoc(term.Offset, term.Size)
	if err != nil {
		return nil, 0, fmt.Errorf("FetchPostings getForwardAddr err: %v", err)
	}

	return decodePostings(bytes.NewBuffer(c))
}

// Flush 落盘操作
func (e *Segment) Flush(PostingsHashBuf InvertedIndexHash) error {
	if len(PostingsHashBuf) == 0 {
		log.LogrusObj.Infof("Flush err: %v", "in.PostingsHashBuf is empty")
		return nil
	}
	for token, invertedIndex := range PostingsHashBuf {
		log.LogrusObj.Infof("token:%s,invertedIndex:%v\n", token, invertedIndex)
		err := e.storagePostings(invertedIndex)
		if err != nil {
			log.LogrusObj.Infof("updatePostings err: %v", err)
			return fmt.Errorf("updatePostings err: %v", err)
		}
	}
	return nil
}

// storagePostings 落盘
func (e *Segment) storagePostings(p *InvertedIndexValue) error {
	if p == nil {
		return fmt.Errorf("updatePostings p is nil")
	}

	// 编码
	buf, err := EncodePostings(p.PostingsList, p.DocCount)
	if err != nil {
		return fmt.Errorf("updatePostings encodePostings err: %v", err)
	}

	// 开始写入数据库
	return e.InvertedDB.StoragePostings(p.Token, buf.Bytes(), p.DocCount)
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
	inDb, forDb := InitSegmentDb(segId)
	return &Segment{
		InvertedDB: inDb,
		ForwardDB:  forDb,
	}
}
