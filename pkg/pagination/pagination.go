package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Pageable represents pagination parameters
type Pageable struct {
	Page int    `json:"page"`
	Size int    `json:"size"`
	Sort string `json:"sort"`
}

// DefaultPageable creates a default pageable
func DefaultPageable() *Pageable {
	return &Pageable{
		Page: 0,
		Size: 10,
		Sort: "",
	}
}

// NewPageableFromContext extracts pagination params from Gin context
func NewPageableFromContext(c *gin.Context) *Pageable {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	sort := c.DefaultQuery("sort", "")

	if page < 0 {
		page = 0
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	return &Pageable{
		Page: page,
		Size: size,
		Sort: sort,
	}
}

// Offset returns the offset for database queries
func (p *Pageable) Offset() int {
	return p.Page * p.Size
}

// Page represents a paginated result
type Page[T any] struct {
	Content          []T   `json:"content"`
	TotalElements    int64 `json:"totalElements"`
	TotalPages       int   `json:"totalPages"`
	Size             int   `json:"size"`
	Number           int   `json:"number"`
	NumberOfElements int   `json:"numberOfElements"`
	First            bool  `json:"first"`
	Last             bool  `json:"last"`
	Empty            bool  `json:"empty"`
}

// NewPage creates a new Page with the given content and pagination info
func NewPage[T any](content []T, totalElements int64, pageable *Pageable) *Page[T] {
	totalPages := int(totalElements) / pageable.Size
	if int(totalElements)%pageable.Size > 0 {
		totalPages++
	}

	return &Page[T]{
		Content:          content,
		TotalElements:    totalElements,
		TotalPages:       totalPages,
		Size:             pageable.Size,
		Number:           pageable.Page,
		NumberOfElements: len(content),
		First:            pageable.Page == 0,
		Last:             pageable.Page >= totalPages-1,
		Empty:            len(content) == 0,
	}
}
