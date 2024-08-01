package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"html/template"
	"net/http"
	"net/url"
)

type DeputyVars struct {
	SuccessMessage    template.HTML
	ActiveClientCount int
	AppVars
}

type DeputyHandler struct {
	router
}

func (h *DeputyHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return StatusError(http.StatusMethodNotAllowed)
	}
	ctx := getContext(r)
	v.PageName = "Deputy details"
	var selectedOrderStatuses []string
	selectedOrderStatuses = append(selectedOrderStatuses, "ACTIVE")

	params := sirius.ClientListParams{
		DeputyId:      v.DeputyId(),
		Search:        1,
		DeputyType:    v.DeputyType(),
		OrderStatuses: selectedOrderStatuses,
	}

	clientList, err := h.Client().GetDeputyClients(ctx, params)
	if err != nil {
		return err
	}

	v.PageName = "Deputy details"

	vars := DeputyVars{
		AppVars:           v,
		ActiveClientCount: clientList.Metadata.TotalActiveClients,
		SuccessMessage:    template.HTML(getSuccessFromUrl(r.URL, v.DeputyDetails)),
	}
	return h.execute(w, r, vars)
}

func getSuccessFromUrl(url *url.URL, deputyDetails sirius.DeputyDetails) string {

	switch url.Query().Get("success") {
	case "deputyDetails":
		return "Deputy details updated"
	case "ecm":
		return "<abbr title='Executive Case Manager'>ECM</abbr> changed to " + deputyDetails.ExecutiveCaseManager.EcmName
	case "importantInformation":
		return "Important information updated"
	case "newFirm":
		return "Firm added"
	case "firm":
		return "Firm changed to " + deputyDetails.Firm.FirmName
	case "teamDetails":
		return "Team details updated"
	default:
		return ""
	}
}
