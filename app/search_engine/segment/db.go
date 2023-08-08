package segment

import (
	"fmt"

	"github.com/CocaineCong/tangseng/app/search_engine/query"
	storage2 "github.com/CocaineCong/tangseng/app/search_engine/storage"
	"github.com/CocaineCong/tangseng/app/search_engine/types"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// InvertedIndexHash 倒排hash
type InvertedIndexHash map[string]*types.InvertedIndexValue

// InitSegmentDb 读取对应segment文件下的db
func InitSegmentDb(segId SegId) (*storage2.InvertedDB, *storage2.ForwardDB) {
	if segId < 0 {
		log.LogrusObj.Infof("db Init :%d<0", segId)
	}
	log.LogrusObj.Infof("index:[termName:%s,invertedName:%s,forwardName:%s]", termName, invertedName, forwardName)
	termName, invertedName, forwardName = GetDbName(segId)
	forwardDB, err := storage2.NewForwardDB(forwardName)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, nil
	}
	return storage2.NewInvertedDB(termName, invertedName), forwardDB
}

// CreateNewInvertedIndex 创建倒排索引
func CreateNewInvertedIndex(token query.Tokenization, docCount int64) *types.InvertedIndexValue {
	return &types.InvertedIndexValue{ // TODO：优化一下结构
		Token:         token.Token,
		PostingsList:  new(types.PostingsList),
		DocCount:      docCount,
		PositionCount: 0,
		TermValues: &types.TermValue{
			DocCount: docCount,
			Offset:   token.Offset,
			Size:     token.Offset - token.Position,
		},
	}
}

// GetDbName 获取db的路径+名称
func GetDbName(segId SegId) (string, string, string) {
	termName = fmt.Sprintf("%s%d%s", config.Conf.SeConfig.StoragePath, segId, TermDbSuffix)
	invertedName = fmt.Sprintf("%s%d%s", config.Conf.SeConfig.StoragePath, segId, InvertedDbSuffix)
	forwardName = fmt.Sprintf("%s%d%s", config.Conf.SeConfig.StoragePath, segId, ForwardDbSuffix)
	return termName, invertedName, forwardName
}
