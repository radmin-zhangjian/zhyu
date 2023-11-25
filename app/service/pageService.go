package service

import (
	"log"
	"math"
)

type pagesModel struct {
	Page      int
	PageSize  int
	PageCount int64
	PageTotal int
	PageNext  int
}

// PagesService 美化分页
func PagesService(page int, pageSize int, total int64) any {
	pageCount := total
	pageNext := 2
	pageTotalFloat := float64(total) / float64(pageSize)
	log.Printf("total: %#v, pageSize: %#v, pageTotalFloat: %#v", total, pageSize, pageTotalFloat)
	pageTotal := int(math.Ceil(pageTotalFloat))
	if page >= pageTotal {
		pageNext = 0
	} else {
		pageNext = page + 1
	}
	pages := pagesModel{
		Page:      page,
		PageSize:  pageSize,
		PageCount: pageCount,
		PageTotal: pageTotal,
		PageNext:  pageNext,
	}
	return pages
}
