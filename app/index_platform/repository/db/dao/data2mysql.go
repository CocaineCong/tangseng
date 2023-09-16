package dao

import (
	"context"
	"sync"
	"time"

	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/repository/mysql/model"
	"github.com/CocaineCong/tangseng/types"
)

type MySqlDirectUpload struct {
	ctx     context.Context
	doneCtx context.Context

	data   []*model.InputData // 数据
	upData []*model.InputData // 上传的数据
	wLock  *sync.Mutex
	upLock *sync.RWMutex
	task   *types.Task

	done func()
}

// NewMySqlDirectUpload 新建一个上传的对象
func NewMySqlDirectUpload(ctx context.Context, task *types.Task) *MySqlDirectUpload {
	ctx, done := context.WithCancel(ctx)

	directUpload := &MySqlDirectUpload{
		ctx:    ctx,
		data:   make([]*model.InputData, 0, 1e5),
		upData: make([]*model.InputData, 0),
		wLock:  &sync.Mutex{},
		upLock: &sync.RWMutex{},
		task:   task,
		done:   done,
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.LogrusObj.Errorf("消费出现错误")
			}
		}()
		directUpload.consume()
	}()

	return directUpload
}

func (d *MySqlDirectUpload) consume() {
	gapTime := 5 * time.Minute
	for {
		select {
		case <-time.After(gapTime):
			log.LogrusObj.Infof("direct upload")
			d.StreamUpload()
			// case <- d.done:
			// 	d.StreamUpload()
		}
	}

}

func (d *MySqlDirectUpload) StreamUpload() (count int, err error) {
	// 写数据库
	d.wLock.Lock()
	if len(d.data) == 0 {
		d.upData = d.data
	} else {
		d.upData = append(d.upData, d.data...)
	}
	d.data = make([]*model.InputData, 0)
	count = len(d.upData)
	d.wLock.Unlock()

	// 开始上报数据
	d.upLock.Lock()
	defer d.upLock.Unlock()

	err = NewInputDataDao(d.ctx).BatchCreateInputData(d.upData)
	if err != nil {
		log.LogrusObj.Error("BatchCreateInputData ", err)
	}

	// 重制 updata
	d.wLock.Lock()
	d.upData = make([]*model.InputData, 0)
	d.wLock.Unlock()

	return
}

func (d *MySqlDirectUpload) Finish() {
	d.done()
}

func (d *MySqlDirectUpload) Push(records *model.InputData) int {
	d.wLock.Lock()
	defer d.wLock.Unlock()
	d.data = append(d.data, records)
	log.LogrusObj.Infof("direct_upload push bi_table:%s", d.task.BiTable)

	return len(d.data)
}
