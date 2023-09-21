package urlbuilder

import (
	"fmt"
	"strconv"
)

type UrlBuilder struct {
	SelectedFilters []Filter
}

func (ub UrlBuilder) buildUrl(filters []Filter) string {
	url := "clients"
	for _, filter := range filters {
		for _, value := range filter.SelectedValues {
			if value != "" {
				url += "&" + filter.Name + "=" + value
			}
		}
	}
	return url
}

func (ub UrlBuilder) GetClearFiltersUrl() string {
	return ub.buildUrl([]Filter{})
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
				retainedValues = append(retainedValues, v)
			}
		}
		if len(retainedValues) > 0 {
			filter.SelectedValues = retainedValues
			retainedFilters = append(retainedFilters, filter)
		}
	}

	return ub.buildUrl(retainedFilters), nil
}
