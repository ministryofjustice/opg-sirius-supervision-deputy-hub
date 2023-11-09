package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/urlbuilder"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type DeputyHubClientInformation interface {
	GetDeputyClients(sirius.Context, sirius.ClientListParams) (sirius.ClientList, error)
	GetAccommodationTypes(sirius.Context, string) ([]model.RefData, error)
}

type ListClientsVars struct {
	Clients sirius.ClientList
	ListPage
	FilterByOrderStatus
	FilterByAccommodation
}

func (lcv ListClientsVars) CreateUrlBuilder() urlbuilder.UrlBuilder {
	return urlbuilder.UrlBuilder{
		OriginalPath: "clients",
		SelectedSort: lcv.Sort,
		SelectedFilters: []urlbuilder.Filter{
			urlbuilder.CreateFilter("order-status", lcv.SelectedOrderStatuses),
			urlbuilder.CreateFilter("accommodation-types", lcv.SelectedAccommodationTypes),
		},
	}
}

func (lcv ListClientsVars) GetAppliedFilters() []string {
	var appliedFilters []string
	for _, u := range lcv.OrderStatuses {
		if u.IsSelected(lcv.SelectedOrderStatuses) {
			appliedFilters = append(appliedFilters, u.Incomplete)
		}
	}

	for _, k := range lcv.AccommodationTypes {
		if k.IsIn(lcv.SelectedAccommodationTypes) {
			appliedFilters = append(appliedFilters, k.Label)
		}
	}
	return appliedFilters
}

func renderTemplateForClientTab(client DeputyHubClientInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		urlParams := r.URL.Query()
		page := paginate.GetRequestedPage(urlParams.Get("page"))
		perPageOptions := []int{25, 50, 100}
		perPage := paginate.GetRequestedElementsPerPage(urlParams.Get("limit"), perPageOptions)
		search, _ := strconv.Atoi(r.FormValue("page"))

		var columnBeingSorted, sortOrder, boolSortOrder = parseUrl(urlParams)

		orderStatuses := []model.OrderStatus{
			{
				"ACTIVE",
				"Active",
				"Active",
				"Active",
				0,
			},
			{
				"CLOSED",
				"Closed",
				"Closed",
				"Closed",
				0,
			},
		}

		selectedOrderStatuses, selectedAccommodationTypes := getFiltersFromParams(urlParams)

		params := sirius.ClientListParams{
			DeputyId:           app.DeputyId(),
			Limit:              perPage,
			Search:             search,
			DeputyType:         app.DeputyType(),
			ColumnBeingSorted:  columnBeingSorted,
			SortOrder:          sortOrder,
			OrderStatuses:      selectedOrderStatuses,
			AccommodationTypes: selectedAccommodationTypes,
		}

		clients, err := client.GetDeputyClients(ctx, params)

		if err != nil {
			return err
		}

		app.PageName = "Clients"

		var vars ListClientsVars

		vars.Clients = clients
		vars.PerPage = perPage

		vars.Sort = urlbuilder.Sort{
			OrderBy:    columnBeingSorted,
			Descending: boolSortOrder,
			SortOrder:  sortOrder,
		}
		vars.AppVars = app
		vars.UrlBuilder = vars.CreateUrlBuilder()

		if page > clients.Pages.PageTotal && clients.Pages.PageTotal > 0 {
			return Redirect(vars.UrlBuilder.GetPaginationUrl(clients.Pages.PageTotal, perPage))
		}

		vars.Pagination = paginate.Pagination{
			CurrentPage:     clients.Pages.PageCurrent,
			TotalPages:      clients.Pages.PageTotal,
			TotalElements:   clients.TotalClients,
			ElementsPerPage: vars.PerPage,
			ElementName:     "clients",
			PerPageOptions:  perPageOptions,
			UrlBuilder:      vars.UrlBuilder,
		}

		selectedOrderStatuses = vars.ValidateSelectedOrderStatuses(selectedOrderStatuses, orderStatuses)
		vars.OrderStatuses = orderStatuses
		vars.SelectedOrderStatuses = selectedOrderStatuses
		vars.SelectedAccommodationTypes = selectedAccommodationTypes

		vars.AppliedFilters = vars.GetAppliedFilters()

		vars.OrderStatusOptions = []model.RefData{
			{
				Handle: "ACTIVE",
				Label:  "Active",
			},
			{
				Handle: "CLOSED",
				Label:  "Closed",
			},
		}

		vars.AccommodationTypes, err = client.GetAccommodationTypes(ctx, "clientAccommodation")
		if err != nil {
			return err
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

func parseUrl(urlParams url.Values) (string, string, bool) {
	sortParam := urlParams.Get("sort")
	boolSortOrder := false
	if sortParam != "" {
		sortParamsArray := strings.Split(sortParam, ":")
		columnBeingSorted := sortParamsArray[0]
		sortOrder := sortParamsArray[1]
		if sortOrder == "asc" {
			boolSortOrder = true
		}
		return columnBeingSorted, sortOrder, boolSortOrder
	}
	return "", "", boolSortOrder
}

func getFiltersFromParams(params url.Values) ([]string, []string) {
	var selectedOrderStatuses []string
	var selectedAccommodationTypes []string

	if params.Has("order-status") {
		selectedOrderStatuses = params["order-status"]
	}
	if params.Has("accommodation") {
		selectedAccommodationTypes = params["accommodation"]
	}
	fmt.Println("selectedAccommodationTypes")
	fmt.Println(selectedAccommodationTypes)
	fmt.Println("selectedOrderStatuses")
	fmt.Println(selectedOrderStatuses)
	return selectedOrderStatuses, selectedAccommodationTypes
}
