// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"crypto/dsa"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

const (
	MOCKMYID_API_ROOT           = "/"
	MOCKMYID_API_LISTEN_ADDRESS = "127.0.0.1"
	MOCKMYID_API_LISTEN_PORT    = 8080
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

type KeyResponse struct {
	Algorithm string `json:"algorithm"`
	X         string `json:"x"`
	Y         string `json:"y"`
	P         string `json:"p"`
	Q         string `json:"q"`
	G         string `json:"g"`
}

func handleKey(w http.ResponseWriter, r *http.Request) {
	keyResponse := KeyResponse{
		Algorithm: "DS",
		X:         fmt.Sprintf("%x", MOCKMYID_KEY.X),
		Y:         fmt.Sprintf("%x", MOCKMYID_KEY.PublicKey.Y),
		P:         fmt.Sprintf("%x", MOCKMYID_KEY.PublicKey.Parameters.P),
		Q:         fmt.Sprintf("%x", MOCKMYID_KEY.PublicKey.Parameters.Q),
		G:         fmt.Sprintf("%x", MOCKMYID_KEY.PublicKey.Parameters.G),
	}

	encodedKeyResponse, err := json.Marshal(keyResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(encodedKeyResponse)
}

type LoginResponse struct {
	Email     string      `json:"email"`
	Audience  string      `json:"audience"`
	Assertion string      `json:"assertion"`
	ClientKey KeyResponse `json:"clientKey"`
}

func handleAssertion(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Get and verify parameters

	if query["email"] == nil {
		http.Error(w, "No email parameter supplied", http.StatusBadRequest)
		return
	}
	email := query["email"][0]

	if query["audience"] == nil {
		http.Error(w, "No audience parameter supplied", http.StatusBadRequest)
		return
	}
	audience := query["audience"][0]

	// Check if the caller wants a unique key

	uniqueClientKey := query["uniqueClientKey"]

	var clientKey *dsa.PrivateKey
	if uniqueClientKey != nil {
		key, err := generateRandomKey()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		clientKey = key
	} else {
		clientKey = globalRandomKey
	}

	// Create an assertion and return it together with the key used

	assertion, err := CreateShortLivedMockMyIDAssertion(*clientKey, email, audience)

	loginResponse := &LoginResponse{
		Email:     email,
		Audience:  audience,
		Assertion: assertion,
		ClientKey: KeyResponse{
			Algorithm: "DS",
			X:         fmt.Sprintf("%x", clientKey.X),
			Y:         fmt.Sprintf("%x", clientKey.PublicKey.Y),
			P:         fmt.Sprintf("%x", clientKey.PublicKey.Parameters.P),
			Q:         fmt.Sprintf("%x", clientKey.PublicKey.Parameters.Q),
			G:         fmt.Sprintf("%x", clientKey.PublicKey.Parameters.G),
		},
	}

	encodedLoginResponse, err := json.Marshal(loginResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(encodedLoginResponse)
}

var globalRandomKey *dsa.PrivateKey

func main() {
	address := flag.String("address", MOCKMYID_API_LISTEN_ADDRESS, "address to listen on")
	port := flag.Int("port", MOCKMYID_API_LISTEN_PORT, "port to listen on")
	root := flag.String("root", MOCKMYID_API_ROOT, "application root (path prefix)")
	flag.Parse()

	prefix := *root
	if prefix == "/" {
		prefix = ""
	}

	log.Print("Generating client DSA key")
	key, err := generateRandomKey()
	if err != nil {
		panic("Cannot generate random key")
	}
	globalRandomKey = key

	http.HandleFunc(prefix+"/assertion", handleAssertion)
	http.HandleFunc(prefix+"/key", handleKey)

	addr := fmt.Sprintf("%s:%d", *address, *port)
	log.Printf("Starting mockmyid-api server on http://%s%s", addr, prefix)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
