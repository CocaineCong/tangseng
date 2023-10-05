package redis

import (
	"fmt"
	"time"
)

// InvertedIndexDbPathKeys 存放倒排索引的路径
var InvertedIndexDbPathKeys = []string{InvertedIndexDbPathDayKey,
	InvertedIndexDbPathMonthKey, InvertedIndexDbPathSeasonKey}

// TireTreeDbPathKey 存放tire tree树的路径
var TireTreeDbPathKey = []string{TireTreeDbPathDayKey,
	TireTreeDbPathMonthKey, TireTreeDbPathSeasonKey,
}

const (
	InvertedIndexDbPathDayKey    = "index_platform:inverted_index:day"       // 天纬度
	InvertedIndexDbPathMonthKey  = "index_platform:inverted_index:month:%s"  // 月纬度
	InvertedIndexDbPathSeasonKey = "index_platform:inverted_index:season:%s" // 季纬度

	TireTreeDbPathDayKey    = "index_platform:tire_tree:day"       // 天纬度
	TireTreeDbPathMonthKey  = "index_platform:tire_tree:month:%s"  // 月纬度
	TireTreeDbPathSeasonKey = "index_platform:tire_tree:season:%s" // 季纬度

	// QueryTokenDocIds 搜索过的token的doc ids query:term --> docs ids
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

func getAllDbPaths(key string) string {
	return fmt.Sprintf(key, "*")
}

func GetInvertedIndexDbPathMonthKey(month string) string {
	return fmt.Sprintf(InvertedIndexDbPathMonthKey, month)
}

func GetInvertedIndexDbPathSeasonKey(season string) string {
	return fmt.Sprintf(InvertedIndexDbPathSeasonKey, season)
}
