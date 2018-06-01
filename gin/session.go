// authors: wangoo
// created: 2018-05-31
// oauth2 server base on gin

package main

import (
	"gopkg.in/session.v2"
	"gopkg.in/go-session/redis.v1"
)

func init() {
	rdsOpt := redis.Options{
		Addr:     rdsAddr,
		Password: rdsPassword,
	}
	session.InitManager(
		session.SetCookieName("session_id"),
		session.SetSign([]byte("sign")),
		session.SetStore(redis.NewRedisStore(&rdsOpt)),
	)
}
