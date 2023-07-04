package handler

import (
	"sync"

	pb "github.com/CocaineCong/tangseng/idl/pb/search_engine"
)

var SearchEngineSrvIns *SearchEngineSrv
var SearchEngineSrvOnce sync.Once

type SearchEngineSrv struct {
	pb.UnimplementedSearchEngineServiceServer
}

func GetSearchEngineSrv() *SearchEngineSrv {
	SearchEngineSrvOnce.Do(func() {
		SearchEngineSrvIns = &SearchEngineSrv{}
	})
	return SearchEngineSrvIns
}
