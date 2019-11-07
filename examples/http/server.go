// authors: wangoo
// created: 2018-05-30
// oauth2 http server using memory store

package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go2s/o2s/captcha"
	"github.com/go2s/o2s/jwtex"
	"github.com/go2s/o2s/o2"
	"github.com/go2s/o2s/o2m"
	"github.com/go2s/o2s/o2x"
	"github.com/golang/glog"
	"gopkg.in/oauth2.v3/store"
)

const (
	//Oauth2ListenAddr listen address
	Oauth2ListenAddr = ":9096"
)

var (
	handleMap = map[string]func(w http.ResponseWriter, r *http.Request){}
)

//HandleHTTP bind http handler
func HandleHTTP(method, pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	if _, exist := handleMap[pattern]; exist {
		return
	}
	glog.Infof("http map %v", pattern)
	handleMap[pattern] = handler
	http.DefaultServeMux.HandleFunc(pattern, handler)
}

func main() {
	flag.Parse()
	flag.Set("logtostderr", "true") // Log to stderr only, instead of file.

	cfg := o2.DefaultServerConfig()
	cfg.ServerName = "Test Memory Oauth2 Server"
	cfg.JWTSupport = true
	cfg.JWT = jwtex.JWTConfig{
		SignedKey:     []byte("go2s"),
		SigningMethod: jwt.SigningMethodHS512,
	}
	ts, err := store.NewMemoryTokenStore()
	if err != nil {
		panic(err)
	}

	cs := store.NewClientStore()
	us := o2x.NewUserStore()
	as := o2x.NewAuthStore()

	svr := o2.InitOauth2Server(cs, ts, us, as, cfg, HandleHTTP)

	mcs, err := o2x.NewMemoryCaptchaStore(time.Minute * 5)
	if err != nil {
		panic(err)
	}
	captcha.EnableCaptchaAuth(svr, mcs, captcha.CaptchaLogSender)

	DemoClient(cs)
	DemoUser(us)

	glog.Info("oauth2 server start on ", Oauth2ListenAddr)
	glog.Fatal(http.ListenAndServe(Oauth2ListenAddr, nil))
}

//DemoClient init demo client
func DemoClient(cs o2x.O2ClientStore) {
	err := cs.Set("000000", &o2m.Oauth2Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "https://localhost",
		Scopes: []string{"manage", "admin", "view", "read"},
	})
	if err != nil {
		log.Printf("%v\n", err)
	}
}

//DemoUser init demo user
func DemoUser(us o2x.UserStore) {
	u := &o2x.SimpleUser{
		UserID: "u1",
		Scopes: map[string]string{
			"000000": "admin,manage",
		},
	}
	u.SetRawPassword("123456")
	u.Mobile = "13344556677"
	err := us.Save(u)
	if err != nil {
		glog.Infof("create demo user error: %v", err)
	}
}
