package sirius

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetUserDetailsReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
	"id": 68,
	"roles": ["Finance Manager", "System Admin"]
	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := UserDetails{
		ID:    68,
		Roles: []string{"Finance Manager", "System Admin"},
	}

	userDetails, err := client.GetUserDetails(getContext(nil))

	assert.Equal(t, expectedResponse, userDetails)
	assert.Equal(t, nil, err)
}

func TestUserDetailsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	userDetails, err := client.GetUserDetails(getContext(nil))

	expectedResponse := UserDetails{ID: 0}

	assert.Equal(t, expectedResponse, userDetails)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/users/current",
		Method: http.MethodGet,
	}, err)
}

func TestUserDetailsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	userDetails, err := client.GetUserDetails(getContext(nil))

	expectedResponse := UserDetails{ID: 0}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, userDetails)
}

func TestUserDetails(t *testing.T) {
	t.Run("with Finance Manager", func(t *testing.T) {
		assert.True(t, UserDetails{Roles: []string{"OPG User", "Finance Manager"}}.IsFinanceManager())
	})

	t.Run("without Finance Manager", func(t *testing.T) {
		assert.False(t, UserDetails{Roles: []string{"OPG User", "Case Manager"}}.IsFinanceManager())
	})

	t.Run("with System Admin", func(t *testing.T) {
		assert.True(t, UserDetails{Roles: []string{"OPG User", "System Admin"}}.IsSystemManager())
	})

	t.Run("without System Admin", func(t *testing.T) {
		assert.False(t, UserDetails{Roles: []string{"OPG User", "Finance Manager"}}.IsSystemManager())
	})
}
