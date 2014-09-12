// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"github.com/st3fan/moz-go-persona"
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

	// Submit the assertion to the Persona verifier to make sure it is ok

	verifier, err := persona.NewVerifier("https://verifier.login.persona.org/verify", "http://localhost:8080")
	if err != nil {
		t.Error("Could not create Verifier")
		return
	}

	personaResponse, err := verifier.VerifyAssertion(assertion)
	if err != nil {
		t.Error("Could not verify assertion")
		return
	}

	if personaResponse.Status != "okay" {
		t.Errorf("Verifier response is not okay: %s / %s", personaResponse.Status, personaResponse.Reason)
	}
}
