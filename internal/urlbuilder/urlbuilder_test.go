package urlbuilder

//
//import (
//	"fmt"
//	"github.com/stretchr/testify/assert"
//	"strconv"
//	"testing"
//)
//
//func TestUrlBuilder_buildUrl(t *testing.T) {
//	tests := []struct {
//		path    string
//		page    int
//		perPage int
//		filters []Filter
//		want    string
//	}{
//		{
//			path:    "slug",
//			page:    11,
//			perPage: 25,
//			filters: nil,
//			want:    "slug?limit=25&page=11",
//		},
//		{
//			path:    "",
//			page:    0,
//			perPage: 0,
//			filters: []Filter{},
//			want:    "?limit=0&page=0",
//		},
//		{
//			path:    "slug",
//			page:    11,
//			perPage: 25,
//			filters: []Filter{
//				{
//					Name:           "test",
//					SelectedValues: nil,
//				},
//			},
//			want: "slug?limit=25&page=11",
//		},
//		{
//			path:    "slug",
//			page:    11,
//			perPage: 25,
//			filters: []Filter{
//				{
//					Name:           "test",
//					SelectedValues: []string{""},
//				},
//			},
//			want: "slug?limit=25&page=11",
//		},
//		{
//			path:    "slug",
//			page:    11,
//			perPage: 25,
//			filters: []Filter{
//				{
//					Name:           "test",
//					SelectedValues: []string{"val"},
//				},
//			},
//			want: "slug?limit=25&page=11&test=val",
//		},
//		{
//			path:    "slug",
//			page:    11,
//			perPage: 25,
//			filters: []Filter{
//				{
//					Name:           "test",
//					SelectedValues: []string{"val1", "val2"},
//				},
//			},
//			want: "slug?limit=25&page=11&test=val1&test=val2",
//		},
//		{
//			path:    "slug",
//			page:    11,
//			perPage: 25,
//			filters: []Filter{
//				{
//					Name:           "test",
//					SelectedValues: []string{"val1", "val2"},
//				},
//				{
//					Name:           "test2",
//					SelectedValues: []string{"val3"},
//				},
//			},
//			want: "slug?limit=25&page=11&test=val1&test=val2&test2=val3",
//		},
//	}
//	for i, test := range tests {
//		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
//			builder := UrlBuilder{OriginalPath: test.path}
//			url := builder.buildUrl(test.page, test.perPage, test.filters)
//			assert.Equal(t, test.want, url)
//		})
//	}
//}
//
//func TestUrlBuilder_GetClearFiltersUrl(t *testing.T) {
//	tests := []struct {
//		urlBuilder UrlBuilder
//		want       string
//	}{
//		{
//			urlBuilder: UrlBuilder{OriginalPath: "page", SelectedPerPage: 50},
//			want:       "page?limit=50&page=1",
//		},
//		{
//			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
//				{
//					Name:           "removed1",
//					SelectedValues: []string{"val1"},
//				},
//				{
//					Name:           "removed2",
//					SelectedValues: []string{"val2"},
//				},
//			}},
//			want: "?limit=0&page=1",
//		},
//	}
//	for i, test := range tests {
//		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
//			assert.Equal(t, test.want, test.urlBuilder.GetClearFiltersUrl())
//		})
//	}
//}
//
//func TestUrlBuilder_GetRemoveFilterUrl(t *testing.T) {
//	tests := []struct {
//		urlBuilder    UrlBuilder
//		name          string
//		value         interface{}
//		want          string
//		expectedError error
//	}{
//		{
//			urlBuilder:    UrlBuilder{OriginalPath: "page", SelectedPerPage: 50},
//			name:          "non-existent-filter",
//			value:         "non-existent-value",
//			want:          "page?limit=50&page=1",
//			expectedError: nil,
//		},
//		{
//			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
//				{
//					Name:           "filter1",
//					SelectedValues: []string{"val1"},
//				},
//			}},
//			name:          "filter1",
//			value:         "non-existent-value",
//			want:          "?limit=0&page=1&filter1=val1",
//			expectedError: nil,
//		},
//		{
//			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
//				{
//					Name:           "filter1",
//					SelectedValues: []string{"val1"},
//				},
//			}},
//			name:          "filter1",
//			value:         "val1",
//			want:          "?limit=0&page=1",
//			expectedError: nil,
//		},
//		{
//			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
//				{
//					Name:           "filter1",
//					SelectedValues: []string{"val1", "val2"},
//				},
//			}},
//			name:          "filter1",
//			value:         "val1",
//			want:          "?limit=0&page=1&filter1=val2",
//			expectedError: nil,
//		},
//		{
//			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
//				{
//					Name:           "filter1",
//					SelectedValues: []string{"val1", "val2"},
//				},
//				{
//					Name:           "filter2",
//					SelectedValues: []string{"val3"},
//				},
//			}},
//			name:          "filter2",
//			value:         "val3",
//			want:          "?limit=0&page=1&filter1=val1&filter1=val2",
//			expectedError: nil,
//		},
//		{
//			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
//				{
//					Name:           "filter1",
//					SelectedValues: []string{"val1", "val2"},
//				},
//				{
//					Name:           "filter2",
//					SelectedValues: []string{"23"},
//				},
//			}},
//			name:          "filter2",
//			value:         23,
//			want:          "?limit=0&page=1&filter1=val1&filter1=val2",
//			expectedError: nil,
//		},
//		{
//			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
//				{
//					Name:           "filter1",
//					SelectedValues: []string{"val1", "val2"},
//				},
//				{
//					Name:           "filter2",
//					SelectedValues: []string{"23", "45", "66"},
//				},
//			}},
//			name:          "filter2",
//			value:         []int{23, 45, 66},
//			want:          "",
//			expectedError: fmt.Errorf("type []int not accepted"),
//		},
//		{
//			urlBuilder: UrlBuilder{SelectedFilters: []Filter{
//				{
//					Name:           "filter1",
//					SelectedValues: []string{"val1", "val2"},
//				},
//			}},
//			name:          "filter2",
//			value:         []string{"val1", "val2", "val3"},
//			want:          "",
//			expectedError: fmt.Errorf("type []string not accepted"),
//		},
//	}
//	for i, test := range tests {
//		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
//			returnedValue, returnedError := test.urlBuilder.GetRemoveFilterUrl(test.name, test.value)
//			assert.Equal(t, test.want, returnedValue)
//			assert.Equal(t, test.expectedError, returnedError)
//		})
//	}
//}
