package handler

import (
	"context"
	"errors"
	inputData "se/internal/inputdata"
	"se/internal/repository"
	"se/internal/service"
	"se/pkg/e"
)

type SearchEngineService struct {
}

func NewSearchEngineService() *SearchEngineService {
	return &SearchEngineService{}
}

func (*SearchEngineService) SearchEngineAdd(ctx context.Context,req *service.SearchEngineRequest) (resp *service.SearchEngineResponse,err error) {
	tableName := req.Table // table
	postKey := req.Key
	data := req.Data
	table := repository.GetTable(tableName)
	resp = new(service.SearchEngineResponse)
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

func (*SearchEngineService) SearchEngineSearch(ctx context.Context,req *service.SearchEngineRequest) (resp *service.SearchEngineResponse,err error) {
	tableName := req.Table
	indexName := req.Key
	valueName := req.Key
	resp = new(service.SearchEngineResponse)
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

func (*SearchEngineService) SearchEngineAllIndex(ctx context.Context,req *service.SearchEngineRequest) (resp *service.SearchEngineResponse,err error) {
	tableName := req.SearchEngineReq.Info
	table := repository.GetTable(tableName)
	table.AllIndex(50)
	resp = new(service.SearchEngineResponse)
	resp.Code = e.SUCCESS
	resp.Msg = e.GetMsg(e.SUCCESS)
	return resp, err
}

func (*SearchEngineService) SearchEngineAllIndexCount(ctx context.Context,req *service.SearchEngineRequest) (resp *service.SearchEngineResponse,err error) {
	tableName := req.SearchEngineReq.Info
	table := repository.GetTable(tableName)
	table.AllIndexCount()
	resp = new(service.SearchEngineResponse)
	resp.Code = e.SUCCESS
	resp.Msg = e.GetMsg(e.SUCCESS)
	return resp, err
}