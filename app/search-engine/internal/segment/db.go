package segment

import (
	"fmt"

	"github.com/CocaineCong/tangseng/app/search-engine/internal/storage"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// InvertedIndexValue 倒排索引
type InvertedIndexValue struct {
	Token         string
	PostingsList  *PostingsList
	DocCount      int64
	PositionCount int64 // 查询使用，写入的时候暂时不用
	TermValues    *storage.TermValue
}

// InvertedIndexHash 倒排hash
type InvertedIndexHash map[string]*InvertedIndexValue

// InitSegmentDb 读取对应segment文件下的db
func InitSegmentDb(segId SegId) (*storage.InvertedDB, *storage.ForwardDB) {
	if segId < 0 {
		log.LogrusObj.Infof("db Init :%d<0", segId)
	}
	log.LogrusObj.Infof("index:[termName:%s,invertedName:%s,forwardName:%s]", termName, invertedName, forwardName)
	termName, invertedName, forwardName = GetDbName(segId)
	forwardDB, err := storage.NewForwardDB(forwardName)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, nil
	}
	return storage.NewInvertedDB(termName, invertedName), forwardDB
}

// CreateNewInvertedIndex 创建倒排索引
func CreateNewInvertedIndex(token string, docCount int64) *InvertedIndexValue {
	p := new(InvertedIndexValue)
	p.DocCount = docCount
	p.Token = token
	p.PositionCount = 0
	p.PostingsList = new(PostingsList)
	return p
}

// GetDbName 获取db的路径+名称
func GetDbName(segId SegId) (string, string, string) {
	termName = fmt.Sprintf("%s%d%s", config.Conf.SeConfig.StoragePath, segId, TermDbSuffix)
	invertedName = fmt.Sprintf("%s%d%s", config.Conf.SeConfig.StoragePath, segId, InvertedDbSuffix)
	forwardName = fmt.Sprintf("%s%d%s", config.Conf.SeConfig.StoragePath, segId, ForwardDbSuffix)
	return termName, invertedName, forwardName
}
