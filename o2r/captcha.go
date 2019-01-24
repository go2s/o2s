// authors: wangoo
// created: 2018-07-20

package o2r

import (
	"github.com/go-redis/redis"
	"time"
)

const (
	redisCaptchaPrefix = "o2x_oauth2_captcha_"
)

type RedisCaptchaStore struct {
	cli            *redis.Client
	expireDuration time.Duration
}

func NewRedisCaptchaStore(cfg *redis.Options, expireDuration time.Duration) (cs *RedisCaptchaStore, err error) {
	if cfg == nil {
		panic("config cannot be nil")
	}
	cli := redis.NewClient(cfg)
	if verr := cli.Ping().Err(); verr != nil {
		err = verr
		return
	}
	cs = &RedisCaptchaStore{cli: cli, expireDuration: expireDuration}
	return
}

func (cs *RedisCaptchaStore) Save(mobile, captcha string) (err error) {
	result := cs.cli.Set(redisCaptchaPrefix+mobile, captcha, cs.expireDuration)
	if verr := result.Err(); verr != nil {
		if verr == redis.Nil {
			return
		}
		err = verr
		return
	}
	return
}

func (cs *RedisCaptchaStore) Remove(mobile string) (err error) {
	result := cs.cli.Del(redisCaptchaPrefix + mobile)
	if verr := result.Err(); verr != nil {
		if verr == redis.Nil {
			return
		}
		err = verr
		return
	}
	return
}

func (cs *RedisCaptchaStore) Valid(mobile, captcha string) (valid bool, err error) {
	valid = false

	result := cs.cli.Get(redisCaptchaPrefix + mobile)
	if verr := result.Err(); verr != nil {
		if verr == redis.Nil {
			return
		}
		err = verr
		return
	}
	b, err := result.Bytes()
	if err != nil {
		return
	}
	ca := string(b)
	if ca == captcha {
		valid = true
		go cs.Remove(mobile)
		return
	}

	return
}
