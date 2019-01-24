// authors: wangoo
// created: 2018-05-29
// oauth2 mongodb token storage

package o2m

import (
	"gopkg.in/oauth2.v3"
	"gopkg.in/mgo.v2"
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/go2s/o2s/o2x"
)

// MgoTokenStore MongoDB storage for OAuth 2.0
type MgoTokenStore struct {
	db         string
	collection string
	session    *mgo.Session
}

// NewTokenStore create a token store instance based on mongodb
func NewTokenStore(session *mgo.Session, db string,
	collection string) (store oauth2.TokenStore) {
	ts := &MgoTokenStore{
		session:    session,
		db:         db,
		collection: collection,
	}

	err := ts.c(ts.collection).EnsureIndex(mgo.Index{
		Key:         []string{"ExpiredAt"},
		ExpireAfter: time.Second * 1,
	})
	if err != nil {
		panic(err)
	}

	err = ts.c(ts.collection).EnsureIndex(mgo.Index{
		Key: []string{"UserID"},
	})
	if err != nil {
		panic(err)
	}

	err = ts.c(ts.collection).EnsureIndex(mgo.Index{
		Key: []string{"ClientId"},
	})
	if err != nil {
		panic(err)
	}

	err = ts.c(ts.collection).EnsureIndex(mgo.Index{
		Key: []string{"Refresh"},
	})
	if err != nil {
		panic(err)
	}
	store = ts
	return
}

func (ts *MgoTokenStore) c(name string) *mgo.Collection {
	return ts.session.DB(ts.db).C(name)
}

func (ts *MgoTokenStore) H(name string, handler func(c *mgo.Collection)) {
	session := ts.session.Clone()
	defer session.Close()
	handler(session.DB(ts.db).C(name))
	return
}

// Create create and store the new token information
func (ts *MgoTokenStore) Create(info oauth2.TokenInfo) (err error) {
	token := Copy(info)
	ts.H(ts.collection, func(c *mgo.Collection) {
		err = c.Insert(token)
	})
	return
}

// RemoveByCode use the authorization code to delete the token information
func (ts *MgoTokenStore) RemoveByCode(code string) (err error) {
	ts.H(ts.collection, func(c *mgo.Collection) {
		mgoErr := c.Remove(bson.M{"Code": code})
		if mgoErr != nil {
			if mgoErr == mgo.ErrNotFound {
				err = o2x.ErrNotFound
				return
			}
			err = mgoErr
		}
	})
	return
}

// RemoveByAccess use the access token to delete the token information
func (ts *MgoTokenStore) RemoveByAccess(access string) (err error) {
	ts.H(ts.collection, func(c *mgo.Collection) {
		mgoErr := c.RemoveId(access)
		if mgoErr != nil {
			if mgoErr == mgo.ErrNotFound {
				err = o2x.ErrNotFound
				return
			}
			err = mgoErr
		}
	})
	return
}

// RemoveByRefresh use the refresh token to delete the token information
func (ts *MgoTokenStore) RemoveByRefresh(refresh string) (err error) {
	ts.H(ts.collection, func(c *mgo.Collection) {
		mgoErr := c.Remove(bson.M{"Refresh": refresh})
		if mgoErr != nil {
			if mgoErr == mgo.ErrNotFound {
				err = o2x.ErrNotFound
				return
			}
			err = mgoErr
		}
	})
	return
}

// RemoveByAccount remove exists token info by userID and clientID
func (ts *MgoTokenStore) RemoveByAccount(userID string, clientID string) (err error) {
	ts.H(ts.collection, func(c *mgo.Collection) {
		err = c.Remove(bson.M{"UserID": userID, "ClientId": clientID})
		if err == nil && bson.IsObjectIdHex(userID) {
			err = c.Remove(bson.M{"UserID": bson.ObjectIdHex(userID), "ClientId": clientID})
		}
	})
	return
}

// GetByField use field value for token information data
func (ts *MgoTokenStore) GetByBson(m bson.M) (ti oauth2.TokenInfo, err error) {
	ts.H(ts.collection, func(c *mgo.Collection) {
		token := &TokenData{}
		mgoErr := c.Find(m).One(token)
		if mgoErr != nil {
			if mgoErr == mgo.ErrNotFound {
				err = o2x.ErrNotFound
				return
			}
			err = mgoErr
			return
		}
		ti = token
	})
	return
}

// GetByField use field value for token information data
func (ts *MgoTokenStore) GetByField(field string, value string) (ti oauth2.TokenInfo, err error) {
	ti, err = ts.GetByBson(bson.M{field: value})
	return
}

// GetByCode use the authorization code for token information data
func (ts *MgoTokenStore) GetByCode(code string) (ti oauth2.TokenInfo, err error) {
	ti, err = ts.GetByField("Code", code)
	return
}

// GetByAccess use the access token for token information data
func (ts *MgoTokenStore) GetByAccess(access string) (ti oauth2.TokenInfo, err error) {
	ti, err = ts.GetByField("_id", access)
	return
}

// GetByRefresh use the refresh token for token information data
func (ts *MgoTokenStore) GetByRefresh(refresh string) (ti oauth2.TokenInfo, err error) {
	ti, err = ts.GetByField("Refresh", refresh)
	return
}

// GetByAccount get the exists token info by userID and clientID
func (ts *MgoTokenStore) GetByAccount(userID string, clientID string) (ti oauth2.TokenInfo, err error) {
	ti, err = ts.GetByBson(bson.M{"UserID": userID, "ClientId": clientID})
	if err == nil && ti == nil && bson.IsObjectIdHex(userID) {
		ti, err = ts.GetByBson(bson.M{"UserID": bson.ObjectIdHex(userID), "ClientId": clientID})
	}
	return
}
