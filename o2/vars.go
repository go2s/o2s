// authors: wangoo
// created: 2018-05-31
// constants

package o2

import (
	"gopkg.in/oauth2.v3/manage"
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
	oauth2UriUserPass   = "/user/pass"
)

// ---------------------------
var (
	oauth2Svr        *Oauth2Server
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
func GetOauth2Svr() *Oauth2Server {
	return oauth2Svr
}

// expose for custom configuration
func GetOauth2Mgr() *manage.Manager {
	return oauth2Mgr
}
