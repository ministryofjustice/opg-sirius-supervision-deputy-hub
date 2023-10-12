package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/urlbuilder"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type DeputyHubClientInformation interface {
	GetDeputyClients(sirius.Context, sirius.ClientListParams) (sirius.ClientList, sirius.AriaSorting, error)
	GetPageDetails(sirius.Context, sirius.ClientList, int, int) sirius.PageDetails
}

type ListClientsVars struct {
	AriaSorting           sirius.AriaSorting
	DeputyClientsDetails  sirius.DeputyClientDetails
	ClientList            sirius.ClientList
	PageDetails           sirius.PageDetails
	ActiveClientCount     int
	ColumnBeingSorted     string
	SortOrder             string
	DisplayClientLimit    int
	SelectedOrderStatuses []string
	OrderStatuses         []OrderStatus
	AppliedFilters        []string
	OrderStatusOptions    []model.RefData
	FilterByOrderStatus
	urlbuilder.UrlBuilder
	AppVars
}

type FilterByOrderStatus struct {
	OrderStatusOptions    []model.RefData
	SelectedOrderStatuses []string
}

func (vars ListClientsVars) CreateUrlBuilder() urlbuilder.UrlBuilder {
	return urlbuilder.UrlBuilder{
		OriginalPath:    "clients",
		SortBy:          vars.SortBy,
		SelectedPerPage: vars.DisplayClientLimit,
		SelectedFilters: []urlbuilder.Filter{
			urlbuilder.CreateFilter("order-status", vars.SelectedOrderStatuses),
		},
	}
}

func (vars ListClientsVars) HasFilterBy(page interface{}, filter string) bool {
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

func (vars ListClientsVars) GetAppliedFilters() []string {
	var appliedFilters []string
	for _, u := range vars.OrderStatuses {
		if u.IsSelected(vars.SelectedOrderStatuses) {
			appliedFilters = append(appliedFilters, u.Incomplete)
		}
	}
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

		search, _ := strconv.Atoi(r.FormValue("page"))
		displayClientLimit, _ := strconv.Atoi(r.FormValue("limit"))
		if displayClientLimit == 0 {
			displayClientLimit = 25
		}

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

		params := sirius.ClientListParams{
			DeputyId:           app.DeputyId(),
			DisplayClientLimit: displayClientLimit,
			Search:             search,
			DeputyType:         app.DeputyType(),
			ColumnBeingSorted:  columnBeingSorted,
			SortOrder:          sortOrder,
			OrderStatuses:      selectedOrderStatuses,
		}

		clientList, ariaSorting, err := client.GetDeputyClients(ctx, params)

		if err != nil {
			return err
		}

		pageDetails := client.GetPageDetails(ctx, clientList, search, displayClientLimit)

		app.PageName = "Clients"

		vars := ListClientsVars{
			DeputyClientsDetails: clientList.Clients,
			ClientList:           clientList,
			PageDetails:          pageDetails,
			AriaSorting:          ariaSorting,
			ColumnBeingSorted:    columnBeingSorted,
			SortOrder:            sortOrder,
			DisplayClientLimit:   displayClientLimit,
			AppVars:              app,
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
