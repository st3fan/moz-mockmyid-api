package main

import (
	"testing"
)

func Test_createShortLivedMockMyIDAssertion(t *testing.T) {
	key, err := generateRandomKey()
	if err != nil {
		t.Error("Could not generate random key")
		return
	}
	assertion, err := CreateShortLivedMockMyIDAssertion(*key, "test@mockmyid.com", "http://localhost:8080")
	if err != nil {
		t.Error("Could not create assertion")
		return
	}
	if len(assertion) == 0 {
		t.Error("Got an empty assertion")
	}
}
