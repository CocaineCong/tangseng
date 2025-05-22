// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package rpc

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/pkg/errors"

	pb "github.com/CocaineCong/tangseng/idl/pb/index_platform"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// BuildIndex 建立索引的RPC调用
func BuildIndex(ctx context.Context, req *pb.BuildIndexReq) (resp *pb.BuildIndexResp, err error) {
	resp, err = IndexPlatformClient.BuildIndexService(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "IndexPlatformClient.BuildIndexService err")
		return
	}

	return
}

func UploadByStream(ctx context.Context, req *pb.BuildIndexReq, file multipart.File, fileSize int64) (resp *pb.UploadResponse, err error) {
	stream, err := IndexPlatformClient.UploadFile(ctx)
	if err != nil {
		err = errors.WithMessage(err, "IndexPlatformClient.UploadStream err")
		return
	}
	buf := make([]byte, 1024*1024) // 1MB chunks
	for {
		n, errx := file.Read(buf)
		if errx == io.EOF {
			break
		}
		if err = stream.Send(&pb.FileChunk{
			Content: buf[:n],
		}); err != nil {
			log.LogrusObj.Error("stream.Send", err)
			return
		}
	}
	resp, err = stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		return nil, err
	}

	return resp, nil
}
