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
	http.HandleFunc(o2.Oauth2UriLogin, o2.LoginHandler)
	http.HandleFunc(o2.Oauth2UriAuth, o2.AuthHandler)
	http.HandleFunc(o2.Oauth2UriAuthorize, o2.AuthorizeRequestHandler)
	http.HandleFunc(o2.Oauth2UriToken, o2.TokenRequestHandler)
	http.HandleFunc(o2.Oauth2UriValid, o2.BearerTokenValidator)

	ts, err := store.NewMemoryTokenStore()
	if err != nil {
		panic(err)
	}
	cs := store.NewClientStore()
	us := o2x.NewUserStore()

	o2.InitOauth2Server(cs, ts, us, nil)

	DemoClient(cs)
	DemoUser(us)

	log.Fatal(http.ListenAndServe(Oauth2ListenAddr, nil))
	log.Println("oauth2 server start on ", Oauth2ListenAddr)
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
