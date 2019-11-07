package o2

import (
	"github.com/go2s/o2s/jwtex"
	"gopkg.in/oauth2.v3/errors"
)

//ParseJWTAccessToken valid jwt access token
func (s *Oauth2Server) ParseJWTAccessToken(access string) (claims *jwtex.Oauth2Claims, err error) {
	if !s.cfg.JWTSupport {
		return nil, errors.ErrInvalidAccessToken
	}

	return jwtex.ParseClaims(s.cfg.JWT.SigningMethod, s.cfg.JWT.SignedKey, access)
}
