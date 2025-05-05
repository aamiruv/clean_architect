package paginate_test

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/amirzayi/clean_architect/pkg/paginate"
)

func ExampleParseFromHttpRequest() {
	r := &http.Request{URL: &url.URL{
		RawQuery: "first_name=amir&first_name=like&last_name=mirzaei&age=30&age=>=",
	}}

	pagination := paginate.ParseFromHttpRequest(r)

	for query, filter := range pagination.Filters {
		fmt.Printf("%s must %s %s.\n", query, filter.Condition, filter.Value)
	}

	// Unordered output:
	// age must >= 30.
	// first_name must like amir.
	// last_name must = mirzaei.
}
