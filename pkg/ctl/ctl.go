package ctl

import (
	"github.com/gin-gonic/gin"

	e2 "github.com/CocaineCong/tangseng/consts/e"
)

// Response 基础序列化器
type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
}

// RespSuccess 带data成功返回
func RespSuccess(ctx *gin.Context, data interface{}, code ...int) *Response {
	status := e2.SUCCESS
	if code != nil {
		status = code[0]
	}

	if data == nil {
		data = "操作成功"
	}

	r := &Response{
		Status: status,
		Data:   data,
		Msg:    e2.GetMsg(status),
	}

	return r
}

func RespError(ctx *gin.Context, err error, data string, code ...int) *Response {
	status := e2.ERROR
	if code != nil {
		status = code[0]
	}

	r := &Response{
		Status: status,
		Data:   data,
		Msg:    e2.GetMsg(status),
		Error:  err.Error(),
	}

	return r
}
