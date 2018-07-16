// authors: wangoo
// created: 2018-07-16

package o2

import (
	"github.com/go2s/o2x"
	"net/http"
)

// ClientScopeHandler check the client allows to use scope
func ClientScopeHandler(clientID, scope string) (allowed bool, err error) {
	cli, err := oauth2ClientStore.GetByID(clientID)
	if err != nil {
		return
	}
	if client, ok := cli.(o2x.Oauth2ClientInfo); ok {
		allowed = o2x.ScopeContains(client.GetScope(), scope)
		return
	}
	allowed = true
	return
}

// RefreshingScopeHandler check the scope of the refreshing token
func RefreshingScopeHandler(newScope, oldScope string) (allowed bool, err error) {
	allowed = o2x.ScopeContains(oldScope, newScope)
	return
}

// AuthorizeScopeHandler set the authorized scope
func AuthorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scp string, err error) {
	reqScope := scope(r)
	if reqScope == "" {
		return
	}

	clientID := clientID(r)
	if clientID == "" {
		return
	}

	allowed, err := ClientScopeHandler(clientID, reqScope)
	if err != nil {
		return
	}
	if allowed {
		scp = reqScope
		return
	}
	return
}
