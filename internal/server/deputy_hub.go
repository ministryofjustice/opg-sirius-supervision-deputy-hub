package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/url"
	"strconv"
)

type DeputyHubInformation interface {
	GetDeputyClients(sirius.Context, int, string, string, string) (sirius.DeputyClientDetails, sirius.AriaSorting, int, error)
	GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error)
}

type deputyHubVars struct {
	Path              string
	XSRFToken         string
	DeputyDetails     sirius.DeputyDetails
	Error             string
	ErrorMessage      string
	Errors            sirius.ValidationErrors
	SuccessMessage    string
	ActiveClientCount int
	IsFinanceManager  bool
}

func renderTemplateForDeputyHub(client DeputyHubInformation, defaultPATeam int, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, deputyDetails sirius.DeputyDetails, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])
		_, _, clientCount, err := client.GetDeputyClients(ctx, deputyId, deputyDetails.DeputyType.Handle, "", "")
		if err != nil {
			return err
		}
		successMessage := getSuccessFromUrl(r.URL, deputyDetails.ExecutiveCaseManager.EcmName, deputyDetails.Firm.FirmName)

		vars := deputyHubVars{
			Path:              r.URL.Path,
			XSRFToken:         ctx.XSRFToken,
			DeputyDetails:     deputyDetails,
			SuccessMessage:    successMessage,
			ActiveClientCount: clientCount,
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

		vars.ErrorMessage = checkForDefaultEcmId(deputyDetails.ExecutiveCaseManager.EcmId, defaultPATeam)

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
		return "Ecm changed to " + ecmName
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

func checkForDefaultEcmId(EcmId, defaultPaTeam int) string {
	if EcmId == defaultPaTeam {
		return "An executive case manager has not been assigned. "
	}
	return ""
}
