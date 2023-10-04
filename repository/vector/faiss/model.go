package faiss

import (
	"context"
	"time"

	"github.com/CocaineCong/tangseng/config"
)

type FaissModel struct {
	name   string
	client *VectorClient
}

func NewFaissModel(name string) *FaissModel {
	return &FaissModel{
		name: name,
	}
}

func (m *FaissModel) Init(ctx context.Context) (err error) {
	vConfig := config.Conf.Vector
	client, err := NewVectorClient(ctx, vConfig.ServerAddress, time.Millisecond*time.Duration(vConfig.Timeout))
	if err != nil {
		return
	}
	m.client = client

	return
}

func (m *FaissModel) Run(data interface{}) (interface{}, error) {
	return m.client.Search(data)
}
