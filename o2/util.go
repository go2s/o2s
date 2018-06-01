// authors: wangoo
// created: 2018-05-29
// util

package o2

import (
	"net/http"
	"os"
	"fmt"
	"net/url"
)

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}

func authQuery(r *http.Request) (q string) {
	redirectUri := r.FormValue("redirect_uri")
	if redirectUri != "" {
		redirectUri = url.QueryEscape(redirectUri)
	}
	clientID := r.FormValue("client_id")
	responseType := r.FormValue("response_type")
	state := r.FormValue("state")
	scope := r.FormValue("scope")

	return fmt.Sprintf("redirect_uri=%v&client_id=%v&response_type=%v&state=%v&scope=%v", redirectUri, clientID, responseType, state, scope)
}
