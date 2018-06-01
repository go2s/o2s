// authors: wangoo
// created: 2018-05-30
// gin oauth2 server using redis store

package main

import (
	"github.com/go2s/o2s/o2"
	"log"
	"github.com/go2s/o2s/engine"
	"github.com/go2s/o2x"
	"github.com/go2s/o2r"
	"github.com/go-redis/redis"
)

const (
	rdsAddr     = "127.0.0.1:6379"
	rdsPassword = ""

	Oauth2ListenAddr = ":9096"
)

func main() {
	engine := engine.NewEngine()

	rdsOpt := redis.Options{
		Addr:     rdsAddr,
		Password: rdsPassword,
	}

	ts, err := o2r.NewTokenStore(&rdsOpt)
	if err != nil {
		panic(err)
	}
	cs, err := o2r.NewClientStore(&rdsOpt)
	if err != nil {
		panic(err)
	}
	us := o2x.NewUserStore()

	o2.InitOauth2Server(cs, ts, us, nil)

	engine.Run(Oauth2ListenAddr)
	log.Println("oauth2 server start on ", Oauth2ListenAddr)
}
