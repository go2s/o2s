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
	"gopkg.in/session.v2"
	"context"
)

// ---------------------------
var oauth2Svr *server.Server
var oauth2Mgr *manage.Manager
var oauth2Cfg *ServerConfig
var defaultOauth2Cfg *ServerConfig

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

// ---------------------------

func DefaultOauth2Config() *ServerConfig {
	if defaultOauth2Cfg == nil {
		defaultOauth2Cfg = &ServerConfig{
			TemplatePrefix: "../template/",
			ServerName:     "Oauth2 Server",
			Logo:           "https://oauth.net/images/oauth-2-sm.png",
			Favicon:        "https://oauth.net/images/oauth-logo-square.png",
		}
	}
	return defaultOauth2Cfg
}

// ---------------------------

func FormatRedirectUri(uri string) string {
	if oauth2Cfg.UriPrefix != "" {
		return oauth2Cfg.UriPrefix + uri
	}
	return uri
}

// ---------------------------
func InitOauth2Server(cs oauth2.ClientStore, ts oauth2.TokenStore, us o2x.UserStore, as o2x.AuthStore, cfg *ServerConfig) {
	if cs == nil || ts == nil || us == nil {
		panic("store is nil")
	}

	oauth2ClientStore = cs
	oauth2TokenStore = ts
	oauth2UserStore = us
	oauth2AuthStore = as

	if oauth2AuthStore == nil {
		oauth2AuthStore = o2x.NewAuthStore()
	}

	if cfg != nil {
		oauth2Cfg = cfg
	} else {
		oauth2Cfg = DefaultOauth2Config()
	}
	InitTemplate()

	o2xTokenStore, o2xTokenAccountSupport = ts.(o2x.Oauth2TokenStore)

	manager := manage.NewDefaultManager()

	manager.MustTokenStorage(ts, nil)
	manager.MustClientStorage(cs, nil)

	DefaultTokenConfig(manager)

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

func CheckUserAuth(w http.ResponseWriter, r *http.Request) (authorized bool, err error) {
	userID, err := oauth2Svr.UserAuthorizationHandler(w, r)
	if err != nil {
		return
	} else if userID == "" {
		return false, nil
	}

	clientID := clientID(r)
	scope := scope(r)

	if clientID != "" && scope != "" {
		authorized = oauth2AuthStore.Exist(&o2x.AuthModel{
			ClientID: clientID,
			UserID:   userID,
			Scope:    scope,
		})
		return
	}
	return false, nil
}

func AuthorizeRequestHandler(w http.ResponseWriter, r *http.Request) {
	authorized, err := CheckUserAuth(w, r)
	if err != nil || !authorized {
		redirectToAuth(w, r)
		return
	}

	if !multipleUserTokenEnable && o2xTokenAccountSupport && o2xTokenStore != nil {
		responseType := responseType(r)
		if responseType == "token" {
			removeAuthToken(w, r)
		}
	}

	err = oauth2Svr.HandleAuthorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func removeAuthToken(w http.ResponseWriter, r *http.Request) {
	clientID := clientID(r)
	if clientID == "" {
		return
	}
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		return
	}
	u, _ := store.Get(SessionUserID)
	if u == nil {
		return
	}
	userID := u.(string)

	o2xTokenStore.RemoveByAccount(userID, clientID)
}

func BearerTokenValidator(w http.ResponseWriter, r *http.Request) {
	tg, err := oauth2Svr.ValidationBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	data := &o2x.ValidResponse{
		ClientID: tg.GetClientID(),
		UserID:   tg.GetUserID(),
		Scope:    tg.GetScope(),
	}

	body, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Write(body)
}
