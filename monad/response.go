package monad

import (
	"encoding/json"

	"github.com/gnomego/sdk/errors"
)

type Response[T any] struct {
	value *T
	ok    bool
	err   *errors.SystemError
}

func OkResponse[T any](value *T) *Response[T] {
	return &Response[T]{
		ok:    true,
		value: value,
	}
}

func ErrorResponse[T any](err *errors.SystemError) *Response[T] {
	return &Response[T]{
		ok:  false,
		err: err,
	}
}

func (r *Response[T]) Ok() bool {
	return r.ok
}

func (r *Response[T]) Value() *T {
	return r.value
}

func (r *Response[T]) Error() *errors.SystemError {
	return r.err
}

func (r *Response[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Ok    bool  `json:"ok"`
		Value *T    `json:"values,omitempty"`
		Error error `json:"error,omitempty"`
	}{
		Ok:    r.Ok(),
		Value: r.Value(),
		Error: r.Error(),
	})
}

func (r *Response[T]) UnmarshalJSON(data []byte) error {
	var v struct {
		Ok    bool                `json:"ok"`
		Value *T                  `json:"values,omitempty"`
		Error *errors.SystemError `json:"error,omitempty"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	r.ok = v.Ok
	r.value = v.Value
	r.err = v.Error

	return nil
}
