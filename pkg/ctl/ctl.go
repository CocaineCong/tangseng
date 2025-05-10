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

package ctl

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CocaineCong/tangseng/pkg/tracing"

	e2 "github.com/CocaineCong/tangseng/consts/e"
)

const (
	defaultSuccessMsg = "操作成功"
)

// Response 基础序列化器
type Response struct {
	Status    int         `json:"status"`
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg"`
	Error     string      `json:"error"`
	Timestamp int64       `json:"timestamps"`
	TraceID   string      `json:"traceID"`
}

// RespSuccess 带data成功返回
func RespSuccess(ctx *gin.Context, data interface{}, code ...int) *Response {
	status := e2.SUCCESS
	if code != nil {
		status = code[0]
	}

	if data == nil {
		data = defaultSuccessMsg
	}

	r := &Response{
		Status:    status,
		Data:      data,
		Msg:       e2.GetMsg(status),
		Timestamp: time.Now().Unix(),
		TraceID:   tracing.GetTraceID(ctx),
	}

	return r
}

func RespError(ctx *gin.Context, err error, data string, code ...int) *Response {
	status := e2.ERROR
	if code != nil {
		status = code[0]
	}

	r := &Response{
		Status:    status,
		Data:      data,
		Msg:       e2.GetMsg(status),
		Error:     err.Error(),
		Timestamp: time.Now().Unix(),
		TraceID:   tracing.GetTraceID(ctx),
	}

	return r
}
