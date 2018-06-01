// authors: wangoo
// created: 2018-05-29
// auth

package o2

import (
	"net/http"
	"gopkg.in/session.v2"
	"context"
	"log"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	q := authQuery(r)
	uid, _ := store.Get(SessionUserID)
	if uid == nil {
		loc := oauth2UriFormatter.FormatRedirectUri(Oauth2UriLogin) + "?" + q
		w.Header().Set("Location", loc)
		w.WriteHeader(http.StatusFound)
		return
	}

	if r.Method == "POST" {
		loc := oauth2UriFormatter.FormatRedirectUri(Oauth2UriAuthorize) + "?" + q
		w.Header().Set("Location", loc)

		w.WriteHeader(http.StatusFound)
		store.Delete(SessionAuthParam)
		err = store.Save()
		if err != nil {
			log.Printf("failed remove authorize parameters:%v\n", err)
		}
		return
	}
	outputHTML(w, r, "auth.html")
}
