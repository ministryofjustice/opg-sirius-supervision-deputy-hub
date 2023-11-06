package server

import (
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
	GetDeputyClients(sirius.Context, sirius.ClientListParams) (sirius.ClientList, sirius.AriaSorting, error)
	//GetAccommodationTypes(sirius.Context, string) ([]model.RefData, error)
}

type ListClientsVars struct {
	AriaSorting          sirius.AriaSorting
	DeputyClientsDetails sirius.DeputyClientDetails
	Clients              sirius.ClientList
	Pagination           paginate.Pagination
	PerPage              int
	ActiveClientCount    int
	ColumnBeingSorted    string
	SortOrder            string
	AppliedFilters       []string
	OrderStatusOptions   []model.RefData
	urlbuilder.UrlBuilder
	//ListPage
	//FilterByAccommodationType
	FilterByOrderStatus
	SelectedOrderStatuses []string
	OrderStatuses         []OrderStatus
	AppVars
}

func (vars ListClientsVars) CreateUrlBuilder() urlbuilder.UrlBuilder {
	return urlbuilder.UrlBuilder{
		OriginalPath: "clients",
		SortBy:       vars.SortBy,
		SelectedFilters: []urlbuilder.Filter{
			urlbuilder.CreateFilter("order-status", vars.SelectedOrderStatuses),
			//urlbuilder.CreateFilter("accommodation-type", vars.SelectedAccommodationTypes),
		},
	}
}

func (vars ListClientsVars) GetAppliedFilters() []string {
	var appliedFilters []string
	for _, u := range vars.OrderStatuses {
		if u.IsSelected(vars.SelectedOrderStatuses) {
			appliedFilters = append(appliedFilters, u.Handle)
		}
	}
	//for _, u := range vars.AccommodationTypes {
	//	if u.IsIn(vars.SelectedOrderStatuses) {
	//		appliedFilters = append(appliedFilters, u.Label)
	//	}
	//}
	return appliedFilters
}

func (vars ListClientsVars) ValidateSelectedOrderStatuses(selectedOrderStatuses []string, orderStatuses []OrderStatus) []string {
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

type OrderStatus struct {
	Handle      string `json:"handle"`
	Incomplete  string `json:"incomplete"`
	Category    string `json:"category"`
	Complete    string `json:"complete"`
	StatusCount int
}

func (os OrderStatus) IsSelected(selectedOrderStatuses []string) bool {
	for _, selectedOrderStatus := range selectedOrderStatuses {
		if os.Handle == selectedOrderStatus {
			return true
		}
	}
	return false
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

		columnBeingSorted, sortOrder := parseUrl(urlParams)

		orderStatuses := []OrderStatus{
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

		var selectedOrderStatuses []string
		if urlParams.Has("order-status") {
			selectedOrderStatuses = urlParams["order-status"]
		}

		//var selectedAccommodationTypes []string
		//if urlParams.Has("accommodation-type") {
		//	selectedAccommodationTypes = urlParams["accommodation-type"]
		//}

		params := sirius.ClientListParams{
			DeputyId:          app.DeputyId(),
			Limit:             perPage,
			Search:            search,
			DeputyType:        app.DeputyType(),
			ColumnBeingSorted: columnBeingSorted,
			SortOrder:         sortOrder,
			OrderStatuses:     selectedOrderStatuses,
			//AccommodationTypes: selectedAccommodationTypes,
		}

		clients, ariaSorting, err := client.GetDeputyClients(ctx, params)

		if err != nil {
			return err
		}

		app.PageName = "Clients"

		vars := ListClientsVars{
			Clients:           clients,
			PerPage:           perPage,
			AriaSorting:       ariaSorting,
			ColumnBeingSorted: columnBeingSorted,
			SortOrder:         sortOrder,
		}

		vars.App = app
		vars.SelectedOrderStatuses = selectedOrderStatuses
		//vars.SelectedAccommodationTypes = selectedAccommodationTypes

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

		vars.UrlBuilder = vars.CreateUrlBuilder()
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

		//vars.AccommodationTypes, err = client.GetAccommodationTypes(ctx, "clientAccommodation")
		//if err != nil {
		//	return err
		//}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

func parseUrl(urlParams url.Values) (string, string) {
	sortParam := urlParams.Get("sort")
	if sortParam != "" {
		sortParamsArray := strings.Split(sortParam, ":")
		columnBeingSorted := sortParamsArray[0]
		sortOrder := sortParamsArray[1]
		return columnBeingSorted, sortOrder
	}
	return "", ""
}
