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

package recall

import (
	"context"

	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/types"
)

// SearchRecall 词条回归
func SearchRecall(ctx context.Context, query string) (res []*types.SearchItem, err error) {
	recallService := NewRecall()
	res, err = recallService.Search(ctx, query)
	if err != nil {
		log.LogrusObj.Errorf("SearchRecall-NewRecallServ:%+v", err)
		return
	}

	return
}

// SearchQuery 词条联想
func SearchQuery(query string) (res []string, err error) {
	recallService := NewRecall()
	res, err = recallService.SearchQueryWord(query)
	if err != nil {
		log.LogrusObj.Errorf("SearchRecall-NewRecallServ:%+v", err)
		return
	}

	return
}
