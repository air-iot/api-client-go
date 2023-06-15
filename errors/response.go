package errors

import (
	"fmt"
	"github.com/air-iot/json"

	"github.com/air-iot/errors"
)

// ResponseError 定义响应错误
type ResponseError struct {
	StatusCode int    `json:"statusCode"` // 错误码
	Code       int    `json:"code"`       // 错误码
	Message    string `json:"message"`    // 错误信息
	Field      string `json:"field"`
	Detail     string `json:"detail"` // 错误详情信息
	ERR        error  `json:"err"`    // 响应错误
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

func ParseBody(statusCode int, body []byte) error {
	res := &ResponseError{}
	if err := json.Unmarshal(body, res); err != nil {
		res.Message = string(body)
		res.StatusCode = statusCode
		return res
	}
	res.StatusCode = statusCode
	return res
}
