// authors: wangoo
// created: 2018-05-29
// user

package o2

import (
	"net/http"
	"gopkg.in/session.v2"
	"context"
	"github.com/go2s/o2x"
	"github.com/golang/glog"
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
	err = o2x.ErrInvalidCredential
	return
}

// add new user handler
func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	clientID, err := ClientBasicAuth(r)
	if err != nil {
		return
	}
	username := username(r)
	password := password(r)
	if anyNil(username, password) {
		err = o2x.ErrValueRequired
		return
	}
	err = AddUser(clientID, username, password)
	if err != nil {
		data, statusCode, _ := oauth2Svr.GetErrorData(err)
		data["user_id"] = username
		response(w, data, statusCode)
		return
	}
	response(w, defaultSuccessResponse(), http.StatusOK)
	return
}

// add new user
func AddUser(clientID, username, password string) (err error) {
	u, err := oauth2UserStore.Find(username)
	if err != nil && err != o2x.ErrNotFound {
		return
	}
	if u != nil {
		err = o2x.ErrDuplicated
		return
	}

	user := &o2x.SimpleUser{
		UserID: username,
	}
	user.SetRawPassword(password)

	glog.Infof("client %v add user %v", clientID, username)
	err = oauth2UserStore.Save(user)
	if err != nil {
		return
	}

	return
}

// remove user handler
func RemoveUserHandler(w http.ResponseWriter, r *http.Request) {
	clientID, err := ClientBasicAuth(r)
	if err != nil {
		return
	}
	username := username(r)
	if anyNil(username) {
		err = o2x.ErrValueRequired
		return
	}
	err = RemoveUser(clientID, username)
	if err != nil {
		data, statusCode, _ := oauth2Svr.GetErrorData(err)
		response(w, data, statusCode)
		return
	}

	response(w, defaultSuccessResponse(), http.StatusOK)
	return
}

// remove a user
func RemoveUser(clientID, username string) (err error) {
	glog.Infof("client %v remove user %v", clientID, username)
	err = oauth2UserStore.Remove(username)
	if err != nil {
		return
	}
	return
}
