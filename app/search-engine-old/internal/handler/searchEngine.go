package handler

import (
	"context"
	"errors"
	"sync"

	inputData "github.com/CocaineCong/Go-SearchEngine/app/search-engine-old/internal/inputdata"
	"github.com/CocaineCong/Go-SearchEngine/app/search-engine-old/internal/repository"
	pb "github.com/CocaineCong/Go-SearchEngine/idl/pb/search_engine"
	"github.com/CocaineCong/Go-SearchEngine/pkg/e"
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
func (s *SearchEngineSrv) SearchEngineAdd(ctx context.Context, req *pb.SearchEngineRequest) (resp *pb.SearchEngineResponse, err error) {
	tableName := req.Table // table
	postKey := req.Key
	data := req.Data
	table := repository.GetTable(tableName)
	resp = new(pb.SearchEngineResponse)
	inData := &inputData.InputData{
		Key:  postKey,
		Data: data,
	}
	resp.Code = e.SUCCESS
	_, err = table.Insert(inData)
	if err != nil {
		resp.Code = e.ERROR
		return resp, err
	}
	table.Save()
	return resp, nil
}

func (s *SearchEngineSrv) SearchEngineSearch(ctx context.Context, req *pb.SearchEngineRequest) (resp *pb.SearchEngineResponse, err error) {
	tableName := req.Table
	indexName := req.Key
	valueName := req.Key
	resp = new(pb.SearchEngineResponse)
	resp.Code = e.SUCCESS
	if indexName == "" {
		resp.Code = e.ERROR
		resp.Msg = errors.New("必须为查询指定一个索引，用法：/:table?index=index1&value=value1").Error()
		return
	}
	if valueName == "" {
		return
	}
	table := repository.GetTable(tableName)
	if !table.CheckIndexExist(indexName) {
		resp.Code = e.ERROR
		resp.Msg = errors.New("index[" + indexName + "]不存在").Error()
		return
	}
	res, err := table.Search(indexName, valueName)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = err.Error()
		return
	}
	resp.Data = res
	return resp, nil
}

func (s *SearchEngineSrv) SearchEngineAllIndex(ctx context.Context, req *pb.SearchEngineRequest) (resp *pb.SearchEngineResponse, err error) {
	tableName := req.SearchEngineReq.Info
	table := repository.GetTable(tableName)
	table.AllIndex(50)
	resp = new(pb.SearchEngineResponse)
	resp.Code = e.SUCCESS
	resp.Msg = e.GetMsg(e.SUCCESS)
	return resp, err
}

func (s *SearchEngineSrv) SearchEngineAllIndexCount(ctx context.Context, req *pb.SearchEngineRequest) (resp *pb.SearchEngineResponse, err error) {
	tableName := req.SearchEngineReq.Info
	table := repository.GetTable(tableName)
	table.AllIndexCount()
	resp = new(pb.SearchEngineResponse)
	resp.Code = e.SUCCESS
	resp.Msg = e.GetMsg(e.SUCCESS)
	return resp, err
}
