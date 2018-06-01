// authors: wangoo
// created: 2018-05-31
// oauth2 session

package main

import (
	"gopkg.in/session.v2"
	"gopkg.in/go-session/redis.v1"
	"time"
)

func initSession() {
	// github.com/go-session/redis
	rdsOpt := redis.Options{
		Addr:     rdsAddr,
		Password: rdsPassword,
	}
	expSeconds := int((time.Minute * 1).Seconds())
	session.InitManager(
		session.SetCookieName("o2s_id"),
		session.SetSign([]byte("sign")),
		session.SetStore(redis.NewRedisStore(&rdsOpt)),
		session.SetCookieLifeTime(expSeconds),
		session.SetExpired(int64(expSeconds)),
	)

	// memory session
	//expSeconds := int((time.Minute * 30).Seconds())
	//session.InitManager(
	//	session.SetCookieName("o2s_id"),
	//	session.SetSign([]byte("sign")),
	//	session.SetCookieLifeTime(expSeconds),
	//	session.SetExpired(int64(expSeconds)),
	//)
}
