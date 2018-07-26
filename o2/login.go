// authors: wangoo
// created: 2018-05-29
// login

package o2

import (
	"net/http"
	"gopkg.in/session.v2"
	"context"
	"log"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		username := username(r)
		password := password(r)

		if username == "" || password == "" {
			showLogin(w, r, "username and password required")
			return
		}

		userID, authErr := PasswordAuthorizationHandler(username, password)
		if authErr != nil {
			showLogin(w, r, authErr.Error())
			return
		}

		store.Set(SessionUserID, userID)
		err = store.Save()
		if err != nil {
			return
		}
		log.Printf("login success userID: %v\n", userID)

		redirectToAuth(w, r)
		return
	}

	showLogin(w, r, "")
}

func showLogin(w http.ResponseWriter, r *http.Request, err string) {
	m := map[string]interface{}{
		"error": err,
	}
	execLoginTemplate(w, r, m)
}
