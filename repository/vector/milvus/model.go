package milvus

import (
	"context"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"

	"github.com/CocaineCong/tangseng/config"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

type MilvusModel struct {
	ctx    context.Context
	name   string
	client client.Client
}

func NewMilvusModel(ctx context.Context, name string) *MilvusModel {
	return &MilvusModel{ctx: ctx, name: name}
}

func (m *MilvusModel) Init() (err error) {
	mConfig := config.Conf.Milvus
	ctx, cancel := context.WithTimeout(m.ctx, time.Millisecond*time.Duration(mConfig.Timeout))
	defer cancel()
	milvusClient, err := client.NewGrpcClient(ctx, mConfig.ServerAddress)
	if err != nil {
		logs.LogrusObj.Errorln(err)
		return
	}
	m.client = milvusClient

	return
}

func (m *MilvusModel) Search(req interface{}) (resp interface{}, err error) {
	request, ok := req.(*MilvusRequest)
	if !ok {
		return
	}

	return m.client.Search(
		m.ctx,
		request.CollectionName,
		request.Partitions,
		request.Expr,
		request.OutputFields,
		request.Vectors,
		request.VectorField,
		request.MetricType,
		request.TopK,
		request.SearchParams,
		nil,
	)
}
