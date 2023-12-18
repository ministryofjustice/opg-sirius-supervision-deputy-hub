package urlbuilder

import (
	"fmt"
	"strconv"
	"strings"
)

type UrlBuilder struct {
	OriginalPath    string
	SelectedPerPage int
	SelectedFilters []Filter
	SelectedSort    Sort
}

func (ub UrlBuilder) buildUrl(page int, perPage int, filters []Filter, sort Sort) string {
	url := fmt.Sprintf("%s?limit=%d&page=%d", ub.OriginalPath, perPage, page)
	for _, filter := range filters {
		for _, value := range filter.SelectedValues {
			if value != "" {
				url += "&" + filter.Name + "=" + value
			}
		}
	}
	if sort.ToURL() != "" {
		url += "&" + sort.ToURL()
	}
	return url
}

func (ub UrlBuilder) GetPaginationUrl(page int, perPage ...int) string {
	selectedPerPage := ub.SelectedPerPage
	if len(perPage) > 0 {
		selectedPerPage = perPage[0]
	}
	return ub.buildUrl(page, selectedPerPage, ub.SelectedFilters, ub.SelectedSort)
}

func (ub UrlBuilder) GetSortUrl(orderBy string) string {
	sort := Sort{OrderBy: orderBy}
	if orderBy == ub.SelectedSort.OrderBy {
		sort.Descending = !ub.SelectedSort.Descending
	}
	return ub.buildUrl(1, ub.SelectedPerPage, ub.SelectedFilters, sort)
}

func (ub UrlBuilder) GetClearFiltersUrl() string {
	return ub.buildUrl(1, ub.SelectedPerPage, []Filter{}, ub.SelectedSort)
}

func (ub UrlBuilder) GetRemoveFilterUrl(name string, value interface{}) (string, error) {
	var retainedFilters []Filter
	var retainedValues []string
	var stringValue string

	switch val := value.(type) {
	case string:
		stringValue = val
	case int:
		stringValue = strconv.Itoa(val)
	default:
		return "", fmt.Errorf("type %T not accepted", val)
	}

	for _, filter := range ub.SelectedFilters {
		retainedValues = []string{}
		for _, v := range filter.SelectedValues {
			if name != filter.Name || stringValue != v {
				formatWhiteSpace := strings.Replace(v, " ", "%20", -1)
				retainedValues = append(retainedValues, formatWhiteSpace)
			}
		}

		if len(retainedValues) > 0 {
			filter.SelectedValues = retainedValues
			retainedFilters = append(retainedFilters, filter)
		}
	}

	return ub.buildUrl(1, ub.SelectedPerPage, retainedFilters, ub.SelectedSort), nil
}
