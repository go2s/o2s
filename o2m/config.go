package o2m

import (
	"time"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

// MongoConfig mongodb configuration parameters
type MongoConfig struct {
	Addrs     []string
	Database  string
	Username  string
	Password  string
	PoolLimit int
}

//NewMongoSession new mongo session
func NewMongoSession(cfg *MongoConfig) *mgo.Session {
	dialInfo := &mgo.DialInfo{
		Addrs:    cfg.Addrs,
		Database: cfg.Database,
		Username: cfg.Username,
		Password: cfg.Password,
		Timeout:  time.Second * 30,
	}

	glog.Infof("mongodb dial info: %+v", dialInfo)

	s, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		glog.Infof("connect mongodb error: %v", err.Error())
		panic(err)
	}
	s.SetMode(mgo.Monotonic, true)
	s.SetPoolLimit(cfg.PoolLimit)

	glog.Infof("mongodb connected")
	return s
}
