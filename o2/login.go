// authors: wangoo
// created: 2018-05-29
// login

package o2

import (
	"net/http"
	"gopkg.in/session.v2"
	"fmt"
	"context"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		store, err := session.Start(context.Background(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "username and password required", http.StatusInternalServerError)
			return
		}

		userID, err := PasswordAuthorizationHandler(username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		store.Set(SessionUserID, userID)
		err = store.Save()
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		w.Header().Set("Location", oauth2UriFormatter.FormatRedirectUri(Oauth2UriAuth))
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, r, "login.html")
}
