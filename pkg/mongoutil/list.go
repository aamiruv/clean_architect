package mongoutil

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/amirzayi/clean_architect/pkg/paginate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PaginatedList[T any](ctx context.Context, col *mongo.Collection,
	pagination *paginate.Pagination, queryableFields map[string]string) ([]T, error) {

	options := options.Find().
		SetLimit(int64(pagination.PerPage)).
		SetSkip(int64((pagination.Page - 1) * pagination.PerPage)).
		SetSort(sortAggregate(pagination.Sort)).
		SetProjection(projectionAggregate(pagination.Fields))

	filterAggregate := filterAggregate(pagination.Filters, queryableFields)
	cursor, err := col.Find(ctx, filterAggregate, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var data []T
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}

	count, err := col.CountDocuments(ctx, filterAggregate)
	if err != nil {
		return nil, err
	}
	pagination.SetTotalItems(count)

	return data, nil
}

func sortAggregate(paginationSort []paginate.Sort) bson.D {
	sortAggregate := bson.D{}
	for _, sort := range paginationSort {
		msort := -1
		if sort.Arrange == paginate.SortOrderAscending {
			msort = 1
		}
		sortAggregate = append(sortAggregate, bson.E{Key: sort.Field, Value: msort})
	}
	return sortAggregate
}

func projectionAggregate(fields []string) bson.D {
	projection := bson.D{}
	for _, field := range fields {
		projection = append(projection, bson.E{Key: field, Value: 1})
	}
	return projection
}

func filterAggregate(filters []paginate.Filter, queryableFields map[string]string) bson.D {
	if len(filters) == 0 {
		return bson.D{}
	}

	filterAggregate := bson.D{}
	match := bson.D{}

	for _, filter := range filters {
		field, ok := queryableFields[filter.Key]
		if !ok {
			continue
		}

		switch filter.Condition {
		case paginate.FilterLike:
			match = bson.D{{Key: field, Value: primitive.Regex{Pattern: filter.Value, Options: "i"}}} // "i" for case insensitive

		case paginate.FilterIn:
			values := strings.Split(filter.Value, ",")
			if len(values) < 2 {
				continue
			}

			arr := bson.A{}
			for _, v := range values {
				arr = append(arr, sanitize(v))
			}
			match = bson.D{{Key: field, Value: bson.D{{Key: "$in", Value: arr}}}}

		case paginate.FilterBetween:
			values := strings.Split(filter.Value, ",")
			if len(values) < 2 {
				continue
			}
			match = bson.D{{Key: field, Value: bson.D{{Key: "$gte", Value: sanitize(values[0])}, {Key: "$lte", Value: sanitize(values[1])}}}}

		default:
			match = bson.D{{Key: field, Value: bson.D{{Key: conditionToNosql(filter.Condition), Value: sanitize(filter.Value)}}}}
		}

		filterAggregate = append(filterAggregate, match...)
	}
	return filterAggregate
}

func conditionToNosql(condition string) string {
	switch condition {
	case paginate.FilterEqual:
		return "$eq"
	case paginate.FilterNotEqual:
		return "$neq"
	case paginate.FilterGreater:
		return "$gt"
	case paginate.FilterGreaterEqual:
		return "$gte"
	case paginate.FilterLess:
		return "$lt"
	case paginate.FilterLessEqual:
		return "$lte"

	default:
		return ""
	}
}
func sanitize(v string) any {
	if digit, err := strconv.ParseFloat(v, 64); err == nil {
		return digit
	}
	t, err := time.Parse(time.RFC3339, v)
	if err == nil {
		return primitive.NewDateTimeFromTime(t)
	}
	return v
}
