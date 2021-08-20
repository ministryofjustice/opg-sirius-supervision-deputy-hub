package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(t *testing.T) {
	t.Run("returns hello world", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		HelloServer(response, request)

		got := response.Body.String()
		want := "Hello world!"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
