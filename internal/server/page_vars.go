package server

import (
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/urlbuilder"
	"reflect"
)

type ListPage struct {
	App            AppVars
	AppliedFilters []string
	Sort           urlbuilder.Sort
	Error          string
	Pagination     paginate.Pagination
	PerPage        int
	UrlBuilder     urlbuilder.UrlBuilder
}

type FilterByOrderStatus struct {
	ListPage
	OrderStatusOptions    []model.RefData
	SelectedOrderStatuses []string
	OrderStatuses         []model.OrderStatus
}

func (lp ListPage) HasFilterBy(page interface{}, filter string) bool {
	filters := map[string]interface{}{
		"order-status": FilterByOrderStatus{},
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

func (vars ListClientsVars) ValidateSelectedOrderStatuses(selectedOrderStatuses []string, orderStatuses []model.OrderStatus) []string {
	var validSelectedOrderStatuses []string
	for _, selectedOrderStatus := range selectedOrderStatuses {
		for _, orderStatus := range orderStatuses {
			if selectedOrderStatus == orderStatus.Handle {
				validSelectedOrderStatuses = append(validSelectedOrderStatuses, selectedOrderStatus)
				break
			}
		}
	}
	return validSelectedOrderStatuses
}
