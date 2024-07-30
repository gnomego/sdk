package errors

import (
	errors2 "errors"
	"fmt"
	"reflect"
)

type SystemError struct {
	Code    string                  `json:"code" yaml:"code"`
	Message string                  `json:"message" yaml:"message"`
	Data    *map[string]interface{} `json:"data,omitempty" yaml:"data,omitempty"`
	Link    *string                 `json:"link,omitempty" yaml:"link,omitempty"`
	Target  *string                 `json:"target,omitempty" yaml:"target,omitempty"`
	Details *[]SystemError          `json:"details,omitempty" yaml:"details,omitempty"`
}

func New(message string) *SystemError {
	return &SystemError{
		Code:    "error",
		Message: message,
	}
}

func Newf(format string, args ...interface{}) *SystemError {
	return &SystemError{
		Code:    "error",
		Message: fmt.Sprintf(format, args...),
	}
}

func NewWithCode(code, message string) *SystemError {
	return &SystemError{
		Code:    code,
		Message: message,
	}
}

func NewWithCodef(code, format string, args ...interface{}) *SystemError {
	return &SystemError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewNotFound(message string) *SystemError {
	return &SystemError{
		Code:    "not_found",
		Message: message,
	}
}

func NewUserUnauthorized(message string) *SystemError {
	return &SystemError{
		Code:    "user_unauthorized",
		Message: message,
	}
}

func NewUserUnauthorizedf(format string, args ...interface{}) *SystemError {
	return &SystemError{
		Code:    "user_unauthorized",
		Message: fmt.Sprintf(format, args...),
	}
}

func Join(errs ...error) error {
	return errors2.Join(errs...)
}

func As(err error, target any) bool {
	return errors2.As(err, target)
}

func Is(err, target error) bool {
	return errors2.Is(err, target)
}

func (e *SystemError) Is(target error) bool {
	if e == nil && target == nil {
		return true
	}

	t, ok := target.(*SystemError)
	if !ok {
		return false
	}
	return (e.Code == t.Code || t.Code == "")
}

func (e *SystemError) SetCode(code string) *SystemError {
	e.Code = code
	return e
}

func (e *SystemError) SetData(data map[string]interface{}) *SystemError {
	e.Data = &data
	return e
}

func (e *SystemError) SetLink(link string) *SystemError {
	e.Link = &link
	return e
}

func (e *SystemError) SetTarget(target string) *SystemError {
	e.Target = &target
	return e
}

func (e *SystemError) AddDetail(detail *SystemError) *SystemError {
	if e.Details == nil {
		e.Details = &[]SystemError{}
	}

	*e.Details = append(*e.Details, *detail)
	return e
}

func (e *SystemError) AddDetails(detail ...*SystemError) *SystemError {
	if e.Details == nil {
		e.Details = &[]SystemError{}
	}

	for _, d := range detail {
		*e.Details = append(*e.Details, *d)
	}
	return e
}

func (e *SystemError) SetDetails(details []SystemError) *SystemError {
	e.Details = &details
	return e
}

func (e *SystemError) Error() string {
	return e.Message
}

func (e *SystemError) Map(err error) {
	if err == nil || e == nil {
		return
	}

	if e2, ok := err.(*SystemError); ok {
		e.Code = e2.Code
		e.Message = e2.Message
		e.Data = e2.Data
		return
	}

	e.Code = reflect.TypeOf(err).Name()
	e.Message = err.Error()
}
