// authors: wangoo
// created: 2018-07-26

package o2

import (
	"gopkg.in/oauth2.v3"
	"github.com/go2s/o2x"
)

// ClientScopeHandler check the client allows to use scope
func ClientScopeHandler(clientID, scope string) (allowed bool, err error) {
	if scope == "" {
		allowed = true
		return
	}
	cli, err := oauth2Svr.clientStore.GetByID(clientID)
	if err != nil {
		return
	}
	if client, ok := cli.(o2x.O2ClientInfo); ok {
		allowed = o2x.ScopeArrContains(client.GetScopes(), scope)
		return
	}
	allowed = true
	return
}

func ClientAuthorizedHandler(clientID string, grantType oauth2.GrantType) (allowed bool, err error) {
	cli, err := oauth2Mgr.GetClient(clientID)
	if err != nil {
		return
	}

	if o2ClientInfo, ok := cli.(o2x.O2ClientInfo); ok {
		if o2ClientInfo.GetGrantTypes() != nil {
			for _, t := range o2ClientInfo.GetGrantTypes() {
				if t == grantType {
					return true, nil
				}
			}
		}
	}

	return
}
