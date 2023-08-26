package retry

import (
	"context"
	"fmt"
	"sync"
	"time"
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
		return
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

	return
}
