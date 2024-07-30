package monad

import (
	"encoding/json"

	"github.com/gnomego/sdk/errors"
)

type PagedResponse[T any] struct {
	items *[]T
	ok    bool
	err   *errors.SystemError
	page  *int32
	size  *int64
	total *int64
}

func OkPagedResponse[T any](items *[]T) *PagedResponse[T] {
	p := int32(1)
	return &PagedResponse[T]{
		items: items,
		ok:    true,
		page:  &p,
	}
}

func ErrorPagedResponse[T any](err *errors.SystemError) *PagedResponse[T] {
	return &PagedResponse[T]{
		ok:  false,
		err: err,
	}
}

func (r *PagedResponse[T]) Ok() bool {
	return r.ok
}

func (r *PagedResponse[T]) Items() *[]T {
	return r.items
}

func (r *PagedResponse[T]) Error() *errors.SystemError {
	return r.err
}

func (r *PagedResponse[T]) Page() *int32 {
	return r.page
}

func (r *PagedResponse[T]) Size() *int64 {
	return r.size
}

func (r *PagedResponse[T]) Total() *int64 {
	return r.total
}

func (r *PagedResponse[T]) SetPagination(page int32, size int64, total int64) *PagedResponse[T] {
	if page < 1 {
		page = 1
	}

	if size < 1 {
		size = 10
	}

	if total < 0 {
		if r.items != nil {
			total = int64(len(*r.items))
		} else {
			total = 0
		}
	}

	r.page = &page
	r.size = &size
	r.total = &total
	return r
}

func (r *PagedResponse[T]) SetPage(page int32) *PagedResponse[T] {
	if page < 1 {
		page = 1
	}

	r.page = &page
	return r
}

func (r *PagedResponse[T]) SetSize(size int64) *PagedResponse[T] {
	if size < 1 {
		size = 10
	}

	r.size = &size
	return r
}

func (r *PagedResponse[T]) SetTotal(total int64) *PagedResponse[T] {
	if total < 0 {
		if r.items != nil {
			total = int64(len(*r.items))
		} else {
			total = 0
		}
	}

	r.total = &total
	return r
}

func (r *PagedResponse[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Ok    bool                `json:"ok"`
		Items *[]T                `json:"values,omitempty"`
		Error *errors.SystemError `json:"error,omitempty"`
		Page  *int32              `json:"page,omitempty"`
		Size  *int64              `json:"size,omitempty"`
		Total *int64              `json:"total,omitempty"`
	}{
		Ok:    r.Ok(),
		Items: r.Items(),
		Error: r.Error(),
		Page:  r.Page(),
		Size:  r.Size(),
		Total: r.Total(),
	})
}

func (r *PagedResponse[T]) UnmarshalJSON(data []byte) error {
	var v struct {
		Ok    bool                `json:"ok"`
		Items *[]T                `json:"values,omitempty"`
		Error *errors.SystemError `json:"error,omitempty"`
		Page  *int32              `json:"page,omitempty"`
		Size  *int64              `json:"size,omitempty"`
		Total *int64              `json:"total,omitempty"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	r.ok = v.Ok
	r.items = v.Items
	r.err = v.Error
	r.page = v.Page
	r.size = v.Size
	r.total = v.Total

	return nil
}
