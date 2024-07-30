package monad

import (
	"github.com/gnomego/sdk/errors"
)

const (
	isNone = int8(0)
	isSome = int8(1)
)

type Option[T any] struct {
	value  *T
	status int8
}

func Some[T any](value *T) *Option[T] {
	return &Option[T]{
		status: isSome,
		value:  value,
	}
}

func None[T any]() *Option[T] {
	return &Option[T]{
		status: isNone,
	}
}

func MapOption[T, U any](o *Option[T], f func(*T) *U) *Option[U] {
	if o.IsNone() {
		return None[U]()
	}
	return Some[U](f(o.value))
}

func ZipOption[T, U any](o1 *Option[T], o2 *Option[U], f func(*T, *U) *V) *Option[Tupple[T, U]] {
	if o1.IsNone() || o2.IsNone() {
		return None[Tupple[T, U]]()
	}

	return Some[Tupple[T, U]](NewTupple(o1.value, o2.value))
}

func (o *Option[T]) IsSome() bool {
	return o.status == isSome
}

func (o *Option[T]) IsNone() bool {
	return o.status == isNone
}

func (o *Option[T]) And(opt *Option[T]) *Option[T] {
	if o.IsSome() {
		return opt
	}
	return o
}

func (o *Option[T]) AndThen(f func(*T) *Option[T]) *Option[T] {
	if o.IsNone() {
		return None[T]()
	}
	return f(o.value)
}

func (o *Option[T]) Or(opt *Option[T]) *Option[T] {
	if o.IsSome() {
		return o
	}
	return opt
}

func (o *Option[T]) OrElse(f func() *Option[T]) *Option[T] {
	if o.IsSome() {
		return o
	}
	return f()
}

func (o *Option[T]) Expect(msg string) *T {
	if o.IsNone() {
		panic(msg)
	}
	return o.value
}

func (o *Option[T]) Inspect(f func(*T)) *Option[T] {
	if o.IsSome() {
		f(o.value)
	}
	return o
}

func (o *Option[T]) Value() *T {
	return o.value
}

func (o *Option[T]) ToResult() *Result[T] {
	if o.IsNone() {
		return Error[T](errors.New("unwrap a None value"))
	}

	return Ok[T](o.value)
}

func (o *Option[T]) ToResultError(err error) *Result[T] {
	if o.IsNone() {
		return Error[T](err)
	}

	return Ok[T](o.value)
}

func (o *Option[T]) ToResponse() *Response[T] {
	if o.IsNone() {
		return ErrorResponse[T](errors.New("unwrap a None value"))
	}

	return OkResponse[T](o.value)
}

func (o *Option[T]) ToResponseWithError(err *errors.SystemError) *Response[T] {
	if o.IsNone() {
		return ErrorResponse[T](err)
	}

	return OkResponse[T](o.value)
}

func (o *Option[T]) Unwrap() *T {
	if o.IsNone() {
		panic("unwrap a None value")
	}
	return o.value
}

func (o *Option[T]) UnwrapOr(def *T) *T {
	if o.IsSome() {
		return o.value
	}
	return def
}

func (o *Option[T]) UnwrapOrElse(f func() *T) *T {
	if o.IsSome() {
		return o.value
	}
	return f()
}
