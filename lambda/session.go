// authors: wangoo
// created: 2018-05-31
// oauth2 session

package main

import (
	"time"
	"gopkg.in/session.v2"
)

func initSession() {
	// github.com/go-session/redis
	//rdsOpt = redis.Options{
	//	Addr:     rdsAddr,
	//	Password: rdsPassword,
	//}
	//session.InitManager(
	//	session.SetCookieName("session_id"),
	//	session.SetSign([]byte("sign")),
	//	session.SetStore(redis.NewRedisStore(&rdsOpt)),
	//)

	// memory session
	expSeconds := int((time.Minute * 30).Seconds())
	session.InitManager(
		session.SetCookieName("session_id"),
		session.SetSign([]byte("sign")),
		session.SetCookieLifeTime(expSeconds),
		session.SetExpired(int64(expSeconds)),
	)
}
