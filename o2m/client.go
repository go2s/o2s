// Oauth2 client mongo store
// authors: wongoo

package o2m

import (
	"context"
	"time"

	"github.com/go2s/o2s/o2x"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/oauth2.v3"
)

const (
	DefaultOauth2ClientDb         = "oauth2"
	DefaultOauth2ClientCollection = "client"
)

var (
	clientCache *cache.Cache
)

func init() {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	clientCache = cache.New(5*time.Minute, 10*time.Minute)
}

func addClientCache(cli oauth2.ClientInfo) {
	clientCache.Add(cli.GetID(), cli, cache.DefaultExpiration)
}

func getClientCache(id string) (cli oauth2.ClientInfo) {
	if c, found := clientCache.Get(id); found {
		cli = c.(oauth2.ClientInfo)
		return
	}
	return
}

// Mongo client store
type MongoClientStore struct {
	db         string
	collection string
	client     *mongo.Client
}

type Oauth2Client struct {
	ID         string             `bson:"_id" json:"id"`
	Secret     string             `bson:"secret" json:"secret"`
	Domain     string             `bson:"domain" json:"domain"`
	Scopes     []string           `bson:"scopes" json:"scopes"`
	GrantTypes []oauth2.GrantType `bson:"grant_types" json:"grant_types"`
	UserID     string             `bson:"user_id,omitempty" json:"user_id,omitempty"`
}

func (c *Oauth2Client) GetID() string {
	return c.ID
}
func (c *Oauth2Client) GetSecret() string {
	return c.Secret
}
func (c *Oauth2Client) GetDomain() string {
	return c.Domain
}
func (c *Oauth2Client) GetScopes() []string {
	return c.Scopes
}
func (c *Oauth2Client) GetGrantTypes() []oauth2.GrantType {
	return c.GrantTypes
}
func (c *Oauth2Client) GetUserID() string {
	return c.UserID
}

func NewClientStore(client *mongo.Client, db string, collection string) (clientStore *MongoClientStore) {
	if client == nil {
		panic("client cannot be nil")
	}
	clientStore = &MongoClientStore{client: client, db: db, collection: collection}
	if clientStore.db == "" {
		clientStore.db = DefaultOauth2ClientDb
	}
	if clientStore.collection == "" {
		clientStore.collection = DefaultOauth2ClientCollection
	}

	return
}

// GetByID according to the ID for the client information
func (cs *MongoClientStore) GetByID(id string) (cli oauth2.ClientInfo, err error) {
	if cli = getClientCache(id); cli != nil {
		return
	}

	c := cs.client.Database(cs.db).Collection(cs.collection)
	query := c.FindOne(context.TODO(), bson.M{"_id": id})
	client := &Oauth2Client{}
	err = query.Decode(client)
	if err != nil {
		return nil, err
	}

	addClientCache(client)
	return client, nil
}

// Add a client info
func (cs *MongoClientStore) Set(id string, cli oauth2.ClientInfo) (err error) {
	client := &Oauth2Client{
		ID:     cli.GetID(),
		UserID: cli.GetUserID(),
		Domain: cli.GetDomain(),
		Secret: cli.GetSecret(),
	}

	if o2ClientInfo, ok := cli.(o2x.O2ClientInfo); ok {
		client.Scopes = o2ClientInfo.GetScopes()
		client.GrantTypes = o2ClientInfo.GetGrantTypes()
	}
	addClientCache(client)
	c := cs.client.Database(cs.db).Collection(cs.collection)
	_, err = c.InsertOne(context.TODO(), client)
	return
}
