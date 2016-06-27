package main

import "net/http"

func login(w http.ResponseWriter, r *http.Request, err string) {
	if r.Header.Get("Authorization") == "" {

		http.Redirect(w, r, "https://login.example.com?redirect_uri=", 302)
	} else {

		http.Error(w, err, http.StatusUnauthorized)
	}
}
