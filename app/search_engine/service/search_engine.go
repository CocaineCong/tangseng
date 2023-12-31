package service

import (
	"context"
	"sync"

	"github.com/pingcap/errors"

	"github.com/CocaineCong/tangseng/app/search_engine/service/recall"
	"github.com/CocaineCong/tangseng/consts/e"
	pb "github.com/CocaineCong/tangseng/idl/pb/search_engine"
	"github.com/CocaineCong/tangseng/types"
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
	sResult, err := recall.SearchRecall(ctx, query)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = err.Error()
		err = errors.WithMessage(err, "SearchEngineSearch-recall.SearchRecall error")
		return
	}

	resp.SearchEngineInfoList, err = BuildSearchEngineResp(sResult)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = err.Error()
		err = errors.WithMessage(err, "SearchEngineSearch-BuildSearchEngineResp error")
		return
	}
	resp.Count = int64(len(sResult))

	return
}

// WordAssociation 词语联想
func (s *SearchEngineSrv) WordAssociation(ctx context.Context, req *pb.SearchEngineRequest) (resp *pb.WordAssociationResponse, err error) {
	resp = new(pb.WordAssociationResponse)
	resp.Code = e.SUCCESS
	query := req.Query
	associationList, err := recall.SearchQuery(query)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = err.Error()
		err = errors.WithMessage(err, "SearchEngineSearch-WordAssociation error")
		return
	}
	resp.WordAssociationList = associationList

	return
}

func BuildSearchEngineResp(item []*types.SearchItem) (resp []*pb.SearchEngineList, err error) {
	resp = make([]*pb.SearchEngineList, 0)
	for _, v := range item {
		resp = append(resp, &pb.SearchEngineList{
			UrlId: v.DocId,
			Desc:  v.Content,
			Score: float32(v.Score),
		})
	}

	return
}
