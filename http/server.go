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
	"fmt"
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
	fmt.Printf("http map %v\n", pattern)
	handleMap[pattern] = handler
	http.DefaultServeMux.HandleFunc(pattern, handler)
}

func main() {
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

	o2.InitOauth2Server(cs, ts, us, as, cfg, HandleHttp)

	DemoClient(cs)
	DemoUser(us)

	log.Println("oauth2 server start on ", Oauth2ListenAddr)
	log.Fatal(http.ListenAndServe(Oauth2ListenAddr, nil))
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
	err := us.Save(u)
	if err != nil {
		log.Printf("create demo user error: %v\n", err)
	}
}
