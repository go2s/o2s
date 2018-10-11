package o2

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
)

//ParseJWTAccessToken valid jwt access token
func (s *Oauth2Server) ParseJWTAccessToken(accessToken string) (claims *generates.JWTAccessClaims, err error) {
	if !s.cfg.JWT.Support {
		return nil, errors.ErrInvalidAccessToken
	}
	token, err := jwt.ParseWithClaims(accessToken, &generates.JWTAccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != s.cfg.JWT.SignMethod {
			return nil, fmt.Errorf("unknown jwt token")
		}
		return s.cfg.JWT.SignKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*generates.JWTAccessClaims)
	if !ok {
		return nil, fmt.Errorf("not jwt access claims")
	}
	err = claims.Valid()
	if err != nil {
		return nil, err
	}
	return
}
