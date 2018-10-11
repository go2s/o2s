// authors: wangoo
// created: 2018-05-29
// util

package o2

import (
	"net/http"
	"fmt"
	"net/url"
	"encoding/json"
	"log"
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

	if redirectUri == "" || clientID == "" || responseType == "" || scope == "" {
		return
	}

	q = fmt.Sprintf("redirect_uri=%v&client_id=%v&response_type=%v&state=%v&scope=%v", redirectUri, clientID, responseType, state, scope)
	return
}

// ---------------------------
func FormatRedirectUri(uri string) string {
	if oauth2Cfg.URIPrefix != "" {
		return oauth2Cfg.URIPrefix + oauth2Cfg.URIContext + uri
	}
	return oauth2Cfg.URIContext + uri
}

func redirectToIndex(w http.ResponseWriter, r *http.Request) {
	redirectToUri(w, r, oauth2UriIndex, "")
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	authRedirect(w, r, oauth2UriLogin)
}

func redirectToAuth(w http.ResponseWriter, r *http.Request) {
	authRedirect(w, r, oauth2UriAuth)
}

func redirectToAuthorize(w http.ResponseWriter, r *http.Request) {
	authRedirect(w, r, oauth2UriAuthorize)
}

func authRedirect(w http.ResponseWriter, r *http.Request, uri string) {
	q := authQuery(r)
	if uri != oauth2UriLogin && q == "" {
		redirectToIndex(w, r)
	} else {
		redirectToUri(w, r, uri, q)
	}
}

func redirectToUri(w http.ResponseWriter, r *http.Request, uri string, query string) {
	loc := FormatRedirectUri(uri)
	if query != "" {
		loc += "?" + query
	}
	w.Header().Set("Location", loc)
	w.WriteHeader(http.StatusFound)
}

func ErrorResponse(w http.ResponseWriter, err error, status int) {
	HttpResponse(w, defaultErrorResponse(err), status)
}

func HttpResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
	}
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

func anyNil(args ... string) bool {
	for _, arg := range args {
		if arg == "" {
			return true
		}
	}
	return false
}
