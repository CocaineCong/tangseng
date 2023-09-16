package segment

import (
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/types"
)

// CreateNewPostingsList 创建倒排索引
func CreateNewPostingsList(docId int64) *types.PostingsList {
	p := new(types.PostingsList)
	p.DocId = docId
	p.PositionCount = 1
	p.Positions = make([]int64, 0)
	return p
}

// MergePostings 合并两个posting
func MergePostings(pa, pb *types.PostingsList) *types.PostingsList {
	ret := new(types.PostingsList)
	p := new(types.PostingsList)
	p = nil
	for pa != nil || pb != nil {
		tmp := new(types.PostingsList)
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
