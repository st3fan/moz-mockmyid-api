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
	"fmt"
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
	Y         string `json:"y"`
	P         string `json:"p"`
	Q         string `json:"q"`
	G         string `json:"g"`
}

func NewCertificatePublicKey(key dsa.PrivateKey) *CertificatePublicKey {
	return &CertificatePublicKey{
		Algorithm: "DS",
		Y:         fmt.Sprintf("%x", key.PublicKey.Y),
		P:         fmt.Sprintf("%x", key.PublicKey.Parameters.P),
		Q:         fmt.Sprintf("%x", key.PublicKey.Parameters.Q),
		G:         fmt.Sprintf("%x", key.PublicKey.Parameters.G),
	}
}

type CertificatePayload struct {
	Principal CertificatePrincipal `json:"principal"`
	PublicKey CertificatePublicKey `json:"public-key"`
	Issuer    string               `json:"iss"`
	IssuedAt  int64                `json:"iat"`
	Audience  string               `json:"aud"`
	ExpiresAt int64                `json:"exp"`
}

func NewCertificatePayload(publicKey dsa.PrivateKey, email string, issuer string, issuedAt time.Time, expiresAt time.Time) *CertificatePayload {
	return &CertificatePayload{
		Principal: *NewCertificatePrincipal(email),
		PublicKey: *NewCertificatePublicKey(publicKey),
		Issuer:    issuer,
		IssuedAt:  issuedAt.UnixNano() / int64(time.Millisecond),
		ExpiresAt: expiresAt.UnixNano() / int64(time.Millisecond),
	}
}

type WebTokenHeader struct {
	Algorithm string `json:"alg"`
}

func (p *CertificatePayload) Encode(signingKey dsa.PrivateKey) (string, error) {
	// Encode the header
	header := WebTokenHeader{Algorithm: fmt.Sprintf("DS%d", (signingKey.PublicKey.Parameters.P.BitLen()+7)/8)}
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
	h := sha1.New()
	h.Write([]byte(message))
	messageHash := h.Sum(nil)

	// Calculate the signature over the message
	r, s, err := dsa.Sign(rand.Reader, &signingKey, messageHash)
	if err != nil {
		return "", err
	}

	signatureData := make([]byte, 40)
	rb := r.Bytes()
	sb := s.Bytes()
	copy(signatureData[20-len(rb):20], rb)
	copy(signatureData[40-len(sb):], sb)

	//signatureData := append(r.Bytes(), s.Bytes()...)

	encodedSignature := base64.URLEncoding.EncodeToString(signatureData)

	token := message + "." + encodedSignature
	// Prepend

	return token, nil
}

func CreateCertificate(publicKey dsa.PrivateKey, email string, issuer string, issuedAt time.Time, expiresAt time.Time, signingKey dsa.PrivateKey) (string, error) {
	certificatePayload := NewCertificatePayload(publicKey, email, issuer, issuedAt, expiresAt)
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
	h := sha1.New()
	h.Write([]byte(message))
	messageHash := h.Sum(nil)

	// Calculate the signature over the message
	r, s, err := dsa.Sign(rand.Reader, &signingKey, messageHash)
	if err != nil {
		return "", err
	}

	signatureData := make([]byte, 40)
	rb := r.Bytes()
	sb := s.Bytes()
	copy(signatureData[20-len(rb):20], rb)
	copy(signatureData[40-len(sb):], sb)

	//signatureData := append(r.Bytes(), s.Bytes()...)

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
