// authors: wangoo
// created: 2018-07-16

package o2

import (
	"github.com/go2s/o2x"
	"net/http"
)

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
