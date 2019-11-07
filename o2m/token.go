// authors: wangoo
// created: 2018-05-29
// oauth2 mongodb token storage

package o2m

import (
	"context"

	"github.com/go2s/o2s/o2x"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"gopkg.in/oauth2.v3"
)

// MgoTokenStore MongoDB storage for OAuth 2.0
type MgoTokenStore struct {
	db         string
	collection string
	client     *mongo.Client
}

// NewTokenStore create a token store instance based on mongodb
func NewTokenStore(client *mongo.Client, db string,
	collection string) (store oauth2.TokenStore) {
	ts := &MgoTokenStore{
		client:     client,
		db:         db,
		collection: collection,
	}

	option := options.Index().SetExpireAfterSeconds(1)
	_, err := ts.c(ts.collection).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bsonx.Doc{bsonx.Elem{Key: "ExpiredAt", Value: bsonx.Int32(1)}},
		Options: option,
	})

	if err != nil {
		panic(err)
	}

	_, err = ts.c(ts.collection).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bsonx.Doc{bsonx.Elem{Key: "UserID", Value: bsonx.Int32(1)}},
	})
	if err != nil {
		panic(err)
	}

	_, err = ts.c(ts.collection).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bsonx.Doc{bsonx.Elem{Key: "ClientId", Value: bsonx.Int32(1)}},
	})
	if err != nil {
		panic(err)
	}

	_, err = ts.c(ts.collection).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bsonx.Doc{bsonx.Elem{Key: "Refresh", Value: bsonx.Int32(1)}},
	})
	if err != nil {
		panic(err)
	}
	store = ts
	return
}

func (ts *MgoTokenStore) c(name string) *mongo.Collection {
	return ts.client.Database(ts.db).Collection(name)
}

func (ts *MgoTokenStore) H(name string, handler func(c *mongo.Collection)) {
	handler(ts.client.Database(ts.db).Collection(name))
}

// Create create and store the new token information
func (ts *MgoTokenStore) Create(info oauth2.TokenInfo) (err error) {
	token := Copy(info)
	ts.H(ts.collection, func(c *mongo.Collection) {
		_, err = c.InsertOne(context.TODO(), token)
	})
	return
}

// RemoveByCode use the authorization code to delete the token information
func (ts *MgoTokenStore) RemoveByCode(code string) (err error) {
	ts.H(ts.collection, func(c *mongo.Collection) {
		_, mgoErr := c.DeleteMany(context.TODO(), bson.M{"Code": code})
		if mgoErr != nil {
			if mgoErr == mongo.ErrNoDocuments {
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
	ts.H(ts.collection, func(c *mongo.Collection) {
		_, mgoErr := c.DeleteOne(context.TODO(), bson.M{"_id": access})
		if mgoErr != nil {
			if mgoErr == mongo.ErrNoDocuments {
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
	ts.H(ts.collection, func(c *mongo.Collection) {
		_, mgoErr := c.DeleteMany(nil, bson.M{"Refresh": refresh})
		if mgoErr != nil {
			if mgoErr == mongo.ErrNoDocuments {
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
	ts.H(ts.collection, func(c *mongo.Collection) {
		res, err := c.DeleteMany(nil, bson.M{"UserID": userID, "ClientId": clientID})
		if err == nil && res.DeletedCount == 0 {
			objectID, err := primitive.ObjectIDFromHex(userID)
			if err == nil {
				_, err = c.DeleteMany(nil, bson.M{"UserID": objectID, "ClientId": clientID})
			}
		}
	})
	return
}

// GetByField use field value for token information data
func (ts *MgoTokenStore) GetByBson(m bson.M) (ti oauth2.TokenInfo, err error) {
	ts.H(ts.collection, func(c *mongo.Collection) {
		token := &TokenData{}
		mgoErr := c.FindOne(nil, m).Decode(token)
		if mgoErr != nil {
			if mgoErr == mongo.ErrNoDocuments {
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
	if err == nil && ti == nil {
		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			ti, err = ts.GetByBson(bson.M{"UserID": objectID, "ClientId": clientID})
		}
	}
	return
}
