package urlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChangeSortButtonDirection(t *testing.T) {
	tests := []struct {
		sortOrder        string
		functionCalling  string
		expectedResponse string
		sortOrderBy      string
	}{
		{sortOrder: "asc", functionCalling: "surname", expectedResponse: "none", sortOrderBy: "reportdue"},
		{sortOrder: "desc", functionCalling: "reportdue", expectedResponse: "none", sortOrderBy: "surname"},
		{sortOrder: "asc", functionCalling: "reportdue", expectedResponse: "descending", sortOrderBy: "reportdue"},
		{sortOrder: "desc", functionCalling: "reportdue", expectedResponse: "ascending", sortOrderBy: "reportdue"},
	}

	for _, tc := range tests {
		s := Sort{OrderBy: tc.sortOrderBy}
		result := s.ChangeSortButtonDirection(tc.functionCalling, tc.sortOrder)
		assert.Equal(t, tc.expectedResponse, result)
	}
}

func TestAriaSorting_GetHTMLSortDirection(t *testing.T) {
	tests := []struct {
		orderingBy       string
		expectedResponse string
		sortOrderBy      string
		descending       bool
	}{
		{orderingBy: "reportdue", expectedResponse: "asc", sortOrderBy: "reportdue", descending: false},
		{orderingBy: "reportdue", expectedResponse: "desc", sortOrderBy: "reportdue", descending: true},
		{orderingBy: "surname", expectedResponse: "asc", sortOrderBy: "reportdue", descending: true},
	}

	for _, tc := range tests {
		s := Sort{OrderBy: tc.sortOrderBy, Descending: tc.descending}
		result := s.GetHTMLSortDirection(tc.orderingBy)
		assert.Equal(t, tc.expectedResponse, result)
	}
}

func TestSort_GetAriaSort(t *testing.T) {
	assert.Equal(t, "none", Sort{}.GetAriaSort("test"))
	assert.Equal(t, "none", Sort{Descending: true}.GetAriaSort("test"))
	assert.Equal(t, "none", Sort{OrderBy: "foo", Descending: true}.GetAriaSort("test"))
	assert.Equal(t, "ascending", Sort{OrderBy: "test"}.GetAriaSort("test"))
	assert.Equal(t, "descending", Sort{OrderBy: "test", Descending: true}.GetAriaSort("test"))
}

func TestSort_ToURL(t *testing.T) {
	assert.Equal(t, "", Sort{}.ToURL())
	assert.Equal(t, "order-by=test&sort=asc", Sort{OrderBy: "test", Descending: false}.ToURL())
	assert.Equal(t, "order-by=test&sort=desc", Sort{OrderBy: "test", Descending: true}.ToURL())
}
