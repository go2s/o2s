// authors: wangoo
// created: 2018-05-30
// gin oauth2 server using redis store

package main

import (
	"github.com/go2s/o2s/o2"
	"github.com/go2s/o2s/engine"
	"github.com/go2s/o2x"
	"github.com/go2s/o2r"
	"github.com/go-redis/redis"
	"github.com/golang/glog"
	"time"
)

const (
	rdsAddr     = "127.0.0.1:6379"
	rdsPassword = ""

	Oauth2ListenAddr = ":9096"
)

func main() {
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
	as := o2x.NewAuthStore()

	cfg := o2.DefaultServerConfig()
	cfg.ServerName = "Test Gin Oauth2 Server"
	cfg.TemplatePrefix = "../template/"

	svr := o2.InitOauth2Server(cs, ts, us, as, cfg, engine.GinMap)

	mcs, err := o2x.NewMemoryCaptchaStore(time.Minute * 5)
	if err != nil {
		panic(err)
	}
	svr.EnableCaptchaAuth(mcs, o2.CaptchaLogSender)

	engine := engine.GetGinEngine()
	engine.Run(Oauth2ListenAddr)
	glog.Info("oauth2 server start on ", Oauth2ListenAddr)
}
