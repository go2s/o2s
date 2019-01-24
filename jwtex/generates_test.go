package jwtex_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go2s/o2s/jwtex"
	"github.com/go2s/o2s/util/timeutil"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

func TestJWTAccess(t *testing.T) {
	Convey("Test JWT Access Generate", t, func() {
		fmt.Println()

		userID := "000000"
		clientID := "123456"
		data := &oauth2.GenerateBasic{
			Client: &models.Client{
				ID:     clientID,
				Secret: "123456",
			},
			UserID: userID,
			TokenInfo: &models.Token{
				ClientID:        clientID,
				UserID:          userID,
				AccessCreateAt:  timeutil.Now(),
				AccessExpiresIn: time.Second * 120,
			},
		}

		signedKey := []byte("00000000")
		method := jwt.SigningMethodHS512
		cfg := jwtex.JWTConfig{
			SigningMethod: method,
			SignedKey:     signedKey,
		}
		gen := jwtex.NewJWTAccessGenerate(cfg)

		access, refresh, err := gen.Token(data, true)
		So(err, ShouldBeNil)
		So(access, ShouldNotBeEmpty)
		So(refresh, ShouldNotBeEmpty)
		fmt.Printf("access:%s\n", access)
		fmt.Printf("refresh:%s\n", refresh)

		claims, err := jwtex.ParseClaims(method, signedKey, access)

		So(err, ShouldBeNil)

		So(claims.Valid(), ShouldBeNil)
		So(claims.Issuer, ShouldEqual, clientID)
		So(claims.Subject, ShouldEqual, userID)
	})
}
