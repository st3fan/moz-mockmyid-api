// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"crypto/dsa"
	"crypto/rand"
	"log"
	"testing"
	"time"
)

func generateRandomKey() (*dsa.PrivateKey, error) {
	params := new(dsa.Parameters)
	if err := dsa.GenerateParameters(params, rand.Reader, dsa.L1024N160); err != nil {
		return nil, err
	}
	priv := new(dsa.PrivateKey)
	priv.PublicKey.Parameters = *params
	if err := dsa.GenerateKey(priv, rand.Reader); err != nil {
		return nil, err
	}
	return priv, nil
}

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
	log.Printf("Token: %s", token)
}
