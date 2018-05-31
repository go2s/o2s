// authors: wangoo
// created: 2018-05-29
// auth

package o2

import (
	"net/http"
	"gopkg.in/session.v2"
	"fmt"
	"net/url"
	"context"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uid, _ := store.Get(SessionUserID)
	if uid == nil {
		w.Header().Set("Location", oauth2UriFormatter.FormatRedirectUri(Oauth2UriLogin))
		w.WriteHeader(http.StatusFound)
		return
	}

	if r.Method == "POST" {
		param, _ := store.Get(SessionAuthorizeParameters)
		if param == nil {
			http.Error(w, "can get authorize parameters", http.StatusInternalServerError)
			return
		}
		u := new(url.URL)
		u.Path = Oauth2UriAuthorize
		u.RawQuery = param.(string)
		w.Header().Set("Location", oauth2UriFormatter.FormatRedirectUri(u.String()))
		w.WriteHeader(http.StatusFound)
		store.Delete(SessionAuthorizeParameters)
		err = store.Save()
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		return
	}
	outputHTML(w, r, "auth.html")
}
