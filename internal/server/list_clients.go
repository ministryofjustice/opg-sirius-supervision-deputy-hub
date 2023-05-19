package server

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type DeputyHubClientInformation interface {
	GetDeputyClients(sirius.Context, int, int, int, string, string, string) (sirius.ClientList, sirius.AriaSorting, error)
	GetPageDetails(sirius.Context, sirius.ClientList, int, int) sirius.PageDetails
}

type ListClientsVars struct {
	Path                 string
	XSRFToken            string
	AriaSorting          sirius.AriaSorting
	DeputyClientsDetails sirius.DeputyClientDetails
	ClientList           sirius.ClientList
	PageDetails          sirius.PageDetails
	DeputyDetails        sirius.DeputyDetails
	Error                string
	ActiveClientCount    int
}

func renderTemplateForClientTab(client DeputyHubClientInformation, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		urlParams := r.URL.Query()

		deputyId, _ := strconv.Atoi(routeVars["id"])
		search, _ := strconv.Atoi(r.FormValue("page"))
		displayClientLimit, _ := strconv.Atoi(r.FormValue("limit"))
		if displayClientLimit == 0 {
			displayClientLimit = 25
		}

		columnBeingSorted, sortOrder := parseUrl(urlParams)

		clientList, ariaSorting, err := client.GetDeputyClients(ctx, deputyId, displayClientLimit, search, deputyDetails.DeputyType.Handle, columnBeingSorted, sortOrder)
		if err != nil {
			return err
		}

		pageDetails := client.GetPageDetails(ctx, clientList, search, displayClientLimit)

		vars := ListClientsVars{
			Path:                 r.URL.Path,
			XSRFToken:            ctx.XSRFToken,
			DeputyClientsDetails: clientList.Clients,
			ClientList:           clientList,
			PageDetails:          pageDetails,
			DeputyDetails:        deputyDetails,
			AriaSorting:          ariaSorting,
			ActiveClientCount:    clientList.Metadata.TotalActiveClients,
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
