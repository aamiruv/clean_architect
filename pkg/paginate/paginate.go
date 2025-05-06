package paginate

import (
	"net/http"
	"slices"
	"strconv"
	"strings"
)

const (
	defaultPerPageSize     = 10
	defaultSortingArrange  = SortOrderDescending
	defaultFilterCondition = FilterEqual

	pageParamName    = "page"
	perPageParamName = "per_page"
	sortParamName    = "sort"
	fieldsParamName  = "fields"

	FilterEqual        = "="
	FilterNotEqual     = "!="
	FilterGreater      = ">"
	FilterGreaterEqual = ">="
	FilterLess         = "<"
	FilterLessEqual    = "<="
	FilterIn           = "in"
	FilterBetween      = "between"
	FilterLike         = "like"

	SortOrderAscending  = "asc"
	SortOrderDescending = "desc"
)

type Filter struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Condition string `json:"condition"`
}

type Sort struct {
	Field   string `json:"field"`
	Arrange string `json:"arrange"`
}

type Pagination struct {
	Page       int      `json:"page,omitempty"`
	PerPage    int      `json:"per_page,,omitempty"`
	Sort       []Sort   `json:"sort,omitempty"`
	Fields     string   `json:"fields,omitempty"`
	Filters    []Filter `json:"filters,omitempty"`
	TotalItems int64    `json:"total_items"`
}

type List struct {
	Data       any         `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func ParseFromHttpRequest(r *http.Request) *Pagination {
	queries := r.URL.Query()

	page := 1
	if queries.Has(pageParamName) {
		pageNo, err := strconv.Atoi(queries.Get(pageParamName))
		if err == nil && pageNo > 0 {
			page = pageNo
		}
	}

	perPage := defaultPerPageSize
	if queries.Has(perPageParamName) {
		perPageNo, err := strconv.Atoi(queries.Get(perPageParamName))
		if err == nil && perPageNo > 0 {
			perPage = perPageNo
		}
	}

	fields := queries.Get(fieldsParamName)

	sort := []Sort{}

	filters := []Filter{}

	for query, values := range queries {
		// prevent sql injection
		if strings.Contains(query, " ") || strings.Contains(query, ";") {
			continue
		}

		// prevent conflicts with pagination
		if slices.Contains([]string{
			pageParamName,
			perPageParamName,
			fieldsParamName,
		}, query) {
			continue
		}

		if query == sortParamName {
			for _, v := range values {
				if isValidSortArrange(v) {
					if len(sort) > 0 {
						sort[len(sort)-1].Arrange = v
					}
					continue
				}
				sorting := Sort{
					Field:   v,
					Arrange: defaultSortingArrange,
				}
				sort = append(sort, sorting)
			}
			continue
		}

		for _, v := range values {
			if isValidFilterCondition(v) {
				if len(filters) > 0 {
					filters[len(filters)-1].Condition = v
				}
				continue
			}

			filter := Filter{
				Key:       query,
				Value:     v,
				Condition: defaultFilterCondition,
			}
			filters = append(filters, filter)
		}
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

func isValidFilterCondition(condition string) bool {
	return slices.Contains([]string{
		FilterEqual,
		FilterNotEqual,
		FilterGreater,
		FilterGreaterEqual,
		FilterLess,
		FilterLessEqual,
		FilterIn,
		FilterBetween,
		FilterLike,
	}, condition)
}

func isValidSortArrange(arrange string) bool {
	return slices.Contains([]string{SortOrderAscending, SortOrderDescending}, arrange)
}
