package o2r

import (
	"encoding/json"
	"time"

	"gopkg.in/oauth2.v3/models"
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	"gopkg.in/oauth2.v3"
)

// RedisTokenStore redis token store
type RedisTokenStore struct {
	cli *redis.Client
}

// NewTokenStore Create a token store instance based on redis
func NewTokenStore(cfg *redis.Options) (ts oauth2.TokenStore, err error) {
	if cfg == nil {
		panic("config cannot be nil")
	}
	cli := redis.NewClient(cfg)
	if verr := cli.Ping().Err(); verr != nil {
		err = verr
		return
	}
	ts = &RedisTokenStore{cli: cli}
	return
}

// Create Create and store the new token information
func (rs *RedisTokenStore) Create(info oauth2.TokenInfo) (err error) {
	ct := time.Now()
	jv, err := json.Marshal(info)
	if err != nil {
		return
	}

	pipe := rs.cli.Pipeline()
	if code := info.GetCode(); code != "" {
		pipe.Set(code, jv, info.GetCodeExpiresIn())
	} else {
		basicID := uuid.NewV4().String()
		aexp := info.GetAccessExpiresIn()
		rexp := aexp

		if refresh := info.GetRefresh(); refresh != "" {
			rexp = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Sub(ct)
			if aexp.Seconds() > rexp.Seconds() {
				aexp = rexp
			}
			pipe.Set(refresh, basicID, rexp)
		}

		pipe.Set(info.GetAccess(), basicID, aexp)
		pipe.Set(basicID, jv, rexp)
	}

	if _, verr := pipe.Exec(); verr != nil {
		err = verr
	}
	return
}

// remove
func (rs *RedisTokenStore) remove(key string) (err error) {
	_, verr := rs.cli.Del(key).Result()
	if verr != redis.Nil {
		err = verr
	}
	return
}

// RemoveByCode Use the authorization code to delete the token information
func (rs *RedisTokenStore) RemoveByCode(code string) (err error) {
	err = rs.remove(code)
	return
}

// RemoveByAccess Use the access token to delete the token information
func (rs *RedisTokenStore) RemoveByAccess(access string) (err error) {
	basicID, err := rs.getBasicID(access)
	if err != nil || basicID == "" {
		return
	}
	rs.remove(access)
	ti, err := rs.getData(basicID)
	if err == nil && ti != nil {
		rs.remove(ti.GetRefresh())
	}
	err = rs.remove(basicID)
	return
}

// RemoveByRefresh Use the refresh token to delete the token information
func (rs *RedisTokenStore) RemoveByRefresh(refresh string) (err error) {
	basicID, err := rs.getBasicID(refresh)
	if err != nil || basicID == "" {
		return
	}
	rs.remove(refresh)
	ti, err := rs.getData(basicID)
	if err == nil && ti != nil {
		rs.remove(ti.GetAccess())
	}
	err = rs.remove(basicID)
	return
}

func (rs *RedisTokenStore) getData(key string) (ti oauth2.TokenInfo, err error) {
	result := rs.cli.Get(key)
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
	var tm models.Token
	if verr := json.Unmarshal(iv, &tm); verr != nil {
		err = verr
		return
	}
	ti = &tm
	return
}

func (rs *RedisTokenStore) getBasicID(token string) (basicID string, err error) {
	tv, verr := rs.cli.Get(token).Result()
	if verr != nil {
		if verr == redis.Nil {
			return
		}
		err = verr
		return
	}
	basicID = tv
	return
}

// GetByCode Use the authorization code for token information data
func (rs *RedisTokenStore) GetByCode(code string) (ti oauth2.TokenInfo, err error) {
	ti, err = rs.getData(code)
	return
}

// GetByAccess Use the access token for token information data
func (rs *RedisTokenStore) GetByAccess(access string) (ti oauth2.TokenInfo, err error) {
	basicID, err := rs.getBasicID(access)
	if err != nil || basicID == "" {
		return
	}
	ti, err = rs.getData(basicID)
	return
}

// GetByRefresh Use the refresh token for token information data
func (rs *RedisTokenStore) GetByRefresh(refresh string) (ti oauth2.TokenInfo, err error) {
	basicID, err := rs.getBasicID(refresh)
	if err != nil || basicID == "" {
		return
	}
	ti, err = rs.getData(basicID)
	return
}
