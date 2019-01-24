// authors: wangoo
// created: 2018-05-31
// mongodb token

package o2m

import (
	"time"
	"gopkg.in/oauth2.v3"
)

type TokenData struct {
	Access           string        `bson:"_id" json:"Access"`
	ClientID         string        `bson:"ClientId" json:"ClientId"`
	UserID           string        `bson:"UserID" json:"UserID"`
	RedirectURI      string        `bson:"RedirectURI,omitempty" json:"RedirectURI,omitempty"`
	Scope            string        `bson:"Scope,omitempty" json:"Scope,omitempty"`
	Code             string        `bson:"Code,omitempty" json:"Code,omitempty"`
	CodeCreateAt     time.Time     `bson:"CodeCreateAt" json:"CodeCreateAt"`
	CodeExpiresIn    time.Duration `bson:"CodeExpiresIn" json:"CodeExpiresIn"`
	AccessCreateAt   time.Time     `bson:"AccessCreateAt" json:"AccessCreateAt"`
	AccessExpiresIn  time.Duration `bson:"AccessExpiresIn" json:"AccessExpiresIn"`
	Refresh          string        `bson:"Refresh,omitempty" json:"Refresh,omitempty"`
	RefreshCreateAt  time.Time     `bson:"RefreshCreateAt,omitempty" json:"RefreshCreateAt,omitempty"`
	RefreshExpiresIn time.Duration `bson:"RefreshExpiresIn,omitempty" json:"RefreshExpiresIn,omitempty"`
	ExpiredAt        time.Time     `bson:"ExpiredAt" json:"ExpiredAt"`
}

func Copy(info oauth2.TokenInfo) (token *TokenData) {
	token = &TokenData{
		ClientID:         info.GetClientID(),
		UserID:           info.GetUserID(),
		RedirectURI:      info.GetRedirectURI(),
		Scope:            info.GetScope(),
		Code:             info.GetCode(),
		CodeCreateAt:     info.GetCodeCreateAt(),
		CodeExpiresIn:    info.GetCodeExpiresIn(),
		Access:           info.GetAccess(),
		AccessCreateAt:   info.GetAccessCreateAt(),
		AccessExpiresIn:  info.GetAccessExpiresIn(),
		Refresh:          info.GetRefresh(),
		RefreshCreateAt:  info.GetRefreshCreateAt(),
		RefreshExpiresIn: info.GetRefreshExpiresIn(),
	}
	if code := info.GetCode(); code != "" {
		token.ExpiredAt = info.GetCodeCreateAt().Add(info.GetCodeExpiresIn())
	} else {
		aexp := info.GetAccessCreateAt().Add(info.GetAccessExpiresIn())
		rexp := aexp
		if info.GetRefresh() != "" && info.GetRefreshExpiresIn() > 0 {
			rexp = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn())
			if aexp.Second() > rexp.Second() {
				aexp = rexp
			}
		}
		token.ExpiredAt = rexp
	}
	return
}

// New create to token model instance
func (t *TokenData) New() oauth2.TokenInfo {
	return &TokenData{}
}

// GetClientID the client id
func (t *TokenData) GetClientID() string {
	return t.ClientID
}

// SetClientID the client id
func (t *TokenData) SetClientID(clientID string) {
	t.ClientID = clientID
}

// GetUserID the user id
func (t *TokenData) GetUserID() string {
	return t.UserID
}

// SetUserID the user id
func (t *TokenData) SetUserID(userID string) {
	t.UserID = userID
}

// GetRedirectURI redirect URI
func (t *TokenData) GetRedirectURI() string {
	return t.RedirectURI
}

// SetRedirectURI redirect URI
func (t *TokenData) SetRedirectURI(redirectURI string) {
	t.RedirectURI = redirectURI
}

// GetScope get scope of authorization
func (t *TokenData) GetScope() string {
	return t.Scope
}

// SetScope get scope of authorization
func (t *TokenData) SetScope(scope string) {
	t.Scope = scope
}

// GetCode authorization code
func (t *TokenData) GetCode() string {
	return t.Code
}

// SetCode authorization code
func (t *TokenData) SetCode(code string) {
	t.Code = code
}

// GetCodeCreateAt create Time
func (t *TokenData) GetCodeCreateAt() time.Time {
	return t.CodeCreateAt
}

// SetCodeCreateAt create Time
func (t *TokenData) SetCodeCreateAt(createAt time.Time) {
	t.CodeCreateAt = createAt
}

// GetCodeExpiresIn the lifetime in seconds of the authorization code
func (t *TokenData) GetCodeExpiresIn() time.Duration {
	return t.CodeExpiresIn
}

// SetCodeExpiresIn the lifetime in seconds of the authorization code
func (t *TokenData) SetCodeExpiresIn(exp time.Duration) {
	t.CodeExpiresIn = exp
}

// GetAccess access Token
func (t *TokenData) GetAccess() string {
	return t.Access
}

// SetAccess access Token
func (t *TokenData) SetAccess(access string) {
	t.Access = access
}

// GetAccessCreateAt create Time
func (t *TokenData) GetAccessCreateAt() time.Time {
	return t.AccessCreateAt
}

// SetAccessCreateAt create Time
func (t *TokenData) SetAccessCreateAt(createAt time.Time) {
	t.AccessCreateAt = createAt
}

// GetAccessExpiresIn the lifetime in seconds of the access token
func (t *TokenData) GetAccessExpiresIn() time.Duration {
	return t.AccessExpiresIn
}

// SetAccessExpiresIn the lifetime in seconds of the access token
func (t *TokenData) SetAccessExpiresIn(exp time.Duration) {
	t.AccessExpiresIn = exp
}

// GetRefresh refresh Token
func (t *TokenData) GetRefresh() string {
	return t.Refresh
}

// SetRefresh refresh Token
func (t *TokenData) SetRefresh(refresh string) {
	t.Refresh = refresh
}

// GetRefreshCreateAt create Time
func (t *TokenData) GetRefreshCreateAt() time.Time {
	return t.RefreshCreateAt
}

// SetRefreshCreateAt create Time
func (t *TokenData) SetRefreshCreateAt(createAt time.Time) {
	t.RefreshCreateAt = createAt
}

// GetRefreshExpiresIn the lifetime in seconds of the refresh token
func (t *TokenData) GetRefreshExpiresIn() time.Duration {
	return t.RefreshExpiresIn
}

// SetRefreshExpiresIn the lifetime in seconds of the refresh token
func (t *TokenData) SetRefreshExpiresIn(exp time.Duration) {
	t.RefreshExpiresIn = exp
}
