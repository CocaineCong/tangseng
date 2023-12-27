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
