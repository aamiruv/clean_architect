package sqlutil

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/amirzayi/clean_architect/pkg/paginate"
	"github.com/jmoiron/sqlx"
)

func PaginatedList[T any](ctx context.Context,
	db *sqlx.DB, table string,
	pagination *paginate.Pagination, queryableFields ...string) ([]T, error) {
	var data []T
	query := BuildPaginationQuery(table, pagination)
	if err := db.SelectContext(ctx, &data, query); err != nil {
		return nil, err
	}

	var count int64
	whereQuery := whereQuery(pagination.Filters, queryableFields...)
	countQuery := fmt.Sprintf("SELECT count(1) FROM %s %s", table, whereQuery)

	if err := db.GetContext(ctx, &count, countQuery); err != nil {
		return nil, err
	}

	pagination.SetTotalItems(count)
	return data, nil
}

func BuildPaginationQuery(table string,
	pagination *paginate.Pagination, queryableFields ...string) string {
	var query strings.Builder

	query.WriteString(selectQuery(table, pagination.Fields))
	query.WriteString("\n")
	whereQuery := whereQuery(pagination.Filters, queryableFields...)
	query.WriteString(whereQuery)
	query.WriteString("\n")
	query.WriteString(orderByQuery(pagination.Sort))
	query.WriteString("\n")
	query.WriteString(limitQuery(pagination.Page, pagination.PerPage))

	return query.String()
}

func selectQuery(table, fields string) string {
	if fields == "" {
		fields = "*"
	}
	return fmt.Sprintf("SELECT %s FROM %s", fields, table)
}

func whereQuery(filters []paginate.Filter, queryableFields ...string) string {
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
		case paginate.FilterBetween:
			values := strings.Split(filter.Value, ",")
			if len(values) < 2 {
				continue
			}
			where = fmt.Sprintf("%s BETWEEN %s AND %s", filter.Key, values[0], values[1])

		case paginate.FilterIn:
			values := strings.Split(filter.Value, ",")
			if len(values) < 2 {
				continue
			}
			args := filter.Value
			_, err := strconv.ParseFloat(values[0], 64)
			if err != nil {
				args = strings.Join(values, `","`)
				args = `"` + args + `"`
			}
			where = fmt.Sprintf("%s IN(%s)", filter.Key, args)

		default:
			where = fmt.Sprintf("%s %s %q", filter.Key, filter.Condition, filter.Value)
		}

		hasAlreadyWhereQuery = true
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

func orderByQuery(sorts map[string]string) string {
	if len(sorts) == 0 {
		return ""
	}

	var query strings.Builder

	query.WriteString("ORDER BY")

	for field, arrange := range sorts {
		query.WriteString(fmt.Sprintf(" %s %s,", field, arrange))
	}

	// remove last "," character at end of query
	orderByQuery := strings.TrimRight(query.String(), ",")
	return orderByQuery
}

func limitQuery(page, perPage int) string {
	return fmt.Sprintf("LIMIT %d offset %d", perPage, (page-1)*perPage)
}
