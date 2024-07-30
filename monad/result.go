package monad

import (
	"encoding/json"
	"fmt"
)

const (
	isError = int8(0)
	isOk    = int8(1)
)

type ResultPair[T any, E any] struct {
	value  *T
	err    *E
	status int8
}

type Result[T any] struct {
	value  *T
	err    error
	status int8
}

func (r *Result[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Ok    bool  `json:"ok"`
		Value *T    `json:"value,omitempty"`
		Error error `json:"error,omitempty"`
	}{
		Ok:    r.IsOk(),
		Value: r.value,
		Error: r.err,
	})
}

func (r *Result[T]) MarshallYAML() ([]byte, error) {
	return json.Marshal(struct {
		Ok    bool  `json:"ok"`
		Value *T    `json:"value,omitempty"`
		Error error `json:"error,omitempty"`
	}{
		Ok:    r.IsOk(),
		Value: r.value,
		Error: r.err,
	})
}

func (r *Result[T]) MarshalText() ([]byte, error) {
	return json.Marshal(struct {
		Ok    bool  `json:"ok"`
		Value *T    `json:"value,omitempty"`
		Error error `json:"error,omitempty"`
	}{
		Ok:    r.IsOk(),
		Value: r.value,
		Error: r.err,
	})
}

func (r *Result[T]) UnmarshalJSON(data []byte) error {
	var v struct {
		Ok    bool  `yaml:"ok"`
		Value *T    `yaml:"value,omitempty"`
		Error error `yaml:"error,omitempty"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v.Ok {
		r.value = v.Value
		r.err = nil
		r.status = isOk
		return nil
	}

	r.value = nil
	r.err = v.Error
	r.status = isError

	return nil
}

func Ok[T any](value *T) *Result[T] {
	return &Result[T]{
		value:  value,
		err:    nil,
		status: isOk,
	}
}

func Error[T any](err error) *Result[T] {
	return &Result[T]{
		value:  nil,
		err:    err,
		status: isError,
	}
}

func MapResult[T any, U any](r Result[T], fn func(*T) *U) *Result[U] {
	if r.status == isError {
		return Error[U](r.err)
	}

	return Ok(fn(r.value))
}

func MapResultOr[T any, U any](r Result[T], def *U, fn func(*T) *U) *Result[U] {
	if r.status == isError {
		return Ok(def)
	}

	return Ok(fn(r.value))
}

func MapResultOrElse[T any, U any](r Result[T], fn func() *U, fn2 func(*T) *U) *Result[U] {
	if r.status == isError {
		return Ok(fn())
	}

	return Ok(fn2(r.value))
}

func MapResultError[T any, U any](r Result[T], fn func(error) error) *Result[U] {
	if r.status == isOk {
		return Error[U](fmt.Errorf("result value is not an error"))
	}

	return Error[U](fn(r.err))
}

func (r *Result[T]) IsOk() bool {
	return r.status == isOk
}

func (r *Result[T]) IsError() bool {
	return r.status == isError
}

func (r *Result[T]) And(other *Result[T]) *Result[T] {
	if r.status == isError {
		return r
	}

	return other
}

func (r *Result[T]) AndThen(fn func(*T) *Result[T]) *Result[T] {
	if r.status == isError {
		return r
	}

	return fn(r.value)
}

func (r *Result[T]) Or(other *Result[T]) *Result[T] {
	if r.status == isOk {
		return r
	}

	return other
}

func (r *Result[T]) OrElse(fn func() *Result[T]) *Result[T] {
	if r.status == isOk {
		return r
	}

	return fn()
}

func (r *Result[T]) Expect(msg string) *T {
	if r.status == isError {
		panic(msg)
	}

	return r.value
}

func (r *Result[T]) ExpectError(msg string) error {
	if r.status == isOk {
		panic(msg)
	}

	return r.err
}

func (r *Result[T]) ExpectErrorf(format string, args ...interface{}) error {
	if r.status == isOk {
		panic(fmt.Sprintf(format, args...))
	}

	return r.err
}

func (r *Result[T]) Expectf(format string, args ...interface{}) *T {
	if r.status == isError {
		panic(fmt.Sprintf(format, args...))
	}

	return r.value
}

func (r *Result[T]) Inspect(fn func(*T)) *Result[T] {
	if r.status == isOk {
		fn(r.value)
	}

	return r
}

func (r *Result[T]) InspectError(fn func(error)) *Result[T] {
	if r.status == isError {
		fn(r.err)
	}

	return r
}

func (r *Result[T]) Unwrap() *T {
	if r.status == isError {
		if r.err == nil {
			panic("result value is nil")
		}

		panic(r.err)
	}

	return r.value
}

func (r *Result[T]) UnwrapOr(def *T) *T {
	if r.status == isError {
		return def
	}

	return r.value
}

func (r *Result[T]) UnwrapOrElse(fn func() *T) *T {
	if r.status == isError {
		return fn()
	}

	return r.value
}

func (r *Result[T]) UnwrapError() error {
	if r.status == isOk {
		panic("result error is nil")
	}

	return r.err
}

func (r *Result[T]) UnwrapErrorOr(def error) error {
	if r.status == isOk {
		return def
	}

	return r.err
}

func (r *Result[T]) UnwrapErrorOrElse(fn func() error) error {
	if r.status == isOk {
		return fn()
	}

	return r.err
}

func (r *Result[T]) Map(fn func(*T) *T) *Result[T] {
	if r.status == isError {
		return Error[T](r.err)
	}

	return Ok(fn(r.value))
}

func (r *Result[T]) MapOr(def *T, fn func(*T) *T) *Result[T] {
	if r.status == isError {
		return Ok(def)
	}

	return Ok(fn(r.value))
}

func (r *Result[T]) MapOrElse(fn func() *T, fn2 func(*T) *T) *Result[T] {
	if r.status == isError {
		return Ok(fn())
	}

	return Ok(fn2(r.value))
}

func (r *Result[T]) MapError(fn func(error) error) *Result[T] {
	if r.status == isOk {
		return Ok(r.value)
	}

	return Error[T](fn(r.err))
}

func (r *Result[T]) Ok() *Option[T] {
	if r.status == isError {
		return None[T]()
	}

	return Some(r.value)
}

func (r *Result[T]) Error() *Option[error] {
	if r.status == isOk {
		return None[error]()
	}

	return Some[error](&r.err)
}
