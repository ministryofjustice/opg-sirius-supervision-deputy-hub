package server

import (
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/urlbuilder"
	"reflect"
)

type ListPage struct {
	App            AppVars
	Error          string
	AppliedFilters []string
	Pagination     paginate.Pagination
	PerPage        int
	UrlBuilder     urlbuilder.UrlBuilder
}

type FilterByOrderStatus struct {
	ListPage
	OrderStatuses         []model.RefData
	SelectedOrderStatuses []string
}

//type FilterByAccommodationType struct {
//	ListPage
//	AccommodationTypes         []model.RefData
//	SelectedAccommodationTypes []string
//}

func (lp ListPage) HasFilterBy(page interface{}, filter string) bool {
	filters := map[string]interface{}{
		"order-status": FilterByOrderStatus{},
		//"accommodation-type": FilterByAccommodationType{},
	}

	extends := func(parent interface{}, child interface{}) bool {
		p := reflect.TypeOf(parent)
		c := reflect.TypeOf(child)
		for i := 0; i < p.NumField(); i++ {
			if f := p.Field(i); f.Type == c && f.Anonymous {
				return true
			}
		}
		return false
	}

	if f, ok := filters[filter]; ok {
		return extends(page, f)
	}
	return false
}
