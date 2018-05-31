// authors: wangoo
// created: 2018-05-29
// user

package o2

import (
	"net/http"
	"gopkg.in/session.v2"
	"context"
	"errors"
)

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(context.Background(), w, r)
	uid, _ := store.Get(SessionUserID)
	if uid == nil {
		if r.Form == nil {
			r.ParseForm()
		}
		f := r.Form
		store.Set(SessionAuthorizeParameters, f.Encode())
		err = store.Save()
		if err != nil {
			return
		}
		w.Header().Set("Location", oauth2UriFormatter.FormatRedirectUri(Oauth2UriLogin))
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	return
}

func PasswordAuthorizationHandler(username, password string) (userID string, err error) {
	u, err := oauth2UserStore.Find(username)
	if err != nil {
		return
	}
	if u != nil && u.Match(password) {
		userID = u.UserID
		return
	}
	err = errors.New("invalid user or password")
	return
}
