// authors: wangoo
// created: 2018-07-20

package o2x

import (
	"time"

	"github.com/golang/glog"
	"github.com/patrickmn/go-cache"
	"gopkg.in/oauth2.v3"
)

const (
	CaptchaCredentials oauth2.GrantType = "captcha"
)

type CaptchaStore interface {
	Save(mobile, captcha string) (err error)
	Remove(mobile string) (err error)
	Valid(mobile, captcha string) (valid bool, err error)
}

type MemoryCaptchaStore struct {
	c *cache.Cache
}

func NewMemoryCaptchaStore(expireDuration time.Duration) (cs *MemoryCaptchaStore, err error) {
	cs = &MemoryCaptchaStore{
		c: cache.New(expireDuration, 2*expireDuration),
	}
	return
}

func (cs *MemoryCaptchaStore) Save(mobile, captcha string) (err error) {
	glog.Infof("save captcha:%v,%v", mobile, captcha)
	cs.c.Add(mobile, captcha, cache.DefaultExpiration)
	return
}

func (cs *MemoryCaptchaStore) Remove(mobile string) (err error) {
	cs.c.Delete(mobile)
	return
}

func (cs *MemoryCaptchaStore) Valid(mobile, captcha string) (valid bool, err error) {
	if c, ok := cs.c.Get(mobile); ok {
		cap := c.(string)
		if cap == captcha {
			valid = true
			go cs.Remove(mobile)
			return
		}
	}

	valid = false
	return
}
