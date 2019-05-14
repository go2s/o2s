// authors: wangoo
// created: 2018-05-31
// oauth2 server base on gin

package main

import (
	"github.com/go-session/redis"
	"github.com/go-session/session"
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
