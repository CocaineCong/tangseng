package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/CocaineCong/tangseng/app/search_engine/repository/db/model"
	"github.com/CocaineCong/tangseng/consts"
)

type InputDataDao struct {
	*gorm.DB
}

func NewInputDataDao(ctx context.Context) *InputDataDao {
	return &InputDataDao{NewDBClient(ctx)}
}

func (d *InputDataDao) CreateInputData(in *model.InputData) (err error) {
	return d.DB.Model(&model.InputData{}).Create(&in).Error
}

func (d *InputDataDao) BatchCreateInputData(in []*model.InputData) (err error) {
	return d.DB.Model(&model.InputData{}).CreateInBatches(&in, consts.BatchCreateSize).Error
}
