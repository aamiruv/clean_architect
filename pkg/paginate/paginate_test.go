package paginate_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/amirzayi/clean_architect/pkg/paginate"
	"github.com/stretchr/testify/require"
)

func TestParseFromRequest(t *testing.T) {
	baseUrl := "/somewhere"

	params := url.Values{}
	params.Add("name", "smith")
	params.Add("name", "like")
	params.Add("age", "36")

	url := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	r, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	require.NoError(t, err)

	pagination := paginate.ParseFromRequest(r)

	require.Contains(t, pagination.Filters, paginate.Filter{
		Key: "name", Value: "smith", Condition: "like",
	})

	require.Contains(t, pagination.Filters, paginate.Filter{
		Key: "age", Value: "36", Condition: paginate.FilterEqual,
	})
}
