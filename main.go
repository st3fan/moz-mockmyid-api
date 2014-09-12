// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	MOCKMYID_API_ROOT           = ""
	MOCKMYID_API_LISTEN_ADDRESS = "0.0.0.0"
	MOCKMYID_API_LISTEN_PORT    = 8124
	MOCKMYID_DOMAIN             = "mockmyid.com"
)

type LoginResponse struct {
	Email     string `json:"email"`
	Assertion string `json:"assertion"`
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query["email"][0]

	loginResponse := &LoginResponse{
		Email:     email,
		Assertion: "cheesefooblah",
	}

	encodedLoginResponse, err := json.Marshal(loginResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(encodedLoginResponse)
}

func main() {
	http.HandleFunc(MOCKMYID_API_ROOT+"/login", handleLogin)
	addr := fmt.Sprintf("%s:%d", MOCKMYID_API_LISTEN_ADDRESS, MOCKMYID_API_LISTEN_PORT)
	log.Printf("Starting tokenserver server on http://%s%s", addr, MOCKMYID_API_ROOT)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
