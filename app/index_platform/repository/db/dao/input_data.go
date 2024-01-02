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

package dao

import (
	"context"

	"github.com/pkg/errors"

	"gorm.io/gorm"

	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/repository/mysql/db"
	"github.com/CocaineCong/tangseng/repository/mysql/model"
)

type InputDataDao struct {
	*gorm.DB
}

func NewInputDataDao(ctx context.Context) *InputDataDao {
	return &InputDataDao{db.NewDBClient(ctx)}
}

func (d *InputDataDao) CreateInputData(in *model.InputData) (err error) {
	err = d.DB.Model(&model.InputData{}).Create(&in).Error
	if err != nil {
		return errors.Wrap(err, "failed to create inputData")
	}
	return
}

func (d *InputDataDao) BatchCreateInputData(in []*model.InputData) (err error) {
	err = d.DB.Model(&model.InputData{}).CreateInBatches(&in, consts.BatchCreateSize).Error
	if err != nil {
		return errors.Wrap(err, "failed to batch create inputData")
	}
	return
}

func (d *InputDataDao) ListInputData() (in []*model.InputData, err error) {
	err = d.DB.Model(&model.InputData{}).Where("is_index = ?", false).
		Find(&in).Error
	if err != nil {
		err = errors.Wrap(err, "failed to query inputData")
	}
	return
}

func (d *InputDataDao) UpdateInputDataByIds(ids []int64) (err error) {
	err = d.DB.Model(&model.InputData{}).Where("id IN ?", ids).
		Update("is_index", true).Error
	if err != nil {
		err = errors.Wrap(err, "failed to update inputData")
	}
	return
}
