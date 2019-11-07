// authors: wangoo
// created: 2018-05-29
// user

package o2

import (
	"context"
	"net/http"

	"github.com/go2s/o2s/o2x"
	"github.com/golang/glog"
	"gopkg.in/session.v3"

	oauth2Error "gopkg.in/oauth2.v3/errors"
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
	u, err := oauth2Svr.userStore.Find(username)
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
	err := AddUserProcessor(w, r)
	if err != nil {
		data, statusCode, _ := oauth2Svr.GetErrorData(err)
		data["user_id"] = username(r)
		HttpResponse(w, data, statusCode)
		return
	}
	HttpResponse(w, defaultSuccessResponse(), http.StatusOK)
	return
}

// add new user
func AddUserProcessor(w http.ResponseWriter, r *http.Request) (err error) {
	clientID, err := ClientBasicAuth(r)
	if err != nil {
		return
	}
	username := username(r)
	password := password(r)
	scope := scope(r)
	if anyNil(username, password) {
		err = o2x.ErrValueRequired
		return
	}
	u, err := oauth2Svr.userStore.Find(username)
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

	if scope != "" {
		user.Scopes[clientID] = scope
	}

	user.SetRawPassword(password)

	glog.Infof("client %v add user %v", clientID, username)
	return oauth2Svr.userStore.Save(user)
}

// remove user processor
func RemoveUserProcessor(w http.ResponseWriter, r *http.Request) (err error) {
	clientID, err := ClientBasicAuth(r)
	if err != nil {
		return
	}
	username := username(r)
	if anyNil(username) {
		err = o2x.ErrValueRequired
		return
	}

	glog.Infof("client %v remove user %v", clientID, username)
	err = oauth2Svr.userStore.Remove(username)
	if err != nil {
		return
	}
	return
}

// update password processor
func UpdatePwdProcessor(w http.ResponseWriter, r *http.Request) (err error) {
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

	glog.Infof("client %v update password of user %v", clientID, username)
	u, err := oauth2Svr.userStore.Find(username)
	if err != nil {
		return
	}
	err = oauth2Svr.userStore.UpdatePwd(u.GetUserID(), password)
	if err != nil {
		return
	}
	return
}

// update scope processor
func UpdateScopeProcessor(w http.ResponseWriter, r *http.Request) (err error) {
	clientID, err := ClientBasicAuth(r)
	if err != nil {
		return
	}
	username := username(r)
	scope := scope(r)
	if anyNil(username, scope) {
		err = o2x.ErrValueRequired
		return
	}

	glog.Infof("client %v update scope of user %v to %v", clientID, username, scope)
	u, err := oauth2Svr.userStore.Find(username)
	if err != nil {
		return
	}

	allow, err := oauth2Svr.ClientScopeHandler(clientID, scope)
	if err != nil {
		return
	}
	if !allow {
		err = oauth2Error.ErrInvalidScope
		return
	}

	err = oauth2Svr.userStore.UpdateScope(u.GetUserID(), clientID, scope)
	if err != nil {
		return
	}
	return
}
