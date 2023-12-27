package faiss

import (
	"context"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, "failed to create new vector client")
	}
	m.client = client

	return
}

func (m *FaissModel) Run(data interface{}) (resp interface{}, err error) {
	resp, err = m.client.Search(data)
	if err != nil {
		err = errors.WithMessage(err, "search error")
	}
	return
}
