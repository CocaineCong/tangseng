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

package starrock

import (
	"bufio"
	"bytes"
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/types"
)

type DirectUpload struct {
	ctx     context.Context
	doneCtx context.Context

	data   []*types.Data2Starrocks // 数据
	upData []*types.Data2Starrocks // 上传的数据
	wLock  *sync.Mutex
	upLock *sync.RWMutex
	task   *types.Task

	done func()
}

// NewDirectUpload 新建一个上传的对象
func NewDirectUpload(ctx context.Context, task *types.Task) *DirectUpload {
	ctx, done := context.WithCancel(ctx)

	directUpload := &DirectUpload{
		ctx:    ctx,
		data:   make([]*types.Data2Starrocks, 0, 1e5),
		upData: make([]*types.Data2Starrocks, 0),
		wLock:  &sync.Mutex{},
		upLock: &sync.RWMutex{},
		task:   task,
		done:   done,
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.LogrusObj.Errorf("NewDirectUpload2-消费出现错误 :%+v", err)
			}
		}()
		directUpload.consume()
	}()

	return directUpload
}

func (d *DirectUpload) consume() {
	gapTime := 5 * time.Minute
	for {
		select {
		case <-time.After(gapTime):
			log.LogrusObj.Infof("direct upload")
			_, _ = d.StreamUpload()
		case <-d.doneCtx.Done():
			_, _ = d.StreamUpload()
		}
	}
}

func (d *DirectUpload) StreamUpload() (count int, err error) {
	// 写数据库
	d.wLock.Lock()
	if len(d.data) == 0 {
		d.upData = d.data
	} else {
		d.upData = append(d.upData, d.data...)
	}
	d.data = make([]*types.Data2Starrocks, 0)
	count = len(d.upData)
	d.wLock.Unlock()

	// 开始上报数据
	d.upLock.Lock()
	defer d.upLock.Unlock()

	if len(d.upData) == 0 {
		log.LogrusObj.Infof("finish upload")
	}

	// 构建csv
	rowDelimiter := "@##@" // 分割线，自定义，后面构建文件流传入即可
	sb := &bytes.Buffer{}
	write := bufio.NewWriter(sb)
	for i := 0; i < count; i++ {
		line := strings.Join([]string{
			cast.ToString(d.upData[i].DocId),
			d.upData[i].Title,
			d.upData[i].Desc,
			d.upData[i].Url,
			cast.ToString(d.upData[i].Score),
		}, ",")
		_, err = write.WriteString(line + rowDelimiter)
		if err != nil {
			err = errors.Wrap(err, "failed to write string")
			return
		}
	}
	err = write.Flush()
	if err != nil {
		err = errors.Wrap(err, "failed to flush data")
		return
	}

	// check 机制
	starrocksClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) (err error) {
			v := via[0]
			req.Header = v.Header
			req.Body, err = v.GetBody()
			if err != nil {
				err = errors.Wrapf(err, "starrock woker")
			}
			return err
		},
		Timeout: time.Minute,
	}

	cli := resty.NewWithClient(starrocksClient)
	cli.Header.Add("format", "CSV")
	cli.Header.Add("column_separator", ",")
	cli.Header.Add("row_separator", rowDelimiter)
	cli.Header.Add("columns", strings.Join(d.task.Columns, ","))
	cli.Header.Add("expect", "100-continue")
	cli.Header.Add("Accept", "*/*")

	srConfig := config.Conf.StarRocks
	hp, err := cli.SetDebug(true).R().SetContext(d.ctx).
		SetBasicAuth(srConfig.UserName, srConfig.Password).
		SetPathParams(map[string]string{
			"host":  srConfig.LoadUrl,
			"db":    srConfig.Database,
			"table": d.task.BiTable,
		}).SetBody(sb.Bytes()).SetContentLength(true).
		Put("https://{host}/api/{db}/{table}/_stream_load")
	if err != nil {
		err = errors.Wrap(err, "failed to load stream")
		return
	}

	if hp.StatusCode() != http.StatusOK {
		return
	}

	// 重制 updata
	d.wLock.Lock()
	d.upData = make([]*types.Data2Starrocks, 0)
	d.wLock.Unlock()

	return
}

func (d *DirectUpload) Finish() {
	d.done()
}

func (d *DirectUpload) Push(records *types.Data2Starrocks) int {
	d.wLock.Lock()
	defer d.wLock.Unlock()
	d.data = append(d.data, records)
	log.LogrusObj.Infof("direct_upload push bi_table:%s", d.task.BiTable)

	return len(d.data)
}
