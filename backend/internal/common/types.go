package common

import "time"

// PaginatedResult wraps list responses with pagination metadata.
type PaginatedResult[T any] struct {
	Items    []T   `json:"items"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

// PaginationParams holds standard list query parameters.
type PaginationParams struct {
	Page     int `form:"page" validate:"omitempty,min=1"`
	PageSize int `form:"pageSize" validate:"omitempty,min=1,max=100"`
}

func (p PaginationParams) Offset() int {
	page := p.Page
	if page < 1 {
		page = 1
	}
	size := p.PageSize
	if size < 1 {
		size = 20
	}
	return (page - 1) * size
}

func (p PaginationParams) Limit() int {
	if p.PageSize < 1 {
		return 20
	}
	return p.PageSize
}

// Timestamps is embedded by domain entities.
type Timestamps struct {
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}
