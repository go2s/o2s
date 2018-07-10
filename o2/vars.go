// authors: wangoo
// created: 2018-05-31
// constants

package o2

import (
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3"
	"github.com/go2s/o2x"
)

const (
	SessionUserID = "UserID"

	oauth2UriIndex      = "/index"
	oauth2UriLogin      = "/login"
	oauth2UriAuth       = "/auth"
	oauth2UriAuthorize  = "/authorize"
	oauth2UriToken      = "/token"
	oauth2UriValid      = "/valid"
	oauth2UriUserAdd    = "/user"
	oauth2UriUserRemove = "/user/remove"
)

// ---------------------------
var (
	oauth2Svr        *server.Server
	oauth2Mgr        *manage.Manager
	oauth2Cfg        *ServerConfig
	defaultOauth2Cfg *ServerConfig
)

func defaultSuccessResponse() map[string]interface{} {
	data := map[string]interface{}{
		"error": "ok",
	}
	return data
}

func defaultErrorResponse(err error) map[string]interface{} {
	data := map[string]interface{}{
		"error":             "server_error",
		"error_description": err.Error(),
	}
	if e, ok := err.(o2x.CodeError); ok {
		data["error"] = e.Code()
	}
	return data
}

// expose for custom configuration
func GetOauth2Svr() *server.Server {
	return oauth2Svr
}

// expose for custom configuration
func GetOauth2Mgr() *manage.Manager {
	return oauth2Mgr
}

// ---------------------------
var oauth2ClientStore oauth2.ClientStore
var oauth2TokenStore oauth2.TokenStore
var oauth2UserStore o2x.UserStore
var oauth2AuthStore o2x.AuthStore

// ---------------------------
// whether the token store support account management
var o2xTokenAccountSupport = false
var o2xTokenStore o2x.Oauth2TokenStore

// ---------------------------
// enable to create multiple token for one user of a client
var multipleUserTokenEnable = false

func EnableMultipleUserToken() {
	multipleUserTokenEnable = true
}

func DisableMultipleUserToken() {
	multipleUserTokenEnable = false
}
