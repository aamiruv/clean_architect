package paginate_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/amirzayi/clean_architect/pkg/paginate"
	"github.com/stretchr/testify/require"
)

func TestParseFromHttpRequest(t *testing.T) {
	baseUrl := "/somewhere"

	params := url.Values{}
	params.Add("name", "smith")
	params.Add("name", "like")
	params.Add("age", "36")

	url := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	r, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	require.NoError(t, err)

	pagination := paginate.ParseFromHttpRequest(r)

	name, ok := pagination.Filters["name"]
	require.True(t, ok)
	require.NotEmpty(t, name)
	require.Equal(t, "smith", name.Value)
	require.Equal(t, "like", name.Condition)

	age, ok := pagination.Filters["age"]
	require.True(t, ok)
	require.NotEmpty(t, age)
	require.Equal(t, "36", age.Value)
	require.Equal(t, "=", age.Condition)
}
