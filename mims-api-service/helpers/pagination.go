package helpers

import (
	"math"
	"strconv"

	"gitlab.com/mims-api-service/responses"
)

const (
	DefaultLimit  int64 = 50
	DefaultOffset int64 = 0
)

func Pagination(data interface{}, limit, page, total int64) responses.Pagination {
	totalPages := calculateTotalPages(total, limit)
	if limit == 0 {
		limit = 50
	}

	//totalPages := int64(math.Ceil(float64(total) / float64(limit)))
	var pagination responses.Pagination
	pagination.Items = data
	pagination.CurrentPage = int64(page)
	pagination.PreviousPage = GetPreviousPage(page)
	pagination.NextPage = GetNextPage(page, totalPages)
	pagination.SizePerPage = limit
	pagination.TotalItems = total
	pagination.TotalPages = totalPages
	return pagination
}

func calculateTotalPages(total, limit int64) int64 {

	if limit <= 0 {
		return 1
	}

	totalPages := int64(math.Ceil(float64(total) / float64(limit)))

	return totalPages
}

func GetPageNumber(pageFromQueryParams string) (int64, error) {
	if pageFromQueryParams == "" {
		return int64(1), nil
	}
	page, err := strconv.Atoi(pageFromQueryParams)
	if err != nil {
		return int64(1), err
	}

	if page <= 0 {
		return int64(1), nil
	}

	return int64(page), nil
}

func GetNextPage(page, totalPages int64) int64 {
	var nextPage int64
	if page >= totalPages {
		nextPage = totalPages
	} else {
		nextPage = int64(page) + 1
	}
	return nextPage
}

func GetPreviousPage(page int64) int64 {
	if page == 0 {
		return 1
	}
	if page-1 == 0 {
		return 1
	}
	return page - 1
}

func GetlimitOffsetPage(limitParam, pageParam string, total int64) (int64, int64, int64) {
	limit := DefaultLimit
	if limitParam != "" {
		limitReq, err := strconv.Atoi(limitParam)
		if err != nil {
			limitReq = 0
		}
		limit = int64(limitReq)
	}

	page, err := GetPageNumber(pageParam)
	if err != nil {
		if err != nil {
			page = 0
		}
	}

	offset := int64((page - 1)) * limit
	totalPages := int64(math.Ceil(float64(total) / float64(limit)))
	if totalPages > 0 && page > totalPages {
		page = totalPages
		offset = int64((page - 1)) * limit
	}

	if totalPages > 0 && page > totalPages {
		page = totalPages
		offset = int64((page - 1)) * limit
	}
	return limit, offset, page
}
