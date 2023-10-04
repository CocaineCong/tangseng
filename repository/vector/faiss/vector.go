package faiss

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/CocaineCong/tangseng/idl/pb/vector_retrieval"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
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
		return
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
		return
	}
	ctx, cancl := context.WithTimeout(c.ctx, c.Timeout)
	defer cancl()
	resp, err = c.VectorClient.Search(ctx, request)
	if err != nil {
		logs.LogrusObj.Errorln("VectorClient-Search ", err)
		return
	}

	return
}
