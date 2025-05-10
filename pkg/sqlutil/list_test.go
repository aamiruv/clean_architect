package sqlutil_test

import (
	"testing"

	"github.com/amirzayi/clean_architect/pkg/paginate"
	"github.com/amirzayi/clean_architect/pkg/sqlutil"
	"github.com/stretchr/testify/require"
)

func TestBuildPaginationQuery(t *testing.T) {
	query, args := sqlutil.BuildPaginationQuery("user", &paginate.Pagination{
		Page:    3,
		PerPage: 15,
		Fields:  []string{"name", "id", "phone", "role", "status"},
		Sort: []paginate.Sort{
			{Field: "id", Arrange: paginate.SortOrderDescending},
			{Field: "name", Arrange: paginate.SortOrderAscending},
		},
		Filters: []paginate.Filter{
			{Key: "name", Value: "amir,admin,test", Condition: paginate.FilterIn},
			{Key: "status", Value: "1,2", Condition: paginate.FilterIn},
		},
	}, map[string]string{
		"id":         "id",
		"name":       "name",
		"phone":      "phone",
		"email":      "email",
		"status":     "status",
		"role":       "role",
		"created_at": "created_at",
	})

	require.NotEmpty(t, query)
	require.NotEmpty(t, args)
	require.Contains(t, query, "SELECT name,id,phone,role,status FROM user")
	require.Contains(t, query, "WHERE name IN(?,?,?) AND status IN(?,?)")
	require.Contains(t, query, "ORDER BY id desc, name asc")
	require.Contains(t, query, "LIMIT ? offset ?")
}
