package retry

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func doSomethingFunc(ctx context.Context, req interface{}) (interface{}, bool, error) {
	a := time.Now().Unix()
	fmt.Println(a)
	if a%2 == 0 {
		return nil, true, nil
	}
	return nil, false, nil
}

func TestRetryOption_Retry(t *testing.T) {
	ctx := context.Background()
	var func_ DelayRetryFunc
	func_ = func(ctx context.Context, req interface{}) (interface{}, bool, error) {
		return doSomethingFunc(ctx, req)
	}
	r := NewRetryOption(ctx, DefaultGapTime, DefaultRetryCount, func_)
	resp, _ := r.Retry(ctx, nil)
	fmt.Println(resp)
}
