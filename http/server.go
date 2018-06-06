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
)

const (
	Oauth2ListenAddr = ":9096"
)

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

	o2.InitOauth2Server(cs, ts, us, as, cfg, http.HandleFunc)

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
	u := &o2x.User{
		UserID:   "u1",
		Nickname: "u1",
	}
	u.SetPassword("123456")
	err := us.Save(u)
	if err != nil {
		log.Printf("create demo user error: %v\n", err)
	}
}