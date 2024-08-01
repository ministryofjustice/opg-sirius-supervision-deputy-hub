package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//func TestGetFirm(t *testing.T) {
//	assert := assert.New(t)
//
//	client := &mockAddFirm{}
//	template := &mockTemplates{}
//	ro := &mockRoute{client: mockApiClient}
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("GET", "/path", nil)
//
//	handler := AddFirmHandler(client, template)
//	err := *handler.render(addFirmAppVars, w, r)
//
//	assert.Nil(err)
//
//	resp := w.Result()
//	assert.Equal(http.StatusOK, resp.StatusCode)
//
//	assert.Equal(1, template.count)
//	assert.Equal("page", template.lastName)
//}

func TestPostAddFirm(t *testing.T) {
	form := url.Values{
		"name":           {"new-firm-name"},
		"address-line-1": {"123 fake street"},
		"address-line-2": {"fake avenue"},
		"town":           {"Springfield"},
		"county":         {"Texas"},
		"postcode":       {"SP1 1TF"},
		"telephone":      {"01224 587452"},
		"email":          {"fake@email.com"},
	}

	client := mockApiClient{}
	ro := &mockRoute{client: client}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/add", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.SetPathValue("deputyId", "75")

	appVars := AppVars{
		Path: "/add",
	}

	appVars.EnvironmentVars.Prefix = "prefix"

	sut := AddFirmHandler{ro}

	err := sut.render(appVars, w, r)

	assert.Nil(t, err)
	assert.Equal(t, "deputy/75?success=newFirm", w.Header().Get("HX-Redirect"))
}

//func TestAddFirmValidationErrors(t *testing.T) {
//	assert := assert.New(t)
//	client := &mockFirmInformation{}
//
//	validationErrors := sirius.ValidationErrors{
//		"firmName": {
//			"stringLengthTooLong": "The firm name must be 255 characters or fewer",
//		},
//	}
//
//	client.AddFirmDetailsErr = sirius.ValidationError{
//		Errors: validationErrors,
//	}
//
//	template := &mockTemplates{}
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))
//	returnedError := renderTemplateForAddFirm(client, template)(AppVars{}, w, r)
//
//	assert.Equal(addFirmVars{
//		AppVars: AppVars{
//			Errors:   validationErrors,
//			PageName: "Create new firm",
//		},
//	}, template.lastVars)
//
//	assert.Nil(returnedError)
//}
//
//func TestErrorAddFirmMessageWhenIsEmpty(t *testing.T) {
//	assert := assert.New(t)
//	client := &mockFirmInformation{}
//
//	validationErrors := sirius.ValidationErrors{
//		"firmName": {
//			"isEmpty": "The firm name is required and can't be empty",
//		},
//	}
//
//	client.AddFirmDetailsErr = sirius.ValidationError{
//		Errors: validationErrors,
//	}
//
//	template := &mockTemplates{}
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))
//
//	var returnedError error
//
//	testHandler := mux.NewRouter()
//	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
//		returnedError = renderTemplateForAddFirm(client, template)(addFirmAppVars, w, r)
//	})
//
//	testHandler.ServeHTTP(w, r)
//
//	expectedValidationErrors := sirius.ValidationErrors{
//		"firmName": {
//			"isEmpty": "The firm name is required and can't be empty",
//		},
//	}
//
//	assert.Equal(addFirmVars{
//		AppVars: AppVars{
//			DeputyDetails: testDeputy,
//			Errors:        expectedValidationErrors,
//			PageName:      "Create new firm",
//		},
//	}, template.lastVars)
//
//	assert.Nil(returnedError)
//}
//
//func TestAddFirmHandlesErrorsInOtherClientFiles(t *testing.T) {
//	returnedError := sirius.StatusError{Code: 500}
//	tests := []struct {
//		Client *mockFirmInformation
//	}{
//		{
//			Client: &mockFirmInformation{
//				AddFirmDetailsErr: returnedError,
//			},
//		},
//		{
//			Client: &mockFirmInformation{
//				AssignDeputyToFirmErr: returnedError,
//			},
//		},
//	}
//	for k, tc := range tests {
//		t.Run("scenario "+strconv.Itoa(k+1), func(t *testing.T) {
//
//			client := tc.Client
//			template := &mockTemplates{}
//			w := httptest.NewRecorder()
//			r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
//
//			addFirmReturnedError := renderTemplateForAddFirm(client, template)(AppVars{}, w, r)
//			assert.Equal(t, returnedError, addFirmReturnedError)
//		})
//	}
//}
