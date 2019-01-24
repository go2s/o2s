// authors: wangoo
// created: 2018-05-30
// test user

package o2x

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	u := &SimpleUser{}
	u.SetUserID("123")

	assert.Equal(t, "123", u.GetUserID())

	password := "my_password"
	u.SetRawPassword(password)

	assert.True(t, u.Match(password), "password should be match")

	js, err := json.Marshal(u)
	assert.Nil(t, err, err)
	fmt.Println(string(js))
}

func TestNewUser(t *testing.T) {
	u := NewUser(SimpleUserPtrType)
	fmt.Println("user:", u)

	u2 := u.(*SimpleUser)
	u2.SetUserID("u2")
	u2.SetRawPassword("pass")

	js, err := json.Marshal(u)
	assert.Nil(t, err, err)
	fmt.Println(string(js))
}

func TestIsUserType(t *testing.T) {
	fmt.Println(SimpleUserPtrType)
	fmt.Println(UserType)
	assert.True(t, IsUserType(SimpleUserPtrType))
}

// http://www.mongodb.org/display/DOCS/Object+IDs
type ObjectId string

// ObjectIdHex returns an ObjectId from the provided hex representation.
// Calling this function with an invalid hex representation will
// cause a runtime panic. See the IsObjectIdHex function.
func ObjectIdHex(s string) ObjectId {
	d, err := hex.DecodeString(s)
	if err != nil || len(d) != 12 {
		panic(fmt.Sprintf("invalid input to ObjectIdHex: %q", s))
	}
	return ObjectId(d)
}

// Hex returns a hex representation of the ObjectId.
func (id ObjectId) Hex() string {
	return hex.EncodeToString([]byte(id))
}

func TestUserIdString(t *testing.T) {
	id := "5ae6b2005946fa106132365c"

	sid, err := UserIdString(id)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, id, sid)

	hexer := ObjectIdHex(id)
	hid, err := UserIdString(hexer)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, id, hid)

}
