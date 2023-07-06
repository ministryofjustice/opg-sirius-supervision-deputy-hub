package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/url"
	"strconv"
	"html/template"
)

type DeputyHubInformation interface {
	GetDeputyClients(sirius.Context, int, int, int, string, string, string) (sirius.ClientList, sirius.AriaSorting, error)
	GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error)
}

type deputyHubVars struct {
	Path              string
	XSRFToken         string
	DeputyDetails     sirius.DeputyDetails
	Error             string
	SuccessMessage    template.HTML
	ActiveClientCount int
	IsFinanceManager  bool
}

func renderTemplateForDeputyHub(client DeputyHubInformation, tmpl Template) Handler {
	return func(deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])
		clientList, _, err := client.GetDeputyClients(ctx, deputyId, 25, 1, deputyDetails.DeputyType.Handle, "", "")
		if err != nil {
			return err
		}
		successMessage := getSuccessFromUrl(r.URL, deputyDetails.ExecutiveCaseManager.EcmName, deputyDetails.Firm.FirmName)

		vars := deputyHubVars{
			Path:              r.URL.Path,
			XSRFToken:         ctx.XSRFToken,
			DeputyDetails:     deputyDetails,
			SuccessMessage:    template.HTML(successMessage),
			ActiveClientCount: clientList.Metadata.TotalActiveClients,
		}

		group, groupCtx := errgroup.WithContext(ctx.Context)
		group.Go(func() error {
			userDetails, err := client.GetUserDetails(ctx.With(groupCtx))
			if err != nil {
				return err
			}

			vars.IsFinanceManager = userDetails.IsFinanceManager()
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

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
