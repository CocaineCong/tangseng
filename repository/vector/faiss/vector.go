package faiss

import (
	"context"
	"github.com/pkg/errors"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/CocaineCong/tangseng/idl/pb/vector_retrieval"
)

type VectorClient struct {
	ctx           context.Context
	ServerAddress string
	Timeout       time.Duration
	VectorClient  pb.VectorRetrievalClient
}

func NewVectorClient(ctx context.Context, address string, timeout time.Duration) (client *VectorClient, err error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return client, errors.Wrap(err, "failed to connect with grpc")
	}

	client = &VectorClient{
		ctx:           ctx,
		ServerAddress: address,
		Timeout:       timeout,
		VectorClient:  pb.NewVectorRetrievalClient(conn),
	}

	return
}

func (c *VectorClient) Search(req interface{}) (resp *pb.VectorResp, err error) {
	request, ok := req.(*pb.VectorReq)
	if !ok {
		return resp, errors.Wrap(errors.New("unexpected request type"), "failed to assert req as pb.VectorReq")
	}
	ctx, cancl := context.WithTimeout(c.ctx, c.Timeout)
	defer cancl()
	resp, err = c.VectorClient.Search(ctx, request)
	if err != nil {
		err = errors.Wrap(err, "failed to VectorClient-search")
	}

	return
}
