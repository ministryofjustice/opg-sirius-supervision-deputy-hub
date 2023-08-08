package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"html/template"
	"net/http"
	"net/url"
)

type DeputyHubInformation interface {
	GetDeputyClients(sirius.Context, int, int, int, string, string, string) (sirius.ClientList, sirius.AriaSorting, error)
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

		vars := deputyHubVars{
			AppVars: app,
		}

		clientList, _, err := client.GetDeputyClients(ctx, app.DeputyId(), 25, 1, app.DeputyDetails.DeputyType.Handle, "", "")
		if err != nil {
			return err
		}
		vars.ActiveClientCount = clientList.Metadata.TotalActiveClients
		vars.SuccessMessage = template.HTML(getSuccessFromUrl(r.URL, app.DeputyDetails.ExecutiveCaseManager.EcmName, app.DeputyDetails.Firm.FirmName))

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

func getSuccessFromUrl(url *url.URL, ecmName string, firmName string) string {
	switch url.Query().Get("success") {
	case "deputyDetails":
		return "Deputy details updated"
	case "ecm":
		return "<abbr title='Executive Case Manager'>ECM</abbr> changed to " + ecmName
	case "importantInformation":
		return "Important information updated"
	case "newFirm":
		return "Firm added"
	case "firm":
		return "Firm changed to " + firmName
	case "teamDetails":
		return "Team details updated"
	default:
		return ""
	}
}
