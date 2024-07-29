package xgin

import (
	"fmt"
	"reflect"

	"github.com/gnomego/apps/gs/einfo"
	"github.com/gnomego/apps/gs/log"
)

type Response[T any] struct {
	Value T                `json:"values,omitempty"`
	Ok    bool             `json:"ok"`
	Error *einfo.ErrorInfo `json:"error,omitempty"`
}

func (r *Response[T]) IsOk() bool {
	return r.Ok
}

func (r *Response[T]) ErrorInfo() *einfo.ErrorInfo {
	return r.Error
}

func (r *Response[T]) Invalid(e *einfo.ErrorInfo) *Response[T] {
	r.Ok = false
	// determine if t is a pointer or not
	// if it is a pointer, then set the value to nil
	// if it is not a pointer, then set the value to the zero value
	// https://stackoverflow.com/a/6395076
	val := reflect.ValueOf(&r.Value).Elem()
	if val.Kind() == reflect.Ptr {
		r.Value = reflect.Zero(val.Type()).Interface().(T)
	}

	r.Error = e
	return r
}

func (r *Response[T]) SetErrorMessage(code string, message string, args ...interface{}) *Response[T] {
	r.Ok = false
	msg := fmt.Sprintf(message, args...)
	r.Error = einfo.NewErrorInfo("error", msg)

	log.Error(fmt.Errorf(msg), msg)

	val := reflect.ValueOf(&r.Value).Elem()
	if val.Kind() == reflect.Ptr {
		r.Value = reflect.Zero(val.Type()).Interface().(T)
	}
	return r
}

func (r *Response[T]) SetError(err error, message string, args ...interface{}) *Response[T] {
	r.Ok = false
	msg := fmt.Sprintf(message, args...)
	r.Error = einfo.NewErrorInfo("error", msg)

	log.Error(err, msg)

	val := reflect.ValueOf(&r.Value).Elem()
	if val.Kind() == reflect.Ptr {
		r.Value = reflect.Zero(val.Type()).Interface().(T)
	}
	return r
}

type ErrorResponse struct {
	Ok    bool             `json:"ok"`
	Error *einfo.ErrorInfo `json:"error,omitempty"`
}

func (r *ErrorResponse) IsOk() bool {
	return r.Ok
}

func (r *ErrorResponse) ErrorInfo() *einfo.ErrorInfo {
	return r.Error
}

type ResponseInfo interface {
	IsOk() bool
	ErrorInfo() *einfo.ErrorInfo
}

func NewErrorResponse(code string, message string) *ErrorResponse {
	return &ErrorResponse{
		Ok:    false,
		Error: &einfo.ErrorInfo{Code: code, Message: message},
	}
}
func NewTargetedErrorResponse(target string, code string, message string) *ErrorResponse {
	return &ErrorResponse{
		Ok:    false,
		Error: &einfo.ErrorInfo{Code: code, Message: message, Target: &target},
	}
}

type PagedResponse[T any] struct {
	Values []T              `json:"values,omitempty"`
	Ok     bool             `json:"ok"`
	Error  *einfo.ErrorInfo `json:"error,omitempty"`
	Page   int              `json:"page,omitempty"`
	Size   int              `json:"size,omitempty"`
	Total  int              `json:"total,omitempty"`
}

func (r *PagedResponse[T]) IsOk() bool {
	return r.Ok
}

func (r *PagedResponse[T]) ErrorInfo() *einfo.ErrorInfo {
	return r.Error
}
