package sqlutil

import (
	"context"
	"fmt"
	"strings"

	"github.com/amirzayi/clean_architect/pkg/paginate"
	"github.com/jmoiron/sqlx"
)

func PaginatedList[T any](ctx context.Context,
	db *sqlx.DB, table string,
	pagination *paginate.Pagination, queryableFields map[string]string) ([]T, error) {
	var data []T
	query, args := BuildPaginationQuery(table, pagination, queryableFields)

	if err := db.SelectContext(ctx, &data, query, args...); err != nil {
		return nil, err
	}
	var count int64
	whereQuery, whereArgs := whereQuery(pagination.Filters, queryableFields)
	countQuery := fmt.Sprintf("SELECT count(1) FROM %s %s", table, whereQuery)

	if err := db.GetContext(ctx, &count, countQuery, whereArgs...); err != nil {
		return nil, err
	}

	pagination.SetTotalItems(count)
	return data, nil
}

func BuildPaginationQuery(table string,
	pagination *paginate.Pagination, queryableFields map[string]string) (string, []any) {
	var query strings.Builder

	var args []any

	query.WriteString(selectQuery(table, pagination.Fields))
	query.WriteString("\n")

	whereQuery, whereArgs := whereQuery(pagination.Filters, queryableFields)
	args = append(args, whereArgs...)
	query.WriteString(whereQuery)
	query.WriteString("\n")

	query.WriteString(orderByQuery(pagination.Sort))
	query.WriteString("\n")

	limit, limitArgs := limitQuery(pagination.Page, pagination.PerPage)
	args = append(args, limitArgs...)
	query.WriteString(limit)

	return query.String(), args
}

func selectQuery(table string, fields []string) string {
	selectFields := "*"
	if len(fields) > 0 {
		selectFields = strings.Join(fields, ",")
	}
	return fmt.Sprintf("SELECT %s FROM %s", selectFields, table)
}

func whereQuery(filters []paginate.Filter, queryableFields map[string]string) (string, []any) {
	if len(filters) == 0 {
		return "", nil
	}

	var (
		query                strings.Builder
		where                string
		hasAlreadyWhereQuery bool
		args                 []any
	)

	query.WriteString("WHERE ")

	for _, filter := range filters {
		field, ok := queryableFields[filter.Key]
		if !ok {
			continue
		}

		switch filter.Condition {
		case paginate.FilterBetween:
			values := strings.Split(filter.Value, ",")
			if len(values) < 2 {
				continue
			}
			where = fmt.Sprintf("%s BETWEEN ? AND ?", field)
			args = append(args, values[0], values[1])

		case paginate.FilterIn:
			values := strings.Split(filter.Value, ",")
			if len(values) < 2 {
				continue
			}

			for _, v := range values {
				args = append(args, v)
			}
			where = fmt.Sprintf("%s IN(?%s)", field, strings.Repeat(",?", len(values)-1))

		default:
			if strings.Contains(filter.Value, ",") {
				continue
			}
			args = append(args, filter.Value)
			where = fmt.Sprintf("%s %s ?", field, conditionToSql(filter.Condition))
		}

		hasAlreadyWhereQuery = true
		query.WriteString(where)
		query.WriteString(" AND ")
	}

	if !hasAlreadyWhereQuery {
		return "", nil
	}

	// remove last " AND " at end of query
	whereQuery := strings.TrimRight(query.String(), " AND ")
	return whereQuery, args
}

func orderByQuery(sorts []paginate.Sort) string {
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

func limitQuery(page, perPage int) (string, []any) {
	return fmt.Sprintf("LIMIT ? offset ?"), []any{perPage, (page - 1) * perPage}
}

func conditionToSql(condition string) string {
	switch condition {
	case paginate.FilterEqual:
		return "="
	case paginate.FilterNotEqual:
		return "<>"
	case paginate.FilterGreater:
		return ">"
	case paginate.FilterGreaterEqual:
		return ">="
	case paginate.FilterLess:
		return "<"
	case paginate.FilterLessEqual:
		return "<="

	default:
		return ""
	}
}
