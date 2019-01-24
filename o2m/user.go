// authors: wangoo
// created: 2018-05-30
// oauth2 user store

package o2m

import (
	"gopkg.in/mgo.v2"
	"github.com/go2s/o2s/o2x"
	"reflect"
	"gopkg.in/mgo.v2/bson"
	"github.com/golang/glog"
	"github.com/patrickmn/go-cache"
	"time"
	"fmt"
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
	session          *mgo.Session
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

func NewUserStore(session *mgo.Session, db, collection string, userCfg *MgoUserCfg) (us *MgoUserStore) {
	if !o2x.IsUserType(userCfg.userType) {
		panic("invalid user type")
	}
	us = &MgoUserStore{
		session:          session,
		db:               db,
		collection:       collection,
		mobileCollection: collection + "_mobile",
		userCfg:          userCfg,
	}

	err := session.DB(us.db).C(us.mobileCollection).EnsureIndex(mgo.Index{
		Key:    []string{"mobile"},
		Unique: true,
	})
	if err != nil {
		panic(err)
	}

	return
}

func (us *MgoUserStore) lockUserMobile(session *mgo.Session, userId, mobile string) (err error) {
	if userId == "" || mobile == "" {
		err = o2x.ErrValueRequired
		return
	}
	c := session.DB(us.db).C(us.mobileCollection)
	userMobile := &MgoUserMobile{
		Id:     userId,
		Mobile: mobile,
	}
	err = c.Insert(userMobile)
	return
}

func (us *MgoUserStore) unlockUserMobile(session *mgo.Session, userId string) (err error) {
	if userId == "" {
		err = o2x.ErrValueRequired
		return
	}
	c := session.DB(us.db).C(us.mobileCollection)
	mgoErr := c.RemoveId(userId)
	if mgoErr != mgo.ErrNotFound {
		err = mgoErr
	}
	return
}

func (us *MgoUserStore) Save(u o2x.User) (err error) {
	session := us.session.Clone()
	defer session.Close()

	if u.GetMobile() != "" {
		err = us.lockUserMobile(session, u.GetID(), u.GetMobile())
		if err != nil {
			return
		}
	}
	c := session.DB(us.db).C(us.collection)
	glog.Infof("insert user:%v", u)
	err = c.Insert(u)

	if err != nil {
		return
	}
	addUserCache(u)

	return
}

func (us *MgoUserStore) Remove(id interface{}) (err error) {
	removeUserCache(id)

	session := us.session.Clone()
	defer session.Close()
	c := session.DB(us.db).C(us.collection)

	sid, err := o2x.UserIdString(id)
	if err != nil {
		return
	}

	glog.Infof("remove user:%v", id)
	us.unlockUserMobile(session, sid)

	mgoErr := c.RemoveId(id)
	if mgoErr != nil && mgoErr == mgo.ErrNotFound {
		// try to find using object id
		if sid, ok := id.(string); ok && bson.IsObjectIdHex(sid) {
			bid := bson.ObjectIdHex(sid)
			mgoErr = c.RemoveId(bid)
		}
	}
	if mgoErr != nil && mgoErr == mgo.ErrNotFound {
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

	session := us.session.Clone()
	defer session.Close()
	c := session.DB(us.db).C(us.collection)

	user := o2x.NewUser(us.userCfg.userType)
	mgoErr := c.FindId(id).One(user)
	if mgoErr != nil && mgoErr == mgo.ErrNotFound {
		// try to find using object id
		if sid, ok := id.(string); ok && bson.IsObjectIdHex(sid) {
			bid := bson.ObjectIdHex(sid)
			mgoErr = c.FindId(bid).One(user)
		}
	}

	if mgoErr != nil {
		if mgoErr == mgo.ErrNotFound {
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
	session := us.session.Clone()
	defer session.Close()
	c := session.DB(us.db).C(us.collection)

	user := o2x.NewUser(us.userCfg.userType)
	mgoErr := c.Find(bson.M{"mobile": mobile}).One(user)
	if mgoErr != nil && mgoErr == mgo.ErrNotFound {
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

	session := us.session.Clone()
	defer session.Close()
	c := session.DB(us.db).C(us.collection)

	bs := bson.M{us.userCfg.passwordName: user.GetPassword(), us.userCfg.saltName: user.GetSalt()}
	bs = bson.M{"$set": bs}
	err = c.UpdateId(user.GetUserID(), bs)

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

	session := us.session.Clone()
	defer session.Close()
	c := session.DB(us.db).C(us.collection)

	bs := bson.M{"scopes." + clientId: scope}
	bs = bson.M{"$set": bs}
	err = c.UpdateId(user.GetUserID(), bs)

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
