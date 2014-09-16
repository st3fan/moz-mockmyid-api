package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	setupHandlers("")
}

func Test_handleKey(t *testing.T) {
	request, _ := http.NewRequest("GET", "/key", nil)
	response := httptest.NewRecorder()

	handleKey(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %v", "200", response.Code)
	}
}

func Test_handleAssertion(t *testing.T) {
	request, _ := http.NewRequest("GET", "/assertion?email=test@mockmyid.com&audience=http://localhost", nil)
	response := httptest.NewRecorder()

	handleAssertion(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %v", "200", response.Code)
	}
}
