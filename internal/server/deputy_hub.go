package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type DeputyHubInformation interface {
	GetDeputyDetails(sirius.Context, int, int) (sirius.DeputyDetails, error)
}

type deputyHubVars struct {
	Path           string
	XSRFToken      string
	DeputyDetails  sirius.DeputyDetails
	Error          string
	ErrorMessage   string
	Errors         sirius.ValidationErrors
	Success        bool
	SuccessMessage string
}

func renderTemplateForDeputyHub(client DeputyHubInformation, defaultPATeam int, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])
		deputyDetails, err := client.GetDeputyDetails(ctx, defaultPATeam, deputyId)
		if err != nil {
			return err
		}

		vars := deputyHubVars{
			Path:          r.URL.Path,
			XSRFToken:     ctx.XSRFToken,
			DeputyDetails: deputyDetails,
		}

		vars.Success, vars.SuccessMessage = createSuccessAndSuccessMessageForVars(r.URL.String(), deputyDetails.ExecutiveCaseManager.EcmName)
		vars.ErrorMessage = checkForDefaultEcmId(deputyDetails.ExecutiveCaseManager.EcmId, defaultPATeam)

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

func createSuccessAndSuccessMessageForVars(url string, EcmName string) (bool, string) {
	splitStringByQuestion := strings.Split(url, "?")
	if len(splitStringByQuestion) > 1 {
		splitString := strings.Split(splitStringByQuestion[1], "=")

		if splitString[1] == "ecm" {
			return true, "Ecm changed to " + EcmName
		} else if splitString[1] == "teamDetails" {
			return true, "Team details updated"
		} else if splitString[1] == "deputyDetails" {
			return true, "Deputy details updated"
		} else if splitString[1] == "importantInformation" {
			return true, "Important information updated"
		}
	}
	return false, ""
}

func checkForDefaultEcmId(EcmId, defaultPaTeam int) string {
	if EcmId == defaultPaTeam {
		return "An executive case manager has not been assigned. "
	}
	return ""
}
