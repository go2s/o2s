package o2m

import (
	"context"
	"time"

	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConfig mongodb configuration parameters
type MongoConfig struct {
	Hosts     []string
	Database  string
	Username  string
	Password  string
	PoolLimit uint64
}

//NewMongoClient new mongo client
func NewMongoClient(cfg *MongoConfig) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	option := options.Client().SetAuth(options.Credential{Username: cfg.Username, Password: cfg.Password, AuthSource: cfg.Database})
	option.SetConnectTimeout(time.Second * 5)
	option.SetMaxPoolSize(cfg.PoolLimit)
	option.SetHosts(cfg.Hosts)

	client, err := mongo.Connect(ctx, option)
	defer cancel()

	if err != nil {
		glog.Infof("connect mongodb error: %v", err.Error())
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	//err = client.Ping(ctx, readpref.Primary())
	glog.Infof("mongodb connected")
	return client
}
