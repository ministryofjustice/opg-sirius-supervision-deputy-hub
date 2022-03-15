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

		hasSuccess, successMessage := createSuccessAndSuccessMessageForVars(r.URL.String(), deputyDetails.ExecutiveCaseManager.EcmName)

		vars := deputyHubVars{
			Path:           r.URL.Path,
			XSRFToken:      ctx.XSRFToken,
			DeputyDetails:  deputyDetails,
			Success:        hasSuccess,
			SuccessMessage: successMessage,
		}

		vars.ErrorMessage = checkForDefaultEcmId(deputyDetails.ExecutiveCaseManager.EcmId, defaultPATeam)

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

func createSuccessAndSuccessMessageForVars(url string, ecmName string) (bool, string) {
	splitStringByQuestion := strings.Split(url, "?")
	if len(splitStringByQuestion) > 1 {
		splitString := strings.Split(splitStringByQuestion[1], "=")

		switch splitString[1] {
		case "deputyDetails":
			return true, "Deputy details updated"
		case "ecm":
			return true, "Ecm changed to " + ecmName
		case "importantInformation":
			return true, "Important information updated"
		case "newFirm":
			return true, "Firm added"
		case "teamDetails":
			return true, "Team details updated"
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
