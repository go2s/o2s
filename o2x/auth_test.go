// authors: wangoo
// created: 2018-06-05
// auth test

package o2x

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	assert.True(t, ScopeContains("s", ""))
	assert.True(t, ScopeContains("", ""))

	as := NewAuthStore()

	a := &AuthModel{
		ClientID: "c1",
		UserID:   "u1",
		Scope:    "s1,s2,s3",
	}

	assert.True(t, a.Contains("s1"))
	assert.True(t, a.Contains("s2"))
	assert.True(t, a.Contains("s3"))
	assert.True(t, a.Contains("s1,s2"))
	assert.True(t, a.Contains("s2,s3"))
	assert.True(t, a.Contains("s1,s3"))

	assert.False(t, a.Contains("s4"))
	assert.False(t, a.Contains("s1,s4"))

	as.Save(a)

	a2, err := as.Find("c1", "u1")
	assert.Nil(t, err)
	assert.Equal(t, "s1,s2,s3", a2.GetScope())
	assert.True(t, as.Exist(a))

	a.SetScope("s1,s2")
	assert.True(t, as.Exist(a))

	a.SetScope("s1,s3")
	assert.True(t, as.Exist(a))

	a.SetScope("s1,s4")
	assert.False(t, as.Exist(a))

}
