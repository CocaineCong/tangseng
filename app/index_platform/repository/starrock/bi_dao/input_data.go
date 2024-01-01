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

package bi_dao

import (
	"context"
	"github.com/pkg/errors"

	"gorm.io/gorm"

	"github.com/CocaineCong/tangseng/types"
)

type StarRocksDao struct {
	*gorm.DB
}

func NewStarRocksDao(ctx context.Context) *StarRocksDao {
	return &StarRocksDao{NewDBClient(ctx)}
}

// ListDataRocks 获取用户信息
func (dao *StarRocksDao) ListDataRocks() (r []*types.Data2Starrocks, err error) {
	sql := "SELECT * FROM input_data"
	err = dao.DB.Raw(sql).Find(&r).Error
	if err != nil {
		err = errors.Wrap(err, "failed to find data")
	}
	return
}
