// authors: wangoo
// created: 2018-05-30
// oauth2 http server using memory store

package main

import (
	"net/http"
	"log"
	"github.com/go2s/o2s/o2"
	"gopkg.in/oauth2.v3/store"
	"github.com/go2s/o2x"
	"github.com/go2s/o2m"
	"github.com/golang/glog"
	"time"
	"flag"
)

const (
	Oauth2ListenAddr = ":9096"
)

var (
	handleMap = map[string]func(w http.ResponseWriter, r *http.Request){}
)

func HandleHttp(method, pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
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

	ts, err := store.NewMemoryTokenStore()
	if err != nil {
		panic(err)
	}
	cs := store.NewClientStore()
	us := o2x.NewUserStore()
	as := o2x.NewAuthStore()

	cfg := o2.DefaultServerConfig()
	cfg.ServerName = "Test Memory Oauth2 Server"
	cfg.TemplatePrefix = "../template/"

	svr := o2.InitOauth2Server(cs, ts, us, as, cfg, HandleHttp)

	mcs, err := o2x.NewMemoryCaptchaStore(time.Minute * 5)
	if err != nil {
		panic(err)
	}
	svr.EnableCaptchaAuth(mcs, o2.CaptchaLogSender)

	DemoClient(cs)
	DemoUser(us)

	glog.Info("oauth2 server start on ", Oauth2ListenAddr)
	glog.Fatal(http.ListenAndServe(Oauth2ListenAddr, nil))
}

func DemoClient(cs o2x.Oauth2ClientStore) {
	err := cs.Set("000000", &o2m.Oauth2Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "https://localhost",
	})
	if err != nil {
		log.Printf("%v\n", err)
	}
}

func DemoUser(us o2x.UserStore) {
	u := &o2x.SimpleUser{
		UserID: "u1",
	}
	u.SetRawPassword("123456")
	u.Mobile = "13344556677"
	err := us.Save(u)
	if err != nil {
		glog.Infof("create demo user error: %v", err)
	}
}
