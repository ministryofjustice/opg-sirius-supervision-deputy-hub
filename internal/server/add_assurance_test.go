package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
)

type mockAddAssuranceClient struct {
	count           int
	lastCtx         sirius.Context
	AddAssuranceErr error
}

func (m *mockAddAssuranceClient) AddAssurance(ctx sirius.Context, assuranceType string, requestedDate string, userId, deputyId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.AddAssuranceErr
}

//func TestGetAddAssurance(t *testing.T) {
//	assert := assert.New(t)
//
//	client := &mockAddAssuranceClient{}
//	ro := &mockRoute{client: client}
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest(http.MethodGet, "", nil)
//	r.SetPathValue("deputyId", "75")
//	r.SetPathValue("assuranceId", "1")
//
//	appVars := AppVars{Path: "/path/"}
//	sut := AddAssuranceHandler{ro}
//	err := sut.render(appVars, w, r)
//
//	assert.Nil(t, err)
//	assert.True(t, ro.executed)
//
//	expected := AddAssuranceVars{
//		appVars,
//	}
//	assert.Equal(t, expected, ro.data)

//handler := AddAssurance(client, template)
//err := handler(AppVars{}, w, r)
//
//assert.Nil(err)
//
//resp := w.Result()
//assert.Equal(http.StatusOK, resp.StatusCode)
//
//assert.Equal(1, template.count)
//assert.Equal("page", template.lastName)
//}

//func TestPostAssurance(t *testing.T) {
//	assert := assert.New(t)
//	client := &mockAddAssuranceClient{}
//	template := &mockTemplates{}
//
//	form := url.Values{}
//	form.Add("assurance-type", "ABC")
//	form.Add("requested-date", "2200/10/20")
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/123/assurances", strings.NewReader(form.Encode()))
//	r.PostForm = form
//
//	var returnedError error
//
//	testHandler := mux.NewRouter()
//	testHandler.HandleFunc("/{id}/assurances", func(w http.ResponseWriter, r *http.Request) {
//		returnedError = AddAssurance(client, template)(AppVars{DeputyDetails: testDeputy}, w, r)
//	})
//
//	testHandler.ServeHTTP(w, r)
//	resp := w.Result()
//	assert.Equal(http.StatusOK, resp.StatusCode)
//	assert.Equal(Redirect("/123/assurances?success=addAssurance"), returnedError)
//}
//
//func TestAddAssuranceHandlesValidationErrorsGeneratedWithinFile(t *testing.T) {
//	assert := assert.New(t)
//	client := &mockAddAssuranceClient{}
//
//	form := url.Values{}
//	form.Add("assurance-type", "")
//	form.Add("requested-date", "")
//
//	template := &mockTemplates{}
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/123/assurances", strings.NewReader(form.Encode()))
//
//	returnedError := AddAssurance(client, template)(AppVars{}, w, r)
//
//	expectedErrors := sirius.ValidationErrors{
//		"assurance-type": {
//			"": "Select an assurance type",
//		},
//		"requested-date": {
//			"": "Enter a requested date",
//		},
//	}
//
//	assert.Equal(AddAssuranceVars{
//		AppVars{
//			Errors:   expectedErrors,
//			PageName: "Add assurance visit",
//		},
//	}, template.lastVars)
//
//	assert.Nil(returnedError)
//}
//
//func TestAddAssuranceHandlesValidationErrorsReturnedFromSiriusCall(t *testing.T) {
//	assert := assert.New(t)
//	client := &mockAddAssuranceClient{}
//
//	validationErrors := sirius.ValidationErrors{
//		"assurance-type": {
//			"": "Select an assurance type",
//		},
//		"requested-date": {
//			"": "Enter a requested date",
//		},
//	}
//
//	client.AddAssuranceErr = sirius.ValidationError{
//		Errors: validationErrors,
//	}
//
//	template := &mockTemplates{}
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/123/assurances", strings.NewReader(""))
//
//	returnedError := AddAssurance(client, template)(AppVars{}, w, r)
//
//	assert.Equal(AddAssuranceVars{
//		AppVars{
//			Errors:   validationErrors,
//			PageName: "Add assurance visit",
//		},
//	}, template.lastVars)
//
//	assert.Nil(returnedError)
//
//}
