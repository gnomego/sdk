package xgin

import "fmt"

type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Target  *string     `json:"target,omitempty"`
	Details []ErrorInfo `json:"details,omitempty"`
}

func NewErrorInfo(code, message string) *ErrorInfo {
	return &ErrorInfo{
		Code:    code,
		Message: message,
	}
}

func (e *ErrorInfo) SetTarget(target string) *ErrorInfo {
	e.Target = &target
	return e
}

func (e *ErrorInfo) AddDetail(detail *ErrorInfo) *ErrorInfo {
	e.Details = append(e.Details, *detail)
	return e
}

func (e *ErrorInfo) Error() string {
	return e.Message
}

func ErrorInfoF(code, format string, args ...interface{}) *ErrorInfo {
	return NewErrorInfo(code, fmt.Sprintf(format, args...))
}

func (e *ErrorInfo) Map(err error) {
	if err == nil || e == nil {
		return
	}

	if e2, ok := err.(*ErrorInfo); ok {
		e.Code = e2.Code
		e.Message = e2.Message
		e.Target = e2.Target
		e.Details = e2.Details
		return
	}

	e.Code = "error"
	e.Message = err.Error()
}

func (e *ErrorInfo) MapWithCode(code string, err error) {
	if err == nil || e == nil {
		return
	}

	if e2, ok := err.(*ErrorInfo); ok {
		e.Code = code
		e.Message = e2.Message
		e.Target = e2.Target
		e.Details = e2.Details
		return
	}

	e.Code = code
	e.Message = err.Error()
}

func (e *ErrorInfo) MapWithTarget(code string, target string, err error) {
	if err == nil || e == nil {
		return
	}

	if e2, ok := err.(*ErrorInfo); ok {
		e.Code = code
		e.Message = e2.Message
		e.Target = &target
		e.Details = e2.Details
		return
	}

	e.Code = code
	e.Message = err.Error()
	e.Target = &target
}
