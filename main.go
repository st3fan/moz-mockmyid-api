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
	"strings"
)

const (
	MOCKMYID_API_ROOT           = "/"
	MOCKMYID_API_LISTEN_ADDRESS = "127.0.0.1"
	MOCKMYID_API_LISTEN_PORT    = 8080
)

type LoginResponse struct {
	Email     string `json:"email"`
	Audience  string `json:"audience"`
	Assertion string `json:"assertion"`
}

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

func handleLogin(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query["email"][0]
	audience := query["audience"][0]
	uniqueKey := query["uniqueKey"]

	var clientKey *dsa.PrivateKey
	if len(uniqueKey) != 0 {
		key, err := generateRandomKey()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		clientKey = key
	} else {
		clientKey = globalRandomKey
	}

	assertion, err := CreateShortLivedMockMyIDAssertion(*clientKey, email, audience)

	loginResponse := &LoginResponse{
		Email:     email,
		Audience:  audience,
		Assertion: assertion,
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
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	log.Print("Generating client DSA key")
	key, err := generateRandomKey()
	if err != nil {
		panic("Cannot generate random key")
	}
	globalRandomKey = key

	http.HandleFunc(prefix+"/login", handleLogin)
	addr := fmt.Sprintf("%s:%d", *address, *port)
	log.Printf("Starting mockmyid-api server on http://%s%s", addr, prefix)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
