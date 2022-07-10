package query

import "fmt"

type Pager interface {
	Page() int
	PageSize() int
	SearchCount() bool
	String() string
}

type PageRequest struct {
	page        int
	pageSize    int
	searchCount bool
}

func NewPageRequest(page int, pageSize int, searchCount bool) *PageRequest {
	return &PageRequest{page: page, pageSize: pageSize, searchCount: searchCount}
}

func (p PageRequest) String() string {
	return fmt.Sprintf("Limit %d, %d", (p.page-1)*p.pageSize, p.pageSize)
}

func (p *PageRequest) Page() int {
	return p.page
}

func (p *PageRequest) SetPage(page int) {
	p.page = page
}

func (p *PageRequest) PageSize() int {
	return p.pageSize
}

func (p *PageRequest) SetPageSize(pageSize int) {
	p.pageSize = pageSize
}

func (p *PageRequest) SearchCount() bool {
	return p.searchCount
}

func (p *PageRequest) SetSearchCount(searchCount bool) {
	p.searchCount = searchCount
}
