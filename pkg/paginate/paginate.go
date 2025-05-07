package paginate

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	defaultPerPageSize     = 10
	defaultSortingArrange  = SortOrderDescending
	defaultFilterCondition = FilterEqual

	pageParamName    = "page"
	perPageParamName = "per_page"
	sortParamName    = "sort"
	fieldsParamName  = "fields"

	FilterEqual        = "eq"
	FilterNotEqual     = "neq"
	FilterGreater      = "gt"
	FilterGreaterEqual = "gte"
	FilterLess         = "lt"
	FilterLessEqual    = "lte"
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
	Fields     string   `json:"fields,omitempty"`
	Sort       []Sort   `json:"sort,omitempty"`
	Filters    []Filter `json:"filters,omitempty"`
	TotalItems int64    `json:"total_items"`
}

type ListResponse struct {
	Data       any         `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func ParseFromRequest(r *http.Request) *Pagination {
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

func SQLList[T any](ctx context.Context, db *sqlx.DB, table string, pagination *Pagination, queryableFields ...string) ([]T, error) {
	var data []T
	var query strings.Builder

	query.WriteString(selectQuery(table, pagination.Fields))
	query.WriteString("\n")
	whereQuery := whereQuery(pagination.Filters)
	query.WriteString(whereQuery)
	query.WriteString("\n")
	query.WriteString(orderByQuery(pagination.Sort))
	query.WriteString("\n")
	query.WriteString(limitQuery(pagination.Page, pagination.PerPage))

	if err := db.SelectContext(ctx, &data, query.String()); err != nil {
		return nil, err
	}

	var count int64
	countQuery := fmt.Sprintf("SELECT count(1) FROM %s %s", table, whereQuery)

	if err := db.GetContext(ctx, &count, countQuery); err != nil {
		return nil, err
	}

	pagination.SetTotalItems(count)
	return data, nil
}

func selectQuery(table, fields string) string {
	if fields == "" {
		fields = "*"
	}
	return fmt.Sprintf("SELECT %s FROM %s", fields, table)
}

func whereQuery(filters []Filter, queryableFields ...string) string {
	if len(filters) == 0 {
		return ""
	}

	var (
		query                strings.Builder
		where                string
		hasAlreadyWhereQuery bool
	)

	query.WriteString("WHERE ")

	for _, filter := range filters {
		if len(queryableFields) > 0 {
			if !slices.Contains(queryableFields, filter.Key) {
				continue
			}
		}
		switch filter.Condition {
		case FilterBetween:
			values := strings.Split(filter.Value, ",")
			if len(values) < 2 {
				continue
			}
			where = fmt.Sprintf("%s BETWEEN %s AND %s", filter.Key, values[0], values[1])

		case FilterIn:
			values := strings.Split(filter.Value, ",")
			if len(values) < 2 {
				continue
			}
			args := filter.Value
			_, err := strconv.Atoi(values[0])
			if err != nil {
				args = strings.Join(values, `","`)
				args = `"` + args + `"`
			}
			where = fmt.Sprintf("%s IN(%s)", filter.Key, args)

		default:
			where = fmt.Sprintf("%s %s %q", filter.Key, filter.Condition, filter.Value)
		}

		query.WriteString(where)
		query.WriteString(" AND ")
	}

	if !hasAlreadyWhereQuery {
		return ""
	}

	// remove last " AND " at end of query
	whereQuery := strings.TrimRight(query.String(), " AND ")
	return whereQuery
}

func orderByQuery(sorts []Sort) string {
	if len(sorts) == 0 {
		return ""
	}

	var query strings.Builder

	query.WriteString("ORDER BY")

	for _, sort := range sorts {
		query.WriteString(fmt.Sprintf(" %s %s,", sort.Field, sort.Arrange))
	}

	// remove last "," character at end of query
	orderByQuery := strings.TrimRight(query.String(), ",")
	return orderByQuery
}

func limitQuery(page, perPage int) string {
	return fmt.Sprintf(" LIMIT %d offset %d", perPage, (page-1)*perPage)
}
