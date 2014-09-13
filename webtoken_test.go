// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"testing"
	"time"
)

func Test_CreateCertificate(t *testing.T) {
	key, err := generateRandomKey()
	if err != nil {
		t.Error("Could not generate random key")
		return
	}

	signingKey, err := generateRandomKey()
	if err != nil {
		t.Error("Could not generate random key")
		return
	}

	token, err := CreateCertificate(*key, "test@mockmyid.com", "https://mockmyid.com", time.Now(), time.Now(), *signingKey)
	if err != nil {
		t.Error("Could not create certificate")
	}
	if len(token) == 0 {
		t.Error("Empty token created")
	}
}
