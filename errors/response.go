package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// ResponseError 定义响应错误
type ResponseError struct {
	Message string // 错误消息
	ERR     error  // 响应错误
}

func (r *ResponseError) Error() string {
	if r.ERR != nil {
		if r.Message != "" {
			return r.Message
		}
		return r.ERR.Error()
	}
	return r.Message
}

func NewError(err error) error {
	res := &ResponseError{
		ERR: errors.WithStack(err),
	}
	return res
}

func NewMsg(msg string, args ...interface{}) error {
	res := &ResponseError{
		Message: fmt.Sprintf(msg, args...),
	}
	return res
}

func NewErrorMsg(err error, msg string, args ...interface{}) error {
	res := &ResponseError{
		ERR:     err,
		Message: fmt.Sprintf(msg, args...),
	}
	return res
}

// UnWrapResponse 解包响应错误
func UnWrapResponse(err error) *ResponseError {
	if v, ok := err.(*ResponseError); ok {
		return v
	}
	return nil
}
