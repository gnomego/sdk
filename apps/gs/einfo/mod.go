package einfo

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

func Sprintf(code, format string, args ...interface{}) *ErrorInfo {
	return NewErrorInfo(code, fmt.Sprintf(format, args...))
}

func Map(err error) *ErrorInfo {
	if err == nil {
		return nil
	}
	if e, ok := err.(*ErrorInfo); ok {
		return e
	}
	return NewErrorInfo("error", err.Error())
}

func MapWithCode(code string, err error) *ErrorInfo {
	if err == nil {
		return nil
	}
	if e, ok := err.(*ErrorInfo); ok {
		return e
	}
	return NewErrorInfo(code, err.Error())
}

func MapWithTarget(code string, target string, err error) *ErrorInfo {
	if err == nil {
		return nil
	}
	if e, ok := err.(*ErrorInfo); ok {
		return e
	}
	return NewErrorInfo(code, err.Error()).SetTarget(target)
}
