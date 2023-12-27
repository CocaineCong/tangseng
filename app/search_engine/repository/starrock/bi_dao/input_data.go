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
		err = errors.Wrap(err, "failed to get data")
	}
	return
}
