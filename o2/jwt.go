package o2

import (
	"github.com/go2s/oauth2/errors"
	"github.com/go2s/oauth2/jwtex"
)

//ParseJWTAccessToken valid jwt access token
func (s *Oauth2Server) ParseJWTAccessToken(access string) (claims *jwtex.Oauth2Claims, err error) {
	if !s.cfg.JWTSupport {
		return nil, errors.ErrInvalidAccessToken
	}

	return jwtex.ParseClaims(s.cfg.JWT.SigningMethod, s.cfg.JWT.SignedKey, access)
}
