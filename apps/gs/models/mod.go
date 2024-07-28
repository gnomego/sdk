package models

const (
	STATUS_ACTIVE          = 1
	STATUS_INACTIVE        = 0
	STATUS_PASSWORD_LOCKED = 2
	STATUS_BANNED          = 3
	STATUS_DELETED         = 4
)

type Page interface {
	Page() int
	Size() int
}

type PageRequest struct {
	Page   int      `json:"page" form:"page"`
	Size   int      `json:"size" form:"size"`
	Expand []string `json:"$expand" form:"expand"`
	Filter []string `json:"$filter" form:"filter"`
}
