package jwtex

import (
	"gopkg.in/oauth2.v3"
)

// NewJWTAccessGenerate create to generate the jwt access token instance
func NewJWTAccessGenerate(cfg JWTConfig) *JWTAccessGenerate {
	return &JWTAccessGenerate{
		cfg: cfg,
	}
}

// JWTAccessGenerate generate the jwt access token
type JWTAccessGenerate struct {
	cfg JWTConfig
}

// Token based on the UUID generated token
// Registered Claim Names: 	https://tools.ietf.org/html/rfc7519#section-4.1
func (a *JWTAccessGenerate) Token(data *oauth2.GenerateBasic, isGenRefresh bool) (access, refresh string, err error) {
	claims := AccessClaims(data.TokenInfo)
	access, err = GenerateJWT(a.cfg.SigningMethod, a.cfg.SignedKey, claims)
	if err != nil {
		return
	}

	if isGenRefresh {
		claims := RefreshClaims(data.TokenInfo)
		refresh, err = GenerateJWT(a.cfg.SigningMethod, a.cfg.SignedKey, claims)
	}

	return
}
