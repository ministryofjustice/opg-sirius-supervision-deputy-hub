package urlbuilder

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestUrlBuilder_buildUrl(t *testing.T) {
	tests := []struct {
		originalPath string
		page         int
		perPage      int
		filters      []Filter
		sort         Sort
		want         string
	}{
		{
			originalPath: "slug",
			page:         11,
			perPage:      25,
			filters:      nil,
			want:         "slug?limit=25&page=11",
		},
		{
			originalPath: "",
			page:         0,
			perPage:      0,
			filters:      []Filter{},
			want:         "?limit=0&page=0",
		},
		{
			originalPath: "slug",
			page:         11,
			perPage:      25,
			filters: []Filter{
				{
					Name:           "test",
					SelectedValues: nil,
				},
			},
			want: "slug?limit=25&page=11",
		},
		{
			originalPath: "slug",
			page:         11,
			perPage:      25,
			filters: []Filter{
				{
					Name:           "test",
					SelectedValues: []string{""},
				},
			},
			want: "slug?limit=25&page=11",
		},
		{
			originalPath: "slug",
			page:         11,
			perPage:      25,
			filters: []Filter{
				{
					Name:           "test",
					SelectedValues: []string{"val"},
				},
			},
			want: "slug?limit=25&page=11&test=val",
		},
		{
			originalPath: "slug",
			page:         11,
			perPage:      25,
			filters: []Filter{
				{
					Name:           "test",
					SelectedValues: []string{"val1", "val2"},
				},
			},
			want: "slug?limit=25&page=11&test=val1&test=val2",
		},
		{
			originalPath: "slug",
			page:         11,
			perPage:      25,
			filters: []Filter{
				{
					Name:           "test",
					SelectedValues: []string{"val1", "val2"},
				},
				{
					Name:           "test2",
					SelectedValues: []string{"val3"},
				},
			},
			want: "slug?limit=25&page=11&test=val1&test=val2&test2=val3",
		},
		{
			originalPath: "",
			page:         2,
			perPage:      15,
			filters:      []Filter{},
			sort:         Sort{OrderBy: "name"},
			want:         "?limit=15&page=2&order-by=name&sort=asc",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			builder := UrlBuilder{OriginalPath: test.originalPath}
			url := builder.buildUrl(test.page, test.perPage, test.filters, test.sort)
			assert.Equal(t, test.want, url)
		})
	}
}

func TestUrlBuilder_GetPaginationUrl(t *testing.T) {
	tests := []struct {
		urlBuilder UrlBuilder
		page       int
		perPage    []int
		want       string
	}{
		{
			urlBuilder: UrlBuilder{OriginalPath: "page", SelectedPerPage: 25},
			page:       2,
			perPage:    []int{25},
			want:       "page?limit=25&page=2",
		},
		{
			urlBuilder: UrlBuilder{SelectedPerPage: 25},
			page:       1,
			perPage:    []int{50},
			want:       "?limit=50&page=1",
		},
		{
			urlBuilder: UrlBuilder{SelectedPerPage: 100},
			page:       2,
			perPage:    nil,
			want:       "?limit=100&page=2",
		},
		{
			urlBuilder: UrlBuilder{SelectedPerPage: 100, SelectedSort: Sort{OrderBy: "name"}},
			page:       2,
			perPage:    nil,
			want:       "?limit=100&page=2&order-by=name&sort=asc",
		},
		{
			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
				{
					Name:                  "retained1",
					SelectedValues:        []string{"val1", "val2"},
					ClearBetweenTeamViews: false,
				},
				{
					Name:                  "retained2",
					SelectedValues:        []string{"val3"},
					ClearBetweenTeamViews: true,
				},
			}},
			page:    2,
			perPage: nil,
			want:    "?limit=0&page=2&retained1=val1&retained1=val2&retained2=val3",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			var result string
			if test.perPage == nil {
				result = test.urlBuilder.GetPaginationUrl(test.page)
			} else {
				result = test.urlBuilder.GetPaginationUrl(test.page, test.perPage[0])
			}
			assert.Equal(t, test.want, result)
		})
	}
}

