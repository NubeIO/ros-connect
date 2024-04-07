package model

import (
	"context"
	"fmt"
)

type Response struct {
	Code int         `json:"code"` // 0 = ok -1 == fail
	Msg  *string     `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(ctx context.Context, data interface{}) *Response {
	return SuccessResponse(data)
}

func Fail(ctx context.Context, msg interface{}) *Response {
	return FailResponse(msg)
}

func SuccessResponse(data interface{}) *Response {
	return &Response{
		Code: 0,
		Data: data,
	}
}

func FailResponse(msg interface{}) *Response {
	_msg := fmt.Sprintf("%s", msg)
	return &Response{
		Code: -1,
		Msg:  &_msg,
	}
}
