package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"io"

	"github.com/stretchr/testify/mock"
)

type mockTemplates struct {
	mock.Mock
	count    int
	lastName string
	lastVars interface{}
	error    error
}

func (m *mockTemplates) ExecuteTemplate(w io.Writer, name string, vars interface{}) error {
	m.count += 1
	m.lastName = name
	m.lastVars = vars

	return nil
}

type mockApiClient struct {
	error              error
	CurrentUserDetails sirius.UserDetails
	DeputyDetails      sirius.DeputyDetails
}

func (m mockApiClient) GetUserDetails(sirius.Context) (sirius.UserDetails, error) {
	return m.CurrentUserDetails, m.error
}

func (m mockApiClient) GetDeputyDetails(sirius.Context, int, int, int) (sirius.DeputyDetails, error) {
	return m.DeputyDetails, m.error
}
