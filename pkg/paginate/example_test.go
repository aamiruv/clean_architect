package paginate_test

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/amirzayi/clean_architect/pkg/paginate"
)

func ExampleParseFromRequest() {
	page := "page=3&per_page=15"
	sort := "sort=age&sort=asc&sort=id"
	fields := "fields=first_name,last_name"
	filters := "first_name=amir&first_name=like&last_name=mirzaei&age=30&age=neq&age=26,33&age=between"
	invalidFilters := "id=1; or drop table user"

	r := &http.Request{URL: &url.URL{
		RawQuery: fmt.Sprintf("%s&%s&%s&%s&%s", page, fields, sort, filters, invalidFilters),
	}}

	pagination := paginate.ParseFromRequest(r)

	fmt.Printf("page: %d, per_page: %d.\n", pagination.Page, pagination.PerPage)

	fmt.Println("fields:")
	fmt.Println(pagination.Fields)

	fmt.Println("filters:")
	for _, filter := range pagination.Filters {
		fmt.Printf("%s %s %s.\n", filter.Key, filter.Condition, filter.Value)
	}

	fmt.Println("sorting:")
	for field, arrange := range pagination.Sort {
		fmt.Printf("%s %s.\n", field, arrange)
	}

	// Output:
	// page: 3, per_page: 15.
	// fields:
	// first_name,last_name
	// filters:
	// first_name like amir.
	// last_name eq mirzaei.
	// age neq 30.
	// age between 26,33.
	// sorting:
	// age asc.
	// id desc.
}
