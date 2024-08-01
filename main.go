package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"encoding/base64"
	"math/rand"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

var log = zap.Must(zap.NewDevelopment()).Sugar()

const (
	clientID     = "<client-id>"
	clientSecret = "<client-secret"
	redirectURI  = "http://localhost:8080/callback" // Replace with your actual redirect URI
	tenantID     = "<tenant-id>"
)

var fscopes = []string{"user.read", "profile", "email", oidc.ScopeOpenID}

func main() {
	config := newConfig(clientID, clientSecret, redirectURI, tenantID)
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/callback", handleAuthCode(config))

	log.Info("Listening on port 8080..")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func newConfig(clientID, clientSecret, redirectURI, tenantID string) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  redirectURI,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       fscopes,
		Endpoint:     endpoints.AzureAD(tenantID),
	}
}

type Response struct {
	Message string `json:"message"`
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	config := newConfig(clientID, clientSecret, redirectURI, tenantID)
	state := generateStateOauthCookie(w)
	// authCode := config.AuthCodeURL(state)
	// // Set the content type header
	// w.Header().Set("Content-Type", "application/json")

	// resp := Response{Message: authCode}

	// // Encode the response to JSON
	// jsonResponse, err := json.Marshal(resp)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// Write the JSON response to the response writer
	// w.Write(jsonResponse)
	http.Redirect(w, r, config.AuthCodeURL(state), http.StatusFound)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func handleAuthCode(config *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Callback invoked")
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "FAILED REQUEST", 499)
			return
		}

		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, ok := token.Extra("id_token").(string)
		if !ok {
			//handle missing token
		}

		// log.Debug("This is the id token", idToken)

		client := config.Client(context.Background(), token)
		resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var user interface{}
		err = json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Hello, %s!\n", user)
	}
}
