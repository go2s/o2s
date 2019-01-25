// authors: wangoo
// created: 2018-05-30
// oauth2 user store

package o2m

import (
	"context"
	"fmt"
	"github.com/go2s/o2s/o2x"
	"github.com/golang/glog"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/patrickmn/go-cache"
	"reflect"
	"time"
)

var (
	userCache *cache.Cache
)

func init() {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	userCache = cache.New(5*time.Minute, 10*time.Minute)
}

func addUserCache(user o2x.User) {
	if user.GetUserID() != nil {
		userCache.Add(fmt.Sprint(user.GetUserID()), user, cache.DefaultExpiration)
	}
}

func getUserCache(id interface{}) (user o2x.User) {
	if c, found := userCache.Get(fmt.Sprint(id)); found {
		user = c.(o2x.User)
		return
	}
	return
}

func removeUserCache(id interface{}) {
	userCache.Delete(fmt.Sprint(id))
}

type MgoUserCfg struct {
	userType reflect.Type

	// password field name
	passwordName string

	// salt field name
	saltName string
}

// used to control the unique mobile for one user if exists
type MgoUserMobile struct {
	Id     string `bson:"_id" json:"_id"`
	Mobile string `bson:"mobile" json:"mobile"`
}

type MgoUserStore struct {
	client           *mongo.Client
	db               string
	collection       string
	mobileCollection string
	userCfg          *MgoUserCfg
}

func DefaultMgoUserCfg() *MgoUserCfg {
	return &MgoUserCfg{
		userType:     o2x.SimpleUserPtrType,
		passwordName: "password",
		saltName:     "salt",
	}
}

func NewUserStore(client *mongo.Client, db, collection string, userCfg *MgoUserCfg) (us *MgoUserStore) {
	if !o2x.IsUserType(userCfg.userType) {
		panic("invalid user type")
	}
	us = &MgoUserStore{
		client:           client,
		db:               db,
		collection:       collection,
		mobileCollection: collection + "_mobile",
		userCfg:          userCfg,
	}

	option := options.Index().SetUnique(true)
	_, err := us.c(us.collection).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bsonx.Doc{bsonx.Elem{Key: "mobile", Value: bsonx.Int32(1)}},
		Options: option,
	})
	if err != nil {
		panic(err)
	}

	return
}

func (us *MgoUserStore) c(name string) *mongo.Collection {
	return us.client.Database(us.db).Collection(name)
}

func (us *MgoUserStore) lockUserMobile(userId, mobile string) (err error) {
	if userId == "" || mobile == "" {
		err = o2x.ErrValueRequired
		return
	}
	c := us.c(us.mobileCollection)
	userMobile := &MgoUserMobile{
		Id:     userId,
		Mobile: mobile,
	}
	count, err := c.Count(context.TODO(), userMobile)
	if count > 0 || err != nil {
		return
	}

	_, err = c.InsertOne(context.TODO(), userMobile)

	return
}

func (us *MgoUserStore) unlockUserMobile(userId string) (err error) {
	if userId == "" {
		err = o2x.ErrValueRequired
		return
	}
	c := us.c(us.mobileCollection)
	_, mgoErr := c.DeleteOne(context.TODO(), bson.M{"_id": userId})
	if mgoErr != mongo.ErrNoDocuments {
		err = mgoErr
	}
	return
}

func (us *MgoUserStore) Save(u o2x.User) (err error) {
	if u.GetMobile() != "" {
		err = us.lockUserMobile(u.GetID(), u.GetMobile())
		if err != nil {
			return
		}
	}
	c := us.c(us.collection)
	glog.Infof("insert user:%v", u)
	_, err = c.InsertOne(nil, u)
	if err != nil {
		return
	}
	addUserCache(u)

	return
}

func (us *MgoUserStore) Remove(id interface{}) (err error) {
	removeUserCache(id)

	sid, err := o2x.UserIdString(id)
	if err != nil {
		return
	}

	glog.Infof("remove user:%v", id)
	us.unlockUserMobile(sid)

	c := us.c(us.collection)

	res, mgoErr := c.DeleteOne(context.TODO(), bson.M{"_id": id})

	if mgoErr == nil && res.DeletedCount == 0 {
		// try to find using object id
		if sid, ok := id.(string); ok {
			objectID, err := primitive.ObjectIDFromHex(sid)
			if err == nil {
				_, mgoErr = c.DeleteOne(nil, bson.M{"_id": objectID})
			}
		}
	}

	if mgoErr != nil && mgoErr == mongo.ErrNoDocuments {
		err = o2x.ErrNotFound
		return
	}
	err = mgoErr
	return
}

func (us *MgoUserStore) Find(id interface{}) (u o2x.User, err error) {
	if u = getUserCache(id); u != nil {
		return
	}

	c := us.c(us.collection)

	user := o2x.NewUser(us.userCfg.userType)
	mgoErr := c.FindOne(context.TODO(), bson.M{"_id": id}).Decode(user)
	if mgoErr != nil && mgoErr == mongo.ErrNoDocuments {
		// try to find using object id
		if sid, ok := id.(string); ok {
			objectID, err := primitive.ObjectIDFromHex(sid)
			if err == nil {
				mgoErr = c.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(user)
			}
		}
	}

	if mgoErr != nil {
		if mgoErr == mongo.ErrNoDocuments {
			err = o2x.ErrNotFound
			return
		}
		err = mgoErr
	}

	if err != nil {
		return
	}

	u = user

	if u != nil {
		addUserCache(u)
	}

	return
}

func (us *MgoUserStore) FindMobile(mobile string) (u o2x.User, err error) {

	c := us.c(us.collection)

	user := o2x.NewUser(us.userCfg.userType)
	mgoErr := c.FindOne(context.TODO(), bson.M{"mobile": mobile}).Decode(user)
	if mgoErr != nil && mgoErr == mongo.ErrNoDocuments {
		err = o2x.ErrNotFound
		return
	}
	u = user
	return
}

func (us *MgoUserStore) UpdatePwd(id interface{}, password string) (err error) {
	user, err := us.Find(id)
	if err != nil {
		return
	}
	glog.Infof("update user password %v", id)
	user.SetRawPassword(password)

	c := us.c(us.collection)

	bs := bson.M{us.userCfg.passwordName: user.GetPassword(), us.userCfg.saltName: user.GetSalt()}
	bs = bson.M{"$set": bs}
	_, err = c.UpdateOne(context.TODO(), bson.M{"_id": user.GetUserID()}, bs)
	if err != nil {
		return
	}
	addUserCache(user)
	return
}

func (us *MgoUserStore) UpdateScope(id interface{}, clientId, scope string) (err error) {
	user, err := us.Find(id)
	if err != nil {
		return
	}
	glog.Infof("update user %v client %v scope %v", id, clientId, scope)

	c := us.c(us.collection)

	bs := bson.M{"scopes." + clientId: scope}
	bs = bson.M{"$set": bs}
	_, err = c.UpdateOne(context.TODO(), bson.M{"_id": user.GetUserID()}, bs)

	if err != nil {
		return
	}

	user, err = us.Find(id)
	if err != nil {
		return
	}

	addUserCache(user)
	return
}
