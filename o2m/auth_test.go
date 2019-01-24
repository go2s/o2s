// authors: wangoo
// created: 2018-06-05
// auth test

package o2m

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMgoAuth(t *testing.T) {
	a := &MgoAuth{}
	a.UserID = "u1"
	a.ClientID = "c1"
	a.Scope = "s1"

	a.UpdateAuthID()

	assert.Equal(t, "c1__u1", a.AuthID)
	assert.True(t, a.Contains("s1"))
	assert.False(t, a.Contains("s2"))

	a = &MgoAuth{
		AuthID: "c2__u2",
		Scope:  "s2,s3,s4",
	}

	assert.Equal(t, "c2", a.GetClientID())
	assert.Equal(t, "u2", a.GetUserID())

	assert.True(t, a.Contains("s2"))
	assert.True(t, a.Contains("s3"))
	assert.True(t, a.Contains("s4"))
	assert.True(t, a.Contains("s2,s3"))
	assert.True(t, a.Contains("s2,s4"))

	assert.False(t, a.Contains("s1"))

}
