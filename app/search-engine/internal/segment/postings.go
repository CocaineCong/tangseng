package segment

import (
	"bytes"
	"encoding/binary"

	log "github.com/CocaineCong/Go-SearchEngine/pkg/logger"
	"github.com/CocaineCong/Go-SearchEngine/pkg/util/se"
)

// PostingsList 倒排列表
type PostingsList struct {
	DocId         int64
	Positions     []int64
	PositionCount int64
	Next          *PostingsList
}

// MergePostings 合并两个东西
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

// decodePostings 解码 return *PostingsList postingslen err
func decodePostings(buf *bytes.Buffer) (*PostingsList, int64, error) {
	if buf == nil || buf.Len() == 0 {
		return nil, 0, nil
	}
	var postingsLen int64
	err := binary.Read(buf, binary.LittleEndian, &postingsLen)
	if err != nil {
		log.LogrusObj.Infoln("binary.Read", err)
		return nil, 0, err
	}
	cp := new(PostingsList)
	p := cp
	for buf.Len() > 0 {
		tmp := new(PostingsList)
		err = binary.Read(buf, binary.LittleEndian, &tmp.DocId)
		if err != nil {
			log.LogrusObj.Infoln("binary.Read", err)
			return nil, 0, err
		}

		err = binary.Read(buf, binary.LittleEndian, &tmp.PositionCount)
		if err != nil {
			return nil, 0, err
		}

		tmp.Positions = make([]int64, tmp.PositionCount)
		err = binary.Read(buf, binary.LittleEndian, &tmp.Positions)
		if err != nil {
			return nil, 0, err
		}
		log.LogrusObj.Infoln("postings", tmp)
		cp.Next = tmp
		cp = tmp

	}
	return p.Next, postingsLen, nil
}

// EncodePostings 编码
func EncodePostings(postings *PostingsList, postingsLen int64) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})
	err := se.BinaryWrite(buf, postingsLen)
	if err != nil {
		return nil, err
	}

	for postings != nil {
		log.LogrusObj.Infof("docid:%d,count:%d,positions:%v \n", postings.DocId, postings.PositionCount, postings.Positions)
		err := se.BinaryWrite(buf, postings.DocId)
		if err != nil {
			return nil, err
		}
		err = se.BinaryWrite(buf, postings.PositionCount)
		if err != nil {
			return nil, err
		}
		err = se.BinaryWrite(buf, postings.Positions)
		if err != nil {
			return nil, err
		}
		postings = postings.Next
	}

	return buf, nil
}

// CreateNewPostingsList 创建倒排索引
func CreateNewPostingsList(docId int64) *PostingsList {
	p := new(PostingsList)
	p.DocId = docId
	p.PositionCount = 1
	p.Positions = make([]int64, 0)
	return p
}
