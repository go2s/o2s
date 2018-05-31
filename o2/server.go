// authors: wangoo
// created: 2018-05-21
// ouath2 server demo based on redis storage

package o2

import (
	"log"
	"net/http"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3"

	"encoding/json"
	"github.com/go2s/o2x"
)

var oauth2Svr *server.Server
var oauth2Mgr *manage.Manager

var oauth2ClientStore oauth2.ClientStore
var oauth2TokenStore oauth2.TokenStore
var oauth2UserStore o2x.UserStore
var oauth2UriFormatter o2x.UriFormatter

type DefaultUriFormatter struct {
}

func (u *DefaultUriFormatter) FormatRedirectUri(uri string) string {
	return uri
}

func InitOauth2Server(cs oauth2.ClientStore, ts oauth2.TokenStore, us o2x.UserStore, formatter o2x.UriFormatter) {
	oauth2ClientStore = cs
	oauth2TokenStore = ts
	oauth2UserStore = us
	oauth2UriFormatter = formatter

	if oauth2UriFormatter == nil {
		oauth2UriFormatter = &DefaultUriFormatter{}
	}

	manager := manage.NewDefaultManager()

	manager.MustTokenStorage(ts, nil)
	manager.MustClientStorage(cs, nil)

	TokenConfig(manager)

	oauth2Svr = server.NewServer(&server.Config{
		TokenType:            "Bearer",
		AllowedResponseTypes: []oauth2.ResponseType{oauth2.Code, oauth2.Token},
		AllowedGrantTypes: []oauth2.GrantType{
			oauth2.AuthorizationCode,
			oauth2.PasswordCredentials,
			oauth2.ClientCredentials,
			oauth2.Refreshing,
			oauth2.Implicit,
		},
	}, manager)

	oauth2Svr.SetAllowGetAccessRequest(true)
	oauth2Svr.SetClientInfoHandler(server.ClientBasicHandler)
	oauth2Svr.SetPasswordAuthorizationHandler(PasswordAuthorizationHandler)
	oauth2Svr.SetUserAuthorizationHandler(userAuthorizeHandler)

	oauth2Svr.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	oauth2Svr.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

}

func TokenRequestHandler(w http.ResponseWriter, r *http.Request) {
	oauth2Svr.HandleTokenRequest(w, r)
}

func AuthorizeRequestHandler(w http.ResponseWriter, r *http.Request) {
	err := oauth2Svr.HandleAuthorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func BearerTokenValidator(w http.ResponseWriter, r *http.Request) {
	tg, err := oauth2Svr.ValidationBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	body, err := json.Marshal(tg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Write(body)
}
