// Oauth2 client storage based on redis
// authors: wongoo

package o2r

import (
	"encoding/json"

	"github.com/go-redis/redis"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

const RedisKeyOauth2Client = "oauth2_client"

// RedisTokenStore redis token store
type RedisClientStore struct {
	cli *redis.Client
}

type ClientValue struct {
	Secret string `bson:"secret" json:"secret"`
	Domain string `bson:"domain" json:"domain"`
	UserID string `bson:"user_id" json:"user_id"`
}

func NewClientStore(cfg *redis.Options) (clientStore *RedisClientStore, err error) {
	if cfg == nil {
		panic("config cannot be nil")
	}
	cli := redis.NewClient(cfg)
	if verr := cli.Ping().Err(); verr != nil {
		err = verr
		return
	}
	clientStore = &RedisClientStore{cli: cli}
	return
}

// GetByID according to the ID for the client information
func (cs *RedisClientStore) GetByID(id string) (cli oauth2.ClientInfo, err error) {
	result := cs.cli.HGet(RedisKeyOauth2Client, id)
	if verr := result.Err(); verr != nil {
		if verr == redis.Nil {
			return
		}
		err = verr
		return
	}
	iv, err := result.Bytes()
	if err != nil {
		return
	}
	var cv ClientValue
	if verr := json.Unmarshal(iv, &cv); verr != nil {
		err = verr
		return
	}

	cli = &models.Client{ID: id, Secret: cv.Secret, Domain: cv.Domain, UserID: cv.UserID}
	return
}

// Set set client information
func (cs *RedisClientStore) Set(id string, cli oauth2.ClientInfo) (err error) {
	cv := &ClientValue{Secret: cli.GetSecret(), Domain: cli.GetDomain(), UserID: cli.GetUserID()}
	v, err := json.Marshal(cv)
	if err != nil {
		return
	}
	result := cs.cli.HSet(RedisKeyOauth2Client, cli.GetID(), v)
	if verr := result.Err(); verr != nil {
		if verr != redis.Nil {
			err = verr
			return
		}
	}
	return
}
