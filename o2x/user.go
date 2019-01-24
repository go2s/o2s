// authors: wangoo
// created: 2018-05-31
// oauth2 user extension

package o2x

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"sync"

	"golang.org/x/crypto/scrypt"
)

const (
	pwSaltByteLen = 16
	pwHashByteLen = 32
)

var (
	UserType          = reflect.TypeOf(new(User)).Elem()
	SimpleUserPtrType = reflect.TypeOf(&SimpleUser{})
)

type UserStore interface {
	Save(u User) (err error)
	Remove(id interface{}) (err error)
	UpdatePwd(id interface{}, password string) (err error)
	UpdateScope(id interface{}, clientId, scope string) (err error)
	Find(id interface{}) (u User, err error)
	FindMobile(mobile string) (u User, err error)
}

type User interface {
	GetID() string
	GetUserID() interface{}
	GetMobile() string
	GetPassword() []byte
	GetSalt() []byte
	GetScopes() map[string]string
	SetRawPassword(password string)
	Match(password string) bool
}

type Hexer interface {
	Hex() string
}

func UserIdString(uid interface{}) (id string, err error) {
	if sid, ok := uid.(string); ok {
		id = sid
		return
	}
	if hexer, ok := uid.(Hexer); ok {
		id = hexer.Hex()
		return
	}
	if stringer, ok := uid.(fmt.Stringer); ok {
		id = stringer.String()
		return
	}
	err = errors.New("unknown user id type")
	return
}

func IsUserType(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct && t.Implements(UserType)
}

func NewUser(t reflect.Type) User {
	return reflect.New(t.Elem()).Interface().(User)
}

// -------------------------------
func NewUserStore() UserStore {
	return &MemoryUserStore{
		data:       make(map[interface{}]User),
		mobileData: make(map[string]User),
	}
}

type MemoryUserStore struct {
	sync.RWMutex
	data       map[interface{}]User
	mobileData map[string]User
}

func (cs *MemoryUserStore) Find(id interface{}) (u User, err error) {
	cs.RLock()
	defer cs.RUnlock()
	if c, ok := cs.data[id]; ok {
		u = c
		return
	}
	err = errors.New("not found")
	return
}

func (cs *MemoryUserStore) FindMobile(mobile string) (u User, err error) {
	cs.RLock()
	defer cs.RUnlock()
	if c, ok := cs.mobileData[mobile]; ok {
		u = c
		return
	}
	err = errors.New("not found")
	return
}

func (cs *MemoryUserStore) Save(u User) (err error) {
	cs.Lock()
	defer cs.Unlock()
	cs.data[u.GetUserID()] = u

	if u.GetMobile() != "" {
		cs.mobileData[u.GetMobile()] = u
	}
	return
}

func (cs *MemoryUserStore) Remove(id interface{}) (err error) {
	cs.Lock()
	defer cs.Unlock()
	cs.data[id] = nil
	return
}

func (cs *MemoryUserStore) UpdatePwd(id interface{}, password string) (err error) {
	cs.Lock()
	defer cs.Unlock()
	if c, ok := cs.data[id]; ok {
		c.SetRawPassword(password)
		return
	}
	err = errors.New("not found")
	return
}

func (cs *MemoryUserStore) UpdateScope(id interface{}, clientId, scope string) (err error) {
	cs.Lock()
	defer cs.Unlock()
	if c, ok := cs.data[id]; ok {
		if c.GetScopes() == nil {
			if u, ok := c.(*SimpleUser); ok {
				u.Scopes = make(map[string]string)
			} else {
				err = errors.New("user scope map is nil")
				return
			}
		}
		c.GetScopes()[clientId] = scope
		return
	}
	err = errors.New("not found")
	return
}

// -------------------------------
type SimpleUser struct {
	UserID   interface{}       `bson:"_id" json:"user_id"`
	Mobile   string            `bson:"mobile" json:"mobile"`
	Password []byte            `bson:"password" json:"password"`
	Salt     []byte            `bson:"salt" json:"salt"`
	Scopes   map[string]string `bson:"scopes" json:"scopes,omitempty"`
}

func (u *SimpleUser) GetUserID() interface{} {
	return u.UserID
}

func (u *SimpleUser) GetID() string {
	id, err := UserIdString(u.UserID)
	if err != nil {
		return ""
	}
	return id
}

func (u *SimpleUser) GetMobile() string {
	return u.Mobile
}

func (u *SimpleUser) SetUserID(userID interface{}) {
	u.UserID = userID
}

func (u *SimpleUser) calcHash(password string) (hash []byte, err error) {
	return scrypt.Key([]byte(password), u.Salt, 1<<14, 8, 1, pwHashByteLen)
}

func (u *SimpleUser) SetRawPassword(password string) {
	salt := make([]byte, pwSaltByteLen)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		log.Fatal(err)
	}
	u.Salt = salt

	hash, err := u.calcHash(password)
	if err != nil {
		log.Fatal(err)
	}
	u.Password = hash
}

func (u *SimpleUser) GetPassword() []byte {
	return u.Password
}

func (u *SimpleUser) GetSalt() []byte {
	return u.Salt
}

func (u *SimpleUser) Match(password string) bool {
	hash, err := u.calcHash(password)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return reflect.DeepEqual(hash, u.Password)
}

func (u *SimpleUser) GetScopes() map[string]string {
	return u.Scopes
}
