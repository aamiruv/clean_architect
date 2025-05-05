package paginate

import (
	"net/http"
	"slices"
	"strconv"
)

const (
	defaultPerPageSize = 10

	pageParamName    = "page"
	perPageParamName = "per_page"
	sortParamName    = "sort"
	fieldsParamName  = "fields"
)

type Filter struct {
	Value     string `json:"value"`
	Condition string `json:"condition"`
}

type Pagination struct {
	Page       int               `json:"page,omitempty"`
	PerPage    int               `json:"per_page,,omitempty"`
	Sort       string            `json:"sort,omitempty"`
	Fields     string            `json:"fields,omitempty"`
	Filters    map[string]Filter `json:"filters,omitempty"`
	TotalItems int64             `json:"total_items"`
}

type List struct {
	Data       any         `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func ParseFromHttpRequest(r *http.Request) *Pagination {
	page := 1
	if r.URL.Query().Has(pageParamName) {
		pageNo, err := strconv.Atoi(r.URL.Query().Get(pageParamName))
		if err == nil && pageNo > 0 {
			page = pageNo
		}
	}

	perPage := defaultPerPageSize
	if r.URL.Query().Has(perPageParamName) {
		perPageNo, err := strconv.Atoi(r.URL.Query().Get(perPageParamName))
		if err == nil && perPageNo > 0 {
			perPage = perPageNo
		}
	}

	sort := r.URL.Query().Get(sortParamName)

	fields := r.URL.Query().Get(fieldsParamName)

	filters := make(map[string]Filter)

	queries := r.URL.Query()

	for query, values := range queries {
		if slices.Contains([]string{
			pageParamName,
			perPageParamName,
			sortParamName,
			fieldsParamName,
		}, query) {
			continue
		}

		param := Filter{
			Value:     values[0],
			Condition: "=",
		}
		if len(values) > 1 {
			param.Condition = values[1]
		}

		filters[query] = param
	}

	return &Pagination{
		Page:    page,
		PerPage: perPage,
		Fields:  fields,
		Sort:    sort,
		Filters: filters,
	}
}

func (p *Pagination) SetTotalItems(totalItems int64) {
	p.TotalItems = totalItems
}
