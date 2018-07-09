// authors: wangoo
// created: 2018-05-29
// user

package o2

import (
	"net/http"
	oauth2_errors "gopkg.in/oauth2.v3/errors"
	"gopkg.in/session.v2"
	"context"
	"github.com/go2s/o2x"
)

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		return
	}
	uid, _ := store.Get(SessionUserID)
	if uid == nil {
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
		uid := u.GetUserID()
		return o2x.UserIdString(uid)
	}
	err = ErrInvalidCredential
	return
}

// add new user
func AddUser(w http.ResponseWriter, r *http.Request) (err error) {
	clientID, clientSecret, err := oauth2Svr.ClientInfoHandler(r)
	if err != nil {
		return
	}
	username := username(r)
	password := password(r)
	if anyNil(clientID, clientSecret, username, password) {
		err = ErrValueRequired
		return
	}
	cli, err := oauth2Mgr.GetClient(clientID)
	if err != nil {
		return
	}
	if clientSecret != cli.GetSecret() {
		err = oauth2_errors.ErrInvalidClient
		return
	}

	u, err := oauth2UserStore.Find(username)
	if err != nil {
		return
	}
	if u != nil {
		data := defaultErrorResponse(ErrDuplicated)
		data["user_id"] = u.GetUserID()
		response(w, data, http.StatusConflict)
		return
	}

	user := &o2x.SimpleUser{
		UserID: username,
	}
	user.SetRawPassword(password)
	err = oauth2UserStore.Save(user)
	if err != nil {
		return
	}

	response(w, defaultSuccessResponse(), http.StatusOK)
	return
}
