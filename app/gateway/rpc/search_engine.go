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

	"github.com/pkg/errors"

	"github.com/CocaineCong/tangseng/consts/e"
	pb "github.com/CocaineCong/tangseng/idl/pb/search_engine"
)

func SearchEngineSearch(ctx context.Context, req *pb.SearchEngineRequest) (r *pb.SearchEngineResponse, err error) {
	r, err = SearchEngineClient.SearchEngineSearch(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "SearchEngineClient.SearchEngineSearch error")
		return
	}

	if r.Code != e.SUCCESS {
		err = errors.Wrap(errors.New(r.Msg), "r.Code is unsuccessful")
		return
	}

	return
}

func WordAssociation(ctx context.Context, req *pb.SearchEngineRequest) (r *pb.WordAssociationResponse, err error) {
	r, err = SearchEngineClient.WordAssociation(ctx, req)
	if err != nil {
		err = errors.WithMessage(err, "SearchEngineClient.WordAssociation error")
		return
	}

	if r.Code != e.SUCCESS {
		err = errors.Wrap(errors.New(r.Msg), "r.Code is unsuccessful")
		return
	}

	return
}
