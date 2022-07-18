package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type DeputyHubClientInformation interface {
	GetDeputyClients(sirius.Context, int, int, int, string, string, string) (sirius.ClientList, sirius.AriaSorting, int, error)
	GetPageDetails(sirius.Context, sirius.ClientList, int, int) sirius.PageDetails
}

type listClientsVars struct {
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
	return func(perm sirius.PermissionSet, deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		var displayClientLimit int
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])
		//search, _ := strconv.Atoi(r.FormValue("page"))
		search := 1
		bothDisplayClientLimits := r.Form["clientsPerPage"]
		currentClientDisplay, _ := strconv.Atoi(r.FormValue("currentClientDisplay"))

		if len(bothDisplayClientLimits) != 0 {
			topDisplayClientLimit, _ := strconv.Atoi(bothDisplayClientLimits[0])
			bottomDisplayClientLimit, _ := strconv.Atoi(bothDisplayClientLimits[1])
			if topDisplayClientLimit != currentClientDisplay {
				displayClientLimit = topDisplayClientLimit
			} else if bottomDisplayClientLimit != currentClientDisplay {
				displayClientLimit = bottomDisplayClientLimit
			} else {
				displayClientLimit = currentClientDisplay
			}
		} else {
			displayClientLimit = 25
		}

		columnBeingSorted, sortOrder := parseUrl(r.URL.String())

		clientList, ariaSorting, activeClientCount, err := client.GetDeputyClients(ctx, deputyId, displayClientLimit, search, deputyDetails.DeputyType.Handle, columnBeingSorted, sortOrder)
		if err != nil {
			return err
		}

		pageDetails := client.GetPageDetails(ctx, clientList, search, displayClientLimit)

		vars := listClientsVars{
			Path:                 r.URL.Path,
			XSRFToken:            ctx.XSRFToken,
			DeputyClientsDetails: clientList.Clients,
			ClientList:           clientList,
			PageDetails:          pageDetails,
			DeputyDetails:        deputyDetails,
			AriaSorting:          ariaSorting,
			ActiveClientCount:    activeClientCount,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

func parseUrl(url string) (string, string) {
	urlQuery := strings.Split(url, "?")
	if len(urlQuery) >= 2 {
		sortParams := urlQuery[1]
		sortParamsArray := strings.Split(sortParams, ":")
		columnBeingSorted := sortParamsArray[0]
		sortOrder := sortParamsArray[1]
		return columnBeingSorted, sortOrder
	}
	return "", ""
}
