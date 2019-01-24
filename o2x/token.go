// authors: wangoo
// created: 2018-06-01
// oauth2 token extension

package o2x

import "gopkg.in/oauth2.v3"

type O2TokenStore interface {
	oauth2.TokenStore

	RemoveByAccount(userID string, clientID string) (err error)
	GetByAccount(userID string, clientID string) (ti oauth2.TokenInfo, err error)
}
