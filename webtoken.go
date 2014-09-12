// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"crypto/dsa"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"time"
)

type CertificatePrincipal struct {
	Email string `json:"email"`
}

func NewCertificatePrincipal(email string) *CertificatePrincipal {
	return &CertificatePrincipal{Email: email}
}

type CertificatePublicKey struct {
	Algorithm string `json:"algorithm"`
	y         string `json:"y"`
	p         string `json:"p"`
	q         string `json:"q"`
	g         string `json:"g"`
}

func NewCertificatePublicKey(priv dsa.PrivateKey) *CertificatePublicKey {
	return &CertificatePublicKey{Algorithm: "DS"}
}

type CertificatePayload struct {
	Principal CertificatePrincipal `json:"principal"`
	PublicKey CertificatePublicKey `json:"public-key"`
	Issuer    string               `json:"iss"`
	IssuedAt  int64                `json:"iat"`
	Audience  string               `json:"aud"`
	ExpiresAt int64                `json:"exp"`
}

func NewCertificatePayload(priv dsa.PrivateKey, email string, issuer string, issuedAt time.Time, expiresAt time.Time) *CertificatePayload {
	return &CertificatePayload{
		Principal: *NewCertificatePrincipal(email),
		PublicKey: *NewCertificatePublicKey(priv),
		Issuer:    issuer,
		IssuedAt:  issuedAt.UnixNano() / int64(time.Millisecond),
		ExpiresAt: expiresAt.UnixNano() / int64(time.Millisecond),
	}
}

type WebTokenHeader struct {
	Algorithm string `json:"alg"`
}

func (p *CertificatePayload) Encode(priv dsa.PrivateKey) (string, error) {
	// Encode the header
	header := WebTokenHeader{Algorithm: "DS128"} // TODO: Create a function to extract this from the dsa.PrivateKey
	headerData, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	encodedHeader := base64.URLEncoding.EncodeToString(headerData)

	// Encode the certificate
	payloadData, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	encodedPayload := base64.URLEncoding.EncodeToString(payloadData)

	// The message is simply header plus payload
	message := encodedHeader + "." + encodedPayload

	// Get the SHA1 hash over the message
	messageHash := sha1.New().Sum([]byte(message))

	// Calculate the signature over the message
	r, s, err := dsa.Sign(rand.Reader, &priv, messageHash)
	if err != nil {
		return "", err
	}
	signatureData := append(r.Bytes(), s.Bytes()...)
	encodedSignature := base64.URLEncoding.EncodeToString(signatureData)

	token := message + "." + encodedSignature
	// Prepend

	return token, nil
}

func CreateCertificate(key dsa.PrivateKey, email string, issuer string, issuedAt time.Time, expiresAt time.Time, signingKey dsa.PrivateKey) (string, error) {
	certificatePayload := NewCertificatePayload(key, email, issuer, issuedAt, expiresAt)
	return certificatePayload.Encode(signingKey)
}

//

type AssertionPayload struct {
	Issuer    string `json:"iss"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	Audience  string `json:"aud"`
}

func NewAssertionPayload(issuer string, issuedAt time.Time, expiresAt time.Time, audience string) (*AssertionPayload, error) {
	return &AssertionPayload{
		Issuer:    issuer,
		IssuedAt:  issuedAt.UnixNano() / int64(time.Millisecond),
		ExpiresAt: expiresAt.UnixNano() / int64(time.Millisecond),
		Audience:  audience,
	}, nil
}

func (p *AssertionPayload) Encode(signingKey dsa.PrivateKey) (string, error) {
	// Encode the header
	header := WebTokenHeader{Algorithm: "DS128"} // TODO: Create a function to extract this from the dsa.PrivateKey
	headerData, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	encodedHeader := base64.URLEncoding.EncodeToString(headerData)

	// Encode the certificate
	payloadData, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	encodedPayload := base64.URLEncoding.EncodeToString(payloadData)

	// The message is simply header plus payload
	message := encodedHeader + "." + encodedPayload

	// Get the SHA1 hash over the message
	messageHash := sha1.New().Sum([]byte(message))

	// Calculate the signature over the message
	r, s, err := dsa.Sign(rand.Reader, &signingKey, messageHash)
	if err != nil {
		return "", err
	}
	signatureData := append(r.Bytes(), s.Bytes()...)
	encodedSignature := base64.URLEncoding.EncodeToString(signatureData)

	token := message + "." + encodedSignature
	// Prepend

	return token, nil
}

func CreateAssertion(signingKey dsa.PrivateKey, certificate string, audience string, issuer string, issuedAt time.Time, expiresAt time.Time) (string, error) {
	assertionPayload, err := NewAssertionPayload(issuer, issuedAt, expiresAt, audience)
	if err != nil {
		return "", err
	}
	assertion, err := assertionPayload.Encode(signingKey)
	if err != nil {
		return "", err
	}
	return certificate + "~" + assertion, nil
}
