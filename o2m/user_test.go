// authors: wangoo
// created: 2018-06-28
// test user

package o2m

import (
	"testing"

	"github.com/go2s/o2s/o2x"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	mgoDatabase  = "oauth2"
	mgoUsername  = "oauth2"
	mgoPassword  = "oauth2"
	mgoPoolLimit = 10
)

var mgoAddress = "mongodb://127.0.0.1:27017"

func TestMgoUserStore(t *testing.T) {
	mgoCfg := MongoConfig{
		Hosts:     mgoAddress,
		Database:  mgoDatabase,
		Username:  mgoUsername,
		Password:  mgoPassword,
		PoolLimit: mgoPoolLimit,
	}

	mgoClient := NewMongoClient(&mgoCfg)

	cfg := DefaultMgoUserCfg()

	us := NewUserStore(mgoClient, mgoDatabase, "user", cfg)

	id := "5c4aa10d97deb33b2a17633a"
	pass := "123456"
	mobile1 := "13344556677"
	mobile2 := "13344556688"

	us.Remove(id)
	user, err := us.Find(id)
	assert.True(t, err != nil)
	assert.Equal(t, "not found", err.Error())

	if user == nil {
		objectID, _ := primitive.ObjectIDFromHex(id)
		user = &o2x.SimpleUser{
			UserID: objectID,
			Mobile: mobile1,
			Scopes: make(map[string]string),
		}
		user.GetScopes()["c1"] = "read"
		err = us.Save(user)
		if err != nil {
			assert.Fail(t, err.Error())
		}
	}

	user, err = us.Find(id)
	if err != nil {
		t.Error(err)
		us.Remove(id)
		t.FailNow()
	}
	assert.Equal(t, "read", user.GetScopes()["c1"])

	//-------------------------------add user with duplicated mobile
	us.Remove("user2")
	user2 := &o2x.SimpleUser{
		UserID: "user2",
		Mobile: mobile1,
	}
	err = us.Save(user2)
	if err == nil {
		assert.Fail(t, "should throw mobile duplicated error")
	}
	//-------------------------------add user with different mobile
	us.Remove("user3")
	user3 := &o2x.SimpleUser{
		UserID: "user3",
		Mobile: mobile2,
	}
	err = us.Save(user3)
	assert.Nil(t, err)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	//-------------------------------

	us.UpdatePwd(id, pass)

	updateUser, _ := us.Find(id)

	assert.True(t, updateUser.Match(pass))
	assert.False(t, updateUser.Match("password"))

	err = us.UpdateScope(id, "c1", "manage,admin")
	err = us.UpdateScope(id, "c2", "operate,view")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	user, err = us.Find(id)
	assert.Equal(t, "manage,admin", user.GetScopes()["c1"])
	assert.Equal(t, "operate,view", user.GetScopes()["c2"])

	err = us.Remove(id)
	if err != nil {
		assert.Fail(t, err.Error())
	}
}
