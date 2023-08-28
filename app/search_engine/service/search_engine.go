package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/CocaineCong/tangseng/app/search_engine/recall"
	"github.com/CocaineCong/tangseng/app/search_engine/types"
	"github.com/CocaineCong/tangseng/consts/e"
	pb "github.com/CocaineCong/tangseng/idl/pb/search_engine"
	log "github.com/CocaineCong/tangseng/pkg/logger"
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

// SearchEngineSearch 搜索
func (s *SearchEngineSrv) SearchEngineSearch(ctx context.Context, req *pb.SearchEngineRequest) (resp *pb.SearchEngineResponse, err error) {
	resp = new(pb.SearchEngineResponse)
	resp.Code = e.SUCCESS
	query := req.Query
	sResult, err := recall.SearchRecall(query)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = err.Error()
		log.LogrusObj.Error("SearchEngineSearch-recall.SearchRecall", err)
		return
	}

	resp.SearchEngineInfoList, err = BuildSearchEngineResp(sResult)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = err.Error()
		log.LogrusObj.Error("SearchEngineSearch-BuildSearchEngineResp", err)
		return
	}

	return
}

// WordAssociation 词语联想
func (s *SearchEngineSrv) WordAssociation(ctx context.Context, req *pb.SearchEngineRequest) (resp *pb.WordAssociationResponse, err error) {
	resp = new(pb.WordAssociationResponse)
	resp.Code = e.SUCCESS
	query := req.Query
	sResult, err := recall.SearchQuery(query)
	wordAssociationList := make([]string, 0)
	for _, v := range sResult {
		if v != nil {
			wordAssociationList = append(wordAssociationList, v.Value)
		}
	}
	resp.WordAssociationList = wordAssociationList
	fmt.Println(resp.WordAssociationList)
	return
}

func BuildSearchEngineResp(item []*types.SearchItem) (resp []*pb.SearchEngineList, err error) {
	resp = make([]*pb.SearchEngineList, 0)
	for _, v := range item {
		resp = append(resp, &pb.SearchEngineList{
			UrlId: v.DocId,
			Desc:  v.Content,
		})
	}

	return
}
