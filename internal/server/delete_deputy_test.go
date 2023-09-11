package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDeleteDeputy struct {
	count             int
	lastCtx           sirius.Context
	deleteDeputyError error
}

func (m *mockDeleteDeputy) DeleteDeputy(ctx sirius.Context, deputyId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.deleteDeputyError
}

func TestGetDeleteDeputy(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeleteDeputy{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeleteDeputy(client, template)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}

func TestPostDeleteDeputy(t *testing.T) {
	assert := assert.New(t)
	client := &mockDeleteDeputy{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	handler := renderTemplateForDeleteDeputy(client, template)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)
	assert.NotEmpty(template.lastVars)
}
