// authors: wangoo
// created: 2018-05-21
// ouath2 server demo based on redis storage

package o2

import (
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3"

	"github.com/go2s/o2x"
	"net/http"
	oauth2Error "gopkg.in/oauth2.v3/errors"
	"github.com/golang/glog"
	"encoding/json"
)

type Oauth2Server struct {
	*server.Server
}

// NewServer create authorization server
func NewServer(cfg *server.Config, manager oauth2.Manager) *Oauth2Server {
	svr := server.NewServer(cfg, manager)
	o2svr := &Oauth2Server{
		svr,
	}

	return o2svr
}

// ValidationTokenRequest the token request validation, add user client scope validation
func (s *Oauth2Server) ValidationTokenRequest(r *http.Request) (gt oauth2.GrantType, tgr *oauth2.TokenGenerateRequest, err error) {
	gt, tgr, err = s.Server.ValidationTokenRequest(r)
	if err != nil {
		return
	}

	if gt != oauth2.PasswordCredentials || tgr.Scope == "" {
		return
	}

	user, err := oauth2UserStore.Find(tgr.UserID)
	if err != nil {
		return
	}
	scope, ok := user.GetScopes()[tgr.ClientID]
	if ok && o2x.ScopeContains(scope, tgr.Scope) {
		return
	}
	glog.Errorf("the scope of user [%v] for client [%v] is [%v], but request [%v]", tgr.UserID, tgr.ClientID, scope, tgr.Scope)
	err = oauth2Error.ErrInvalidScope
	return
}

// HandleTokenRequest token request handling
func (s *Oauth2Server) HandleTokenRequest(w http.ResponseWriter, r *http.Request) (err error) {
	gt, tgr, verr := s.ValidationTokenRequest(r)
	if verr != nil {
		err = s.tokenError(w, verr)
		return
	}

	ti, verr := s.GetAccessToken(gt, tgr)
	if verr != nil {
		err = s.tokenError(w, verr)
		return
	}

	err = s.token(w, s.GetTokenData(ti), nil)
	return
}

// override Server.tokenError
func (s *Oauth2Server) tokenError(w http.ResponseWriter, err error) (uerr error) {
	data, statusCode, header := s.GetErrorData(err)

	uerr = s.token(w, data, header, statusCode)
	return
}

// override Server.token
func (s *Oauth2Server) token(w http.ResponseWriter, data map[string]interface{}, header http.Header, statusCode ...int) (err error) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	for key := range header {
		w.Header().Set(key, header.Get(key))
	}

	status := http.StatusOK
	if len(statusCode) > 0 && statusCode[0] > 0 {
		status = statusCode[0]
	}

	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(data)
	return
}

// ---------------------------
func InitOauth2Server(cs oauth2.ClientStore, ts oauth2.TokenStore, us o2x.UserStore, as o2x.AuthStore, cfg *ServerConfig, mapper HandleMapper) {
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

	InitServerConfig(cfg, mapper)

	o2xTokenStore, o2xTokenAccountSupport = ts.(o2x.Oauth2TokenStore)

	oauth2Mgr = manage.NewDefaultManager()

	oauth2Mgr.MustTokenStorage(ts, nil)
	oauth2Mgr.MustClientStorage(cs, nil)

	DefaultTokenConfig(oauth2Mgr)

	oauth2Svr = NewServer(&server.Config{
		TokenType:            "Bearer",
		AllowedResponseTypes: []oauth2.ResponseType{oauth2.Code, oauth2.Token},
		AllowedGrantTypes: []oauth2.GrantType{
			oauth2.AuthorizationCode,
			oauth2.PasswordCredentials,
			oauth2.ClientCredentials,
			oauth2.Refreshing,
			oauth2.Implicit,
		},
	}, oauth2Mgr)

	oauth2Svr.SetAllowGetAccessRequest(true)
	oauth2Svr.SetClientInfoHandler(server.ClientBasicHandler)
	oauth2Svr.SetPasswordAuthorizationHandler(PasswordAuthorizationHandler)
	oauth2Svr.SetUserAuthorizationHandler(userAuthorizeHandler)
	oauth2Svr.SetInternalErrorHandler(InternalErrorHandler)
	oauth2Svr.SetResponseErrorHandler(ResponseErrorHandler)
	oauth2Svr.SetClientScopeHandler(ClientScopeHandler)
	oauth2Svr.SetRefreshingScopeHandler(RefreshingScopeHandler)
	oauth2Svr.SetAuthorizeScopeHandler(AuthorizeScopeHandler)
}
