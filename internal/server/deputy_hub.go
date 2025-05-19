package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"html"
	"html/template"
	"net/http"
	"net/url"
)

type DeputyHubInformation interface {
	GetDeputyClients(sirius.Context, sirius.ClientListParams) (sirius.ClientList, error)
}

type deputyHubVars struct {
	SuccessMessage    template.HTML
	ActiveClientCount int
	AppVars
}

func renderTemplateForDeputyHub(client DeputyHubInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		var selectedOrderStatuses []string
		selectedOrderStatuses = append(selectedOrderStatuses, "ACTIVE")

		params := sirius.ClientListParams{
			DeputyId:      app.DeputyId(),
			Search:        1,
			DeputyType:    app.DeputyType(),
			OrderStatuses: selectedOrderStatuses,
		}

		clientList, err := client.GetDeputyClients(ctx, params)
		if err != nil {
			return err
		}

		app.PageName = "Deputy details"

		vars := deputyHubVars{
			AppVars:           app,
			ActiveClientCount: clientList.Metadata.TotalActiveClients,
			SuccessMessage:    template.HTML(getSuccessFromUrl(r.URL, app.DeputyDetails)),
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

func getSuccessFromUrl(url *url.URL, deputyDetails sirius.DeputyDetails) string {

	switch url.Query().Get("success") {
	case "deputyDetails":
		return "Deputy details updated"
	case "ecm":
		return "<abbr title='Executive Case Manager'>ECM</abbr> changed to " + html.EscapeString(deputyDetails.ExecutiveCaseManager.EcmName)
	case "importantInformation":
		return "Important information updated"
	case "newFirm":
		return "Firm added"
	case "firm":
		return "Firm changed to " + html.EscapeString(deputyDetails.Firm.FirmName)
	case "teamDetails":
		return "Team details updated"
	default:
		return ""
	}
}