func TestUrlBuilder_GetClearFiltersUrl(t *testing.T) {
	tests := []struct {
		urlBuilder UrlBuilder
		want       string
	}{
		{
			urlBuilder: UrlBuilder{OriginalPath: "page", SelectedPerPage: 50, SelectedSort: Sort{OrderBy: "name"}},
			want:       "page?limit=50&page=1&order-by=name&sort=asc",
		},
		{
			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
				{
					Name:                  "removed1",
					SelectedValues:        []string{"val1"},
					ClearBetweenTeamViews: true,
				},
				{
					Name:                  "removed2",
					SelectedValues:        []string{"val2"},
					ClearBetweenTeamViews: false,
				},
			}},
			want: "?limit=0&page=1",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, test.urlBuilder.GetClearFiltersUrl())
		})
	}
}

func TestUrlBuilder_GetRemoveFilterUrl(t *testing.T) {
	tests := []struct {
		urlBuilder    UrlBuilder
		name          string
		value         interface{}
		want          string
		expectedError error
	}{
		{
			urlBuilder:    UrlBuilder{OriginalPath: "page", SelectedPerPage: 50, SelectedSort: Sort{OrderBy: "name"}},
			name:          "non-existent-filter",
			value:         "non-existent-value",
			want:          "page?limit=50&page=1&order-by=name&sort=asc",
			expectedError: nil,
		},
		{
			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
				{
					Name:           "filter1",
					SelectedValues: []string{"val1"},
				},
			}},
			name:          "filter1",
			value:         "non-existent-value",
			want:          "?limit=0&page=1&filter1=val1",
			expectedError: nil,
		},
		{
			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
				{
					Name:           "filter1",
					SelectedValues: []string{"val1"},
				},
			}},
			name:          "filter1",
			value:         "val1",
			want:          "?limit=0&page=1",
			expectedError: nil,
		},
		{
			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
				{
					Name:           "filter1",
					SelectedValues: []string{"val1", "val2"},
				},
			}},
			name:          "filter1",
			value:         "val1",
			want:          "?limit=0&page=1&filter1=val2",
			expectedError: nil,
		},
		{
			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
				{
					Name:           "filter1",
					SelectedValues: []string{"val1", "val2"},
				},
				{
					Name:           "filter2",
					SelectedValues: []string{"val3"},
				},
			}},
			name:          "filter2",
			value:         "val3",
			want:          "?limit=0&page=1&filter1=val1&filter1=val2",
			expectedError: nil,
		},
		{
			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
				{
					Name:           "filter1",
					SelectedValues: []string{"val1", "val2"},
				},
				{
					Name:           "filter2",
					SelectedValues: []string{"23"},
				},
			}},
			name:          "filter2",
			value:         23,
			want:          "?limit=0&page=1&filter1=val1&filter1=val2",
			expectedError: nil,
		},
		{
			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
				{
					Name:           "filter1",
					SelectedValues: []string{"val1", "val2"},
				},
				{
					Name:           "filter2",
					SelectedValues: []string{"23", "45", "66"},
				},
			}},
			name:          "filter2",
			value:         []int{23, 45, 66},
			want:          "",
			expectedError: fmt.Errorf("type []int not accepted"),
		},
		{
			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
				{
					Name:           "filter1",
					SelectedValues: []string{"val1", "val2"},
				},
			}},
			name:          "filter2",
			value:         []string{"val1", "val2", "val3"},
			want:          "",
			expectedError: fmt.Errorf("type []string not accepted"),
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			returnedValue, returnedError := test.urlBuilder.GetRemoveFilterUrl(test.name, test.value)
			assert.Equal(t, test.want, returnedValue)
			assert.Equal(t, test.expectedError, returnedError)
		})
	}
}

func TestUrlBuilder_GetSortUrl(t *testing.T) {
	tests := []struct {
		urlBuilder UrlBuilder
		orderBy    string
		want       string
	}{
		{
			urlBuilder: UrlBuilder{},
			orderBy:    "test",
			want:       "?limit=0&page=1&order-by=test&sort=asc",
		},
		{
			urlBuilder: UrlBuilder{SelectedSort: Sort{OrderBy: "test2", Descending: true}},
			orderBy:    "test",
			want:       "?limit=0&page=1&order-by=test&sort=asc",
		},
		{
			urlBuilder: UrlBuilder{SelectedSort: Sort{OrderBy: "test"}},
			orderBy:    "test",
			want:       "?limit=0&page=1&order-by=test&sort=desc",
		},
		{
			urlBuilder: UrlBuilder{SelectedSort: Sort{OrderBy: "test", Descending: true}},
			orderBy:    "test",
			want:       "?limit=0&page=1&order-by=test&sort=asc",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, test.urlBuilder.GetSortUrl(test.orderBy))
		})
	}
}
