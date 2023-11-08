package urlbuilder

import (
	"fmt"
)

type Sort struct {
	OrderBy    string
	SortOrder  string
	Descending bool
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
