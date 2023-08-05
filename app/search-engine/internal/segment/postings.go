package segment

import (
	"bytes"
	"encoding/gob"

	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/util/codec"
)

// PostingsList 倒排列表
type PostingsList struct {
	DocId         int64
	Positions     []int64
	PositionCount int64
	Next          *PostingsList
}

// MergePostings 合并两个posting
func MergePostings(pa, pb *PostingsList) *PostingsList {
	ret := new(PostingsList)
	p := new(PostingsList)
	p = nil
	for pa != nil || pb != nil {
		tmp := new(PostingsList)
		if pb == nil || (pa != nil && pa.DocId <= pb.DocId) {
			tmp, pa = pa, pa.Next
		} else if pa == nil || (pa != nil && pa.DocId > pb.DocId) {
			tmp, pb = pb, pb.Next
		} else {
			log.LogrusObj.Infoln("break")
			break
		}
		tmp.Next = nil

		if p == nil {
			ret.Next = tmp
		} else {
			p.Next = tmp
		}
		p = tmp
	}

	return ret.Next
}

// MergeInvertedIndex 合并两个倒排索引
func MergeInvertedIndex(base, toBeAdd InvertedIndexHash) {
	for token, index := range base {
		if toBeAddedIndex, ok := (toBeAdd)[token]; ok {
			index.PostingsList = MergePostings(index.PostingsList, toBeAddedIndex.PostingsList)
			index.DocCount += toBeAddedIndex.DocCount
			delete(toBeAdd, token)
		}
	}

	for tokenId, index := range toBeAdd {
		(base)[tokenId] = index
	}

}

// DecodePostings 解码 return *PostingsList postingslen err
func DecodePostings(buf *bytes.Buffer) (p *PostingsList, postingsLen int64, err error) {
	if buf == nil || buf.Len() == 0 {
		log.LogrusObj.Infoln("DecodePostings-buf 为空")
		return
	}

	dec := gob.NewDecoder(buf)
	err = dec.Decode(&postingsLen)
	if err != nil {
		log.LogrusObj.Errorln("binary.Read", err)
		return
	}

	cp := new(PostingsList)
	p = cp
	for buf.Len() > 0 {
		tmp := new(PostingsList)
		err = dec.Decode(&tmp.DocId)
		if err != nil {
			log.LogrusObj.Errorln("binary.Read", err)
			return
		}

		err = dec.Decode(&tmp.PositionCount)
		if err != nil {
			log.LogrusObj.Errorln("binary.Read", err)
			return
		}

		tmp.Positions = make([]int64, tmp.PositionCount)
		err = dec.Decode(&tmp.Positions)
		if err != nil {
			log.LogrusObj.Errorln("binary.Read", err)
			return
		}
		log.LogrusObj.Infoln("postings", tmp)
		cp.Next = tmp
		cp = tmp

	}

	return p.Next, postingsLen, nil
}

// EncodePostings 编码
func EncodePostings(postings *PostingsList, postingsLen int64) (buf *bytes.Buffer, err error) {
	buf, err = codec.GobWrite(postingsLen)
	if err != nil {
		return
	}

	for postings != nil {
		log.LogrusObj.Infof("docid:%d,count:%d,positions:%v \n", postings.DocId, postings.PositionCount, postings.Positions)
		buf, err = codec.GobWrite(postings.DocId)
		if err != nil {
			return
		}
		buf, err = codec.GobWrite(postings.PositionCount)
		if err != nil {
			return
		}
		buf, err = codec.GobWrite(postings.Positions)
		if err != nil {
			return
		}
		postings = postings.Next
	}

	return
}

// CreateNewPostingsList 创建倒排索引
func CreateNewPostingsList(docId int64) *PostingsList {
	p := new(PostingsList)
	p.DocId = docId
	p.PositionCount = 1
	p.Positions = make([]int64, 0)
	return p
}
