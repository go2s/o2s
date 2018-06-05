// authors: wangoo
// created: 2018-05-29
// util

package o2

import (
	"net/http"
	"fmt"
	"net/url"
)

func authQuery(r *http.Request) (q string) {
	redirectUri := redirectUri(r)
	if redirectUri != "" {
		redirectUri = url.QueryEscape(redirectUri)
	}
	clientID := clientID(r)
	responseType := responseType(r)
	state := state(r)
	scope := scope(r)

	return fmt.Sprintf("redirect_uri=%v&client_id=%v&response_type=%v&state=%v&scope=%v", redirectUri, clientID, responseType, state, scope)
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	redirectToUri(w, r, Oauth2UriLogin)
}

func redirectToAuth(w http.ResponseWriter, r *http.Request) {
	redirectToUri(w, r, Oauth2UriAuth)
}

func redirectToAuthorize(w http.ResponseWriter, r *http.Request) {
	redirectToUri(w, r, Oauth2UriAuthorize)
}

func redirectToUri(w http.ResponseWriter, r *http.Request, uri string) {
	q := authQuery(r)
	loc := FormatRedirectUri(uri) + "?" + q
	w.Header().Set("Location", loc)
	w.WriteHeader(http.StatusFound)
}

func redirectUri(r *http.Request) string {
	return r.FormValue("redirect_uri")
}

func clientID(r *http.Request) string {
	return r.FormValue("client_id")
}

func state(r *http.Request) string {
	return r.FormValue("state")
}

func scope(r *http.Request) string {
	return r.FormValue("scope")
}

func responseType(r *http.Request) string {
	return r.FormValue("response_type")
}

func username(r *http.Request) string {
	return r.FormValue("username")
}

func password(r *http.Request) string {
	return r.FormValue("password")
}
