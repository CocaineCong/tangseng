package milvus

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"

	"github.com/CocaineCong/tangseng/config"
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
	milvusClient, err := client.NewGrpcClient(ctx, fmt.Sprintf("%s:%s", mConfig.Host, mConfig.Port))
	if err != nil {
		return errors.Wrap(err, "failed to create new grpc client")
	}
	m.client = milvusClient

	return
}

func (m *MilvusModel) Search(req interface{}) (resp interface{}, err error) {
	request, ok := req.(*MilvusRequest)
	if !ok {
		err = errors.Wrap(errors.New("unexpected request type"), "failed to assert req as MilvusRequest")
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
