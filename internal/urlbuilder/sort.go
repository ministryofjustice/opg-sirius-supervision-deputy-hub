package urlbuilder

import (
	"fmt"
	"net/url"
)

type Sort struct {
	OrderBy    string
	SortOrder  string
	Descending bool
}

func CreateSortFromURL(values url.Values, validOptions []string) Sort {
	if len(validOptions) == 0 {
		return Sort{}
	}
	sort := Sort{
		OrderBy:    values.Get("order-by"),
		Descending: values.Get("sort") == "desc",
	}
	for _, validSort := range validOptions {
		if sort.OrderBy == validSort {
			return sort
		}
	}
	return Sort{OrderBy: validOptions[0]}
}

func (s Sort) GetDirection() string {
	if s.Descending {
		return "desc"
	}
	return "asc"
}

func (s Sort) ToURL() string {
	if s.OrderBy == "" {
		return ""
	}
	return fmt.Sprintf("order-by=%s&sort=%s", s.OrderBy, s.GetDirection())
}

func (s Sort) GetAriaSort(orderBy string) string {
	if s.OrderBy != orderBy {
		return "none"
	}
	if s.Descending {
		return "descending"
	}
	return "ascending"
}

func (s Sort) GetHTMLSortDirection(orderingBy string) string {
	if orderingBy == s.OrderBy {
		return s.GetDirection()
	} else {
		return "asc"
	}
}

func (s Sort) ChangeSortButtonDirection(orderingBy string, currentDirection string) string {
	if orderingBy == s.OrderBy {
		if currentDirection != "asc" {
			return "ascending"
		} else if currentDirection != "desc" {
			return "descending"
		}
		return "none"
	} else {
		return "none"
	}
}
