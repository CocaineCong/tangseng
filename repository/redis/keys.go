package redis

import (
	"fmt"
	"time"
)

const (
	// InvertedIndexDbPathKey 存放倒排索引的路径
	InvertedIndexDbPathKey = "index_platform:inverted_index"
	// TireTreeDbPathKey 存放tire tree树的路径
	TireTreeDbPathKey = "index_platform:tire_tree"

	// QueryTokenDocIds 搜索过的token的doc ids query:term --> doc ids
	QueryTokenDocIds = "query_doc_id:%s"
	// UserQueryToken 用户搜索过的token query:user_id --> term
	UserQueryToken = "query_token:%d"
)

const (
	QueryTokenDocIdsDefaultTimeout = 10 * time.Minute
)

func getQueryTokenDocIdsKey(term string) string {
	return fmt.Sprintf(QueryTokenDocIds, term)
}

func getUserQueryTokenKey(userId int64) string {
	return fmt.Sprintf(UserQueryToken, userId)
}
