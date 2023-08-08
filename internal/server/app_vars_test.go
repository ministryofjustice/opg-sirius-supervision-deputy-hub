package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type mockAppVarsClient struct {
	lastCtx sirius.Context
	err     error
	user    sirius.UserDetails
	deputy  sirius.DeputyDetails
}

func (m *mockAppVarsClient) GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.lastCtx = ctx

	return m.user, m.err
}

func (m *mockAppVarsClient) GetDeputyDetails(ctx sirius.Context, defaultPaTeam int, defaultProTeam int, deputyId int) (sirius.DeputyDetails, error) {
	m.lastCtx = ctx

	return m.deputy, m.err
}

var mockUserDetails = sirius.UserDetails{
	ID: 1,
}

var mockDeputyDetails = sirius.DeputyDetails{
	ID: 2,
}

func TestNewAppVars(t *testing.T) {
	client := &mockAppVarsClient{user: mockUserDetails, deputy: mockDeputyDetails}
	r, _ := http.NewRequest("GET", "/path", nil)

	envVars := EnvironmentVars{DefaultPaTeam: 3, DefaultProTeam: 4}
	vars, err := NewAppVars(client, r, envVars)

	assert.Nil(t, err)
	assert.Equal(t, AppVars{
		Path:            "/path",
		XSRFToken:       "",
		UserDetails:     mockUserDetails,
		DeputyDetails:   mockDeputyDetails,
		Error:           "",
		Errors:          nil,
		EnvironmentVars: envVars,
	}, *vars)
}
