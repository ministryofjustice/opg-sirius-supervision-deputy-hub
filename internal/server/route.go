package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"net/http"
)

type PageData struct {
	Data           AppVars
	SuccessMessage string
	MyDetails      sirius.UserDetails
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
func (r route) execute(w http.ResponseWriter, req *http.Request, data any, appVars AppVars) error {
	if r.isHxRequest(req) {
		return r.tmpl.ExecuteTemplate(w, r.partial, data)
	} else {
		storedData := appVars
		fmt.Println("do i have user details")
		fmt.Println(storedData.UserDetails)
		fmt.Println("do i have deputy details")
		fmt.Println(storedData.DeputyDetails)
		//data.SuccessMessage = r.getSuccess(req)
		fmt.Println("in route go")
		fmt.Println(data)
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
