package dao

import (
	"context"
	"github.com/pkg/errors"

	"gorm.io/gorm"

	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/repository/mysql/db"
	"github.com/CocaineCong/tangseng/repository/mysql/model"
	"github.com/CocaineCong/tangseng/types"
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
		return errors.Wrap(err, "failed to batch create inputData ")
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

// ListInputDataByDocIds 根据传进来的 docs id 获取所有的信息
func (d *InputDataDao) ListInputDataByDocIds(docIds []int64) (in []*types.SearchItem, err error) {
	err = d.DB.Model(&model.InputData{}).
		Where("doc_id IN ?", docIds).
		Select("doc_id," +
			"title," +
			"body AS content," +
			"url," +
			"score AS content_score").
		Find(&in).Error
	if err != nil {
		err = errors.Wrap(err, "failed to query inputData by docIds")
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
