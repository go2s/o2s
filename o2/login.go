// authors: wangoo
// created: 2018-05-29
// login

package o2

import (
	"net/http"
	"gopkg.in/session.v2"
	"fmt"
	"context"
	"log"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "username and password required", http.StatusInternalServerError)
			return
		}
		log.Printf("login request username: %v\n", username)

		userID, err := PasswordAuthorizationHandler(username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		store.Set(SessionUserID, userID)
		err = store.Save()
		if err != nil {
			log.Printf("login failed: %v\n", err)
			fmt.Fprint(w, err)
			return
		}
		log.Printf("login success userID: %v\n", userID)

		q := authQuery(r)

		loc := oauth2UriFormatter.FormatRedirectUri(Oauth2UriAuth) + "?" + q
		w.Header().Set("Location", loc)
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, r, "login.html")
}
