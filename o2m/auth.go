// authors: wangoo
// created: 2018-06-05
// mongodb auth store

package o2m

import (
	"context"
	"strings"

	"github.com/go2s/o2s/o2x"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DefaultOauth2AuthDb         = "oauth2"
	DefaultOauth2AuthCollection = "auth"
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
	client     *mongo.Client
	db         string
	collection string
}

func NewAuthStore(client *mongo.Client, db string, collection string) (store *MgoAuthStore) {
	if client == nil {
		panic("client cannot be nil")
	}
	store = &MgoAuthStore{client: client, db: db, collection: collection}
	if store.db == "" {
		store.db = DefaultOauth2AuthDb
	}
	if store.collection == "" {
		store.collection = DefaultOauth2AuthCollection
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

	c := s.client.Database(s.db).Collection(s.collection)
	_, err := c.InsertOne(context.TODO(), mau)
	return err
}

// find auth by clientID and userID
func (s *MgoAuthStore) Find(clientId string, userID string) (auth o2x.Auth, err error) {
	auth = &MgoAuth{}
	filter := bson.M{"_id": buildAuthID(clientId, userID)}
	err = s.client.Database(s.db).Collection(s.collection).FindOne(context.TODO(), filter).Decode(auth)
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
