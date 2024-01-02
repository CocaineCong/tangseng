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

package retry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const (
	DefaultRetryCount = 3               // 默认的超时重试次数
	DefaultGapTime    = 3 * time.Second // 默认的超时时间间隔
)

var instance *RetryOption
var once sync.Once

type DelayRetryFunc func(context.Context, interface{}) (interface{}, bool, error)

type RetryOption struct {
	GapTime    time.Duration  // 重试间隔时间
	RetryCount int            // 重试次数
	RetryFunc  DelayRetryFunc // 重试函数

	ctx context.Context
}

func NewRetryOption(ctx context.Context, gapTime time.Duration, retryCount int, func_ DelayRetryFunc) *RetryOption {
	once.Do(func() {
		instance = &RetryOption{
			GapTime:    gapTime,
			RetryCount: retryCount,
			RetryFunc:  func_,
			ctx:        ctx,
		}
	})

	return instance
}

func (r *RetryOption) Retry(ctx context.Context, req interface{}) (resp interface{}, err error) {
	if r.RetryFunc == nil {
		return resp, errors.Wrap(errors.New("RetryFunc is nil"), "failed to retry")
	}

	for i := 0; i < r.RetryCount; i++ {
		res, needRetry, errx := r.RetryFunc(ctx, req)
		if needRetry || errx != nil {
			err = errx
			fmt.Printf("retry count %d ", i+1)
			time.Sleep(r.GapTime)
			continue
		}
		resp = res
		break
	}

	return resp, errors.Wrap(err, "needRetry or errx is not nil")
}
