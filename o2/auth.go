// authors: wangoo
// created: 2018-05-29
// auth

package o2

import (
	"net/http"
	"gopkg.in/session.v2"
	"context"
	"github.com/go2s/o2x"
	oauth2Errors "gopkg.in/oauth2.v3/errors"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		ErrorResponse(w, err, http.StatusInternalServerError)
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
	exists := oauth2Svr.authStore.Exist(auth)
	if exists {
		redirectToAuthorize(w, r)
		return
	}

	if r.Method == "POST" {
		oauth2Svr.authStore.Save(auth)
		redirectToAuthorize(w, r)
		return
	}

	m := map[string]interface{}{
		"client": clientID,
		"scope":  scope,
	}

	execAuthTemplate(w, r, m)
}

func ClientBasicAuth(r *http.Request) (cid string, err error) {
	clientID, clientSecret, err := oauth2Svr.ClientInfoHandler(r)
	if err != nil {
		return
	}
	if anyNil(clientID, clientSecret) {
		err = o2x.ErrValueRequired
		return
	}
	cli, err := oauth2Mgr.GetClient(clientID)
	if err != nil {
		return
	}
	if clientSecret != cli.GetSecret() {
		err = oauth2Errors.ErrInvalidClient
		return
	}
	cid = clientID
	return
}
