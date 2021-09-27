package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type UserDetailsClient interface {
	UserDetails(sirius.Context) (sirius.UserDetails, error)
}

func userDetails(client UserDetailsClient, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		userDetails, err := client.UserDetails(ctx)
		if err != nil {
			return err
		}

		vars := userDetails.ID

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
