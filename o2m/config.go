package o2m

import (
	"context"
	"github.com/golang/glog"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"time"
)

// MongoConfig mongodb configuration parameters
type MongoConfig struct {
	Address   string
	Database  string
	Username  string
	Password  string
	PoolLimit int
}

//NewMongoClient new mongo client
func NewMongoClient(cfg *MongoConfig) *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	option := options.Client().SetAuth(options.Credential{Username: cfg.Username, Password: cfg.Password, AuthSource: cfg.Database})
	option.SetConnectTimeout(time.Second * 5)
	option.SetMaxPoolSize(uint16(cfg.PoolLimit))
	client, err := mongo.Connect(ctx, cfg.Address, option)

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
