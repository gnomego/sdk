package stores

import "github.com/gnomego/apps/gs/einfo"

type Response[T any] struct {
	Value T                `json:"values,omitempty"`
	Ok    bool             `json:"ok"`
	Error *einfo.ErrorInfo `json:"error,omitempty"`
}

type PagedResponse[T any] struct {
	Values []T              `json:"values,omitempty"`
	Ok     bool             `json:"ok"`
	Error  *einfo.ErrorInfo `json:"error,omitempty"`
	Page   int              `json:"page,omitempty"`
	Size   int              `json:"size,omitempty"`
	Total  int              `json:"total,omitempty"`
}
