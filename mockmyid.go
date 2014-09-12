// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"crypto/dsa"
	"math/big"
	"strings"
	"time"
)

const (
	CERTIFICATE_DURATION      = time.Duration(60*60) * time.Second
	ASSERTION_DURATION        = time.Duration(60*60) * time.Second
	CERTIFICATE_ISSUED_OFFSET = time.Duration(30) * time.Second
	ASSERTION_ISSUES_OFFSET   = time.Duration(15) * time.Second
)

func createMockMyIDCertificate(key dsa.PrivateKey, username string, issuedAt time.Time, duration time.Duration) (string, error) {
	if !strings.HasSuffix(username, "@mockmyid.com") {
		username = username + "@mockmyid.com"
	}

	expiresAt := issuedAt.Add(duration)
	return CreateCertificate(key, username, "@mockmyid.com", issuedAt, expiresAt, MOCKMYID_KEY) // From webtoken.go
}

func CreateMockMyIDAssertion(key dsa.PrivateKey, username, audience string, certificateIssuedAt time.Time, certificateDuration time.Duration, assertionIssuedAt time.Time, assertionDuration time.Duration) (string, error) {
	certificate, err := createMockMyIDCertificate(key, username, certificateIssuedAt, certificateDuration)
	if err != nil {
		return "", err
	}
	assertionExpiresAt := assertionIssuedAt.Add(assertionDuration)
	return CreateAssertion(key, certificate, audience, "127.0.0.1", assertionIssuedAt, assertionExpiresAt) // From webtoken.go
}

func CreateShortLivedMockMyIDAssertion(key dsa.PrivateKey, username, audience string) (string, error) {
	now := time.Now()
	return CreateMockMyIDAssertion(key, username, audience, now.Add(-CERTIFICATE_ISSUED_OFFSET), CERTIFICATE_DURATION, now.Add(-ASSERTION_ISSUES_OFFSET), ASSERTION_DURATION)
}

//

var MOCKMYID_KEY dsa.PrivateKey

func stringToBig(s string) *big.Int {
	n := new(big.Int)
	n.SetString(s, 16)
	return n
}

func init() {
	// TODO: It would be nice if we could fetch this from the web
	MOCKMYID_KEY = dsa.PrivateKey{
		PublicKey: dsa.PublicKey{
			Parameters: dsa.Parameters{
				P: stringToBig("ff600483db6abfc5b45eab78594b3533d550d9f1bf2a992a7a8daa6dc34f8045ad4e6e0c429d334eeeaaefd7e23d4810be00e4cc1492cba325ba81ff2d5a5b305a8d17eb3bf4a06a349d392e00d329744a5179380344e82a18c47933438f891e22aeef812d69c8f75e326cb70ea000c3f776dfdbd604638c2ef717fc26d02e17"),
				Q: stringToBig("e21e04f911d1ed7991008ecaab3bf775984309c3"),
				G: stringToBig("c52a4a0ff3b7e61fdf1867ce84138369a6154f4afa92966e3c827e25cfa6cf508b90e5de419e1337e07a2e9e2a3cd5dea704d175f8ebf6af397d69e110b96afb17c7a03259329e4829b0d03bbc7896b15b4ade53e130858cc34d96269aa89041f409136c7242a38895c9d5bccad4f389af1d7a4bd1398bd072dffa896233397a"),
			},
			Y: stringToBig("738ec929b559b604a232a9b55a5295afc368063bb9c20fac4e53a74970a4db7956d48e4c7ed523405f629b4cc83062f13029c4d615bbacb8b97f5e56f0c7ac9bc1d4e23809889fa061425c984061fca1826040c399715ce7ed385c4dd0d402256912451e03452d3c961614eb458f188e3e8d2782916c43dbe2e571251ce38262"),
		},
		X: stringToBig("385cb3509f086e110c5e24bdd395a84b335a09ae"),
	}
}
