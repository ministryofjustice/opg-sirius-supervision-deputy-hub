package server

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type DeputyHubClientInformation interface {
	GetDeputyClients(sirius.Context, int, int, int, string, string, string) (sirius.ClientList, sirius.AriaSorting, error)
	GetPageDetails(sirius.Context, sirius.ClientList, int, int) sirius.PageDetails
}

type ListClientsVars struct {
	AriaSorting          sirius.AriaSorting
	DeputyClientsDetails sirius.DeputyClientDetails
	ClientList           sirius.ClientList
	PageDetails          sirius.PageDetails
	ActiveClientCount    int
	AppVars
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

		clientList, ariaSorting, err := client.GetDeputyClients(ctx, app.DeputyId(), displayClientLimit, search, app.DeputyType(), columnBeingSorted, sortOrder)
		if err != nil {
			return err
		}

		pageDetails := client.GetPageDetails(ctx, clientList, search, displayClientLimit)

		vars := ListClientsVars{
			DeputyClientsDetails: clientList.Clients,
			ClientList:           clientList,
			PageDetails:          pageDetails,
			AriaSorting:          ariaSorting,
			ActiveClientCount:    clientList.Metadata.TotalActiveClients,
			AppVars:              app,
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
