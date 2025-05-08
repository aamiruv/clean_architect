package sqlutil_test

import (
	"fmt"

	"github.com/amirzayi/clean_architect/pkg/paginate"
	"github.com/amirzayi/clean_architect/pkg/sqlutil"
)

func ExampleBuildPaginationQuery() {
	query := sqlutil.BuildPaginationQuery("user", &paginate.Pagination{
		Page:    3,
		PerPage: 15,
		Fields:  "name,id,phone,role,status",
		Sort: map[string]string{
			"id":   paginate.SortOrderDescending,
			"name": paginate.SortOrderAscending,
		},
		Filters: []paginate.Filter{
			{Key: "name", Value: "amir,admin,test", Condition: paginate.FilterIn},
			{Key: "status", Value: "1,2", Condition: paginate.FilterIn},
		},
	})
	fmt.Print(query)

	// Unordered output:
	// SELECT name,id,phone,role,status FROM user
	// WHERE name IN("amir","admin","test") AND status IN(1,2)
	// ORDER BY id desc, name asc
	// LIMIT 15 offset 30
}
