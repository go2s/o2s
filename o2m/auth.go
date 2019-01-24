// authors: wangoo
// created: 2018-06-05
// mongodb auth store

package o2m

import (
	"strings"
	"gopkg.in/mgo.v2"
	"github.com/go2s/o2s/o2x"
)

const (
	DefaultOuath2AuthDb         = "oauth2"
	DefaultOuath2AuthCollection = "auth"
	idSplit                     = "__"
)

// AuthID = ClientId + "__" + UserID
type MgoAuth struct {
	AuthID   string `bson:"_id" json:"-"`
	ClientID string `bson:"-" json:"client_id"`
	UserID   string `bson:"-" json:"user_id"`
	Scope    string `bson:"scope" json:"scope"`
}

func buildAuthID(clientID string, userID string) string {
	return clientID + idSplit + userID
}

func (au *MgoAuth) UpdateAuthID() {
	au.AuthID = buildAuthID(au.GetClientID(), au.GetUserID())
}

func (au *MgoAuth) GetClientID() string {
	if au.ClientID == "" {
		idx := strings.Index(au.AuthID, idSplit)
		if idx > 0 {
			au.ClientID = au.AuthID[:idx]
		}
	}
	return au.ClientID
}

func (au *MgoAuth) SetClientID(id string) {
	au.ClientID = id
	au.UpdateAuthID()
}

func (au *MgoAuth) GetUserID() string {
	if au.UserID == "" {
		idx := strings.Index(au.AuthID, idSplit)
		if idx > 0 {
			au.UserID = au.AuthID[idx+len(idSplit):]
		}
	}
	return au.UserID
}

func (au *MgoAuth) SetUserID(id string) {
	au.UserID = id
	au.UpdateAuthID()
}

func (au *MgoAuth) GetScope() string {
	return au.Scope
}

func (au *MgoAuth) SetScope(scope string) {
	au.Scope = scope
}

func (au *MgoAuth) Contains(scope string) bool {
	return o2x.ScopeContains(au.Scope, scope)
}

type MgoAuthStore struct {
	session    *mgo.Session
	db         string
	collection string
}

func NewAuthStore(session *mgo.Session, db string, collection string) (store *MgoAuthStore) {
	if session == nil {
		panic("session cannot be nil")
	}
	store = &MgoAuthStore{session: session, db: db, collection: collection}
	if store.db == "" {
		store.db = DefaultOuath2AuthDb
	}
	if store.collection == "" {
		store.collection = DefaultOuath2AuthCollection
	}

	return
}

func (s *MgoAuthStore) Save(auth o2x.Auth) error {
	mau := &MgoAuth{
		ClientID: auth.GetClientID(),
		UserID:   auth.GetUserID(),
		Scope:    auth.GetScope(),
	}
	mau.UpdateAuthID()

	session := s.session.Clone()
	defer session.Close()

	c := session.DB(s.db).C(s.collection)
	return c.Insert(mau)
}

// find auth by clientID and userID
func (s *MgoAuthStore) Find(clientId string, userID string) (auth o2x.Auth, err error) {
	session := s.session.Clone()
	defer session.Close()

	auth = &MgoAuth{}
	err = session.DB(s.db).C(s.collection).FindId(buildAuthID(clientId, userID)).One(auth)
	if err != nil {
		return nil, err
	}
	auth.GetUserID()
	auth.GetClientID()
	return
}

// whether the auth already exists
func (s *MgoAuthStore) Exist(auth o2x.Auth) bool {
	au, err := s.Find(auth.GetClientID(), auth.GetUserID())
	if err != nil {
		return false
	}
	return au.Contains(auth.GetScope())
}
