package helper

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Pagination struct {
	Page       int32 `json:"page"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
	TotalCount int32 `json:"total_count"`
}

func PaginationNew(context *fiber.Ctx) *Pagination {
	pagination := &Pagination{}

	page, err := strconv.Atoi(context.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pagination.Page = int32(page)

	limit, err := strconv.Atoi(context.Query("limit"))
	if err != nil || limit < -1 {
		limit = 10
	}
	pagination.Limit = int32(limit)

	offset := (page - 1) * limit
	pagination.Offset = int32(offset)
	pagination.TotalCount = 0

	return pagination
}
