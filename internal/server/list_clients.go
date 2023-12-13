package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/urlbuilder"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type DeputyHubClientInformation interface {
	GetDeputyClients(sirius.Context, sirius.ClientListParams) (sirius.ClientList, error)
	GetAccommodationTypes(sirius.Context) ([]model.RefData, error)
	GetSupervisionLevels(sirius.Context) ([]model.RefData, error)
}

type ListClientsVars struct {
	Clients sirius.ClientList
	ListPage
	FilterByOrderStatus
	FilterByAccommodation
	FilterBySupervisionLevel
}

func (lcv ListClientsVars) CreateUrlBuilder() urlbuilder.UrlBuilder {
	return urlbuilder.UrlBuilder{
		OriginalPath:    "clients",
		SelectedPerPage: lcv.PerPage,
		SelectedSort:    lcv.Sort,
		SelectedFilters: []urlbuilder.Filter{
			urlbuilder.CreateFilter("order-status", lcv.SelectedOrderStatuses),
			urlbuilder.CreateFilter("accommodation", lcv.SelectedAccommodationTypes),
			urlbuilder.CreateFilter("supervision-level", lcv.SelectedSupervisionLevels),
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

	for _, k := range lcv.SupervisionLevels {
		if k.IsIn(lcv.SelectedSupervisionLevels) {
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
		group, groupCtx := errgroup.WithContext(ctx.Context)
		urlParams := r.URL.Query()
		page := paginate.GetRequestedPage(urlParams.Get("page"))
		perPageOptions := []int{25, 50, 100}
		perPage := paginate.GetRequestedElementsPerPage(urlParams.Get("limit"), perPageOptions)
		search, _ := strconv.Atoi(r.FormValue("page"))

		orderStatuses := []model.OrderStatus{
			{
				Handle:      "ACTIVE",
				Incomplete:  "Active",
				Category:    "Active",
				Complete:    "Active",
				StatusCount: 0,
			},
			{
				Handle:      "CLOSED",
				Incomplete:  "Closed",
				Category:    "Closed",
				Complete:    "Closed",
				StatusCount: 0,
			},
		}

		selectedOrderStatuses, selectedAccommodationTypes, selectedSupervisionLevels := getFiltersFromParams(urlParams)

		sort := urlbuilder.CreateSortFromURL(urlParams, []string{"surname", "reportDue"}) //this will have all the thing you can sort on ie risk

		params := sirius.ClientListParams{
			DeputyId:           app.DeputyId(),
			Limit:              perPage,
			Search:             search,
			DeputyType:         app.DeputyType(),
			Sort:               fmt.Sprintf("%s:%s", sort.OrderBy, sort.GetDirection()),
			OrderStatuses:      selectedOrderStatuses,
			AccommodationTypes: selectedAccommodationTypes,
			SupervisionLevels:  selectedSupervisionLevels,
		}

		var vars ListClientsVars

		group.Go(func() error {
			clients, err := client.GetDeputyClients(ctx.With(groupCtx), params)
			if err != nil {
				return err
			}
			vars.Clients = clients
			return nil
		})

		group.Go(func() error {
			accommodationTypes, err := client.GetAccommodationTypes(ctx.With(groupCtx))
			if err != nil {
				return err
			}
			vars.AccommodationTypes = accommodationTypes
			return nil
		})

		group.Go(func() error {
			supervisionLevels, err := client.GetSupervisionLevels(ctx.With(groupCtx))
			if err != nil {
				return err
			}
			vars.SupervisionLevels = supervisionLevels
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		app.PageName = "Clients"
		vars.PerPage = perPage

		selectedOrderStatuses = vars.ValidateSelectedOrderStatuses(selectedOrderStatuses, orderStatuses)
		vars.OrderStatuses = orderStatuses
		vars.SelectedOrderStatuses = selectedOrderStatuses
		vars.SelectedAccommodationTypes = selectedAccommodationTypes
		vars.SelectedSupervisionLevels = selectedSupervisionLevels

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

		vars.AppVars = app
		vars.Sort = sort
		vars.UrlBuilder = vars.CreateUrlBuilder()

		if page > vars.Clients.Pages.PageTotal && vars.Clients.Pages.PageTotal > 0 {
			return Redirect(vars.UrlBuilder.GetPaginationUrl(vars.Clients.Pages.PageTotal, perPage))
		}

		vars.Pagination = paginate.Pagination{
			CurrentPage:     vars.Clients.Pages.PageCurrent,
			TotalPages:      vars.Clients.Pages.PageTotal,
			TotalElements:   vars.Clients.TotalClients,
			ElementsPerPage: vars.PerPage,
			ElementName:     "clients",
			PerPageOptions:  perPageOptions,
			UrlBuilder:      vars.UrlBuilder,
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

func getFiltersFromParams(params url.Values) ([]string, []string, []string) {
	var selectedOrderStatuses, selectedAccommodationTypes, selectedSupervisionLevels []string

	if params.Has("order-status") {
		selectedOrderStatuses = params["order-status"]
	}
	if params.Has("accommodation") {
		selectedAccommodationTypes = params["accommodation"]
	}
	if params.Has("supervision-level") {
		selectedSupervisionLevels = params["supervision-level"]
	}
	return selectedOrderStatuses, selectedAccommodationTypes, selectedSupervisionLevels
}
