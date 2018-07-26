// authors: wangoo
// created: 2018-07-26

package o2

import (
	"gopkg.in/oauth2.v3"
	"net/http"
)

type GrantTypeRequestValidator func(r *http.Request) (gt oauth2.GrantType, tgr *oauth2.TokenGenerateRequest, err error)

var (
	customGrantRequestValidatorMap = make(map[oauth2.GrantType]GrantTypeRequestValidator)
)
