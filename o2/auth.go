// authors: wangoo
// created: 2018-05-29
// auth

package o2

import (
	"net/http"
	"gopkg.in/session.v2"
	"context"
	"github.com/go2s/o2x"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	uid, _ := store.Get(SessionUserID)
	if uid == nil {
		redirectToLogin(w, r)
		return
	}

	clientID := clientID(r)
	scope := scope(r)

	auth := &o2x.AuthModel{
		ClientID: clientID,
		UserID:   uid.(string),
		Scope:    scope,
	}
	exists := oauth2AuthStore.Exist(auth)
	if exists {
		redirectToAuthorize(w, r)
		return
	}

	if r.Method == "POST" {
		oauth2AuthStore.Save(auth)
		redirectToAuthorize(w, r)
		return
	}

	m := map[string]interface{}{
		"client": clientID,
		"scope":  scope,
	}

	execAuthTemplate(w, r, m)
}
