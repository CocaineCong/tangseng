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

func TestRetryOptionRetry(t *testing.T) {
	ctx := context.Background()
	func_ := func(ctx context.Context, req interface{}) (interface{}, bool, error) {
		return doSomethingFunc(ctx, req)
	}
	r := NewRetryOption(ctx, DefaultGapTime, DefaultRetryCount, func_)
	resp, _ := r.Retry(ctx, nil)
	fmt.Println(resp)
}
