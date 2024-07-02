package server

import (
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
)

//type AppVars struct {
//	Path          string
//	XSRFToken     string
//	UserDetails   sirius.UserDetails
//	DeputyDetails sirius.DeputyDetails
//	PageName      string
//	Error         string
//	Errors        sirius.ValidationErrors
//	EnvironmentVars
//}

type PageData struct {
	Data           AppVars
	SuccessMessage string
}

type route struct {
	client  ApiClient
	tmpl    Template
	partial string
}

func (r route) Client() ApiClient {
	return r.client
}

// execute is an abstraction of the Template execute functions in order to conditionally render either a full template or
// a block, in response to a header added by HTMX. If the header is not present, the function will also fetch all
// additional data needed by the page for a full page load.
func (r route) execute(w http.ResponseWriter, req *http.Request, data any, envVars EnvironmentVars) error {
	if r.isHxRequest(req) {
		return r.tmpl.ExecuteTemplate(w, r.partial, data)
	} else {
		ctx := getContext(req)
		group, groupCtx := errgroup.WithContext(ctx.Context)
		deputyId, _ := strconv.Atoi(mux.Vars(req)["id"])

		pageInfo := PageData{
			Data: AppVars{
				Path:            req.URL.Path,
				XSRFToken:       ctx.XSRFToken,
				EnvironmentVars: envVars,
			},
			SuccessMessage: "",
		}

		group.Go(func() error {
			user, err := r.client.GetUserDetails(ctx.With(groupCtx))
			if err != nil {
				return err
			}
			pageInfo.Data.UserDetails = user
			return nil
		})
		group.Go(func() error {
			deputy, err := r.client.GetDeputyDetails(ctx.With(groupCtx), pageInfo.Data.DefaultPaTeam, pageInfo.Data.DefaultProTeam, deputyId)
			if err != nil {
				return err
			}
			pageInfo.Data.DeputyDetails = deputy
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		pageInfo.SuccessMessage = r.getSuccess(req)
		return r.tmpl.Execute(w, data)
	}
}

func (r route) getSuccess(req *http.Request) string {
	switch req.URL.Query().Get("success") {
	case "invoice-adjustment[CREDIT WRITE OFF]":
		return "Write-off successfully created"
	case "invoice-adjustment[CREDIT MEMO]":
		return "Manual credit successfully created"
	}
	return ""
}

func (r route) isHxRequest(req *http.Request) bool {
	return req.Header.Get("HX-Request") == "true"
}
