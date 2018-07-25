// authors: wangoo
// created: 2018-07-24

package o2

import (
	"net/http"
	"github.com/go2s/o2x"
	"github.com/golang/glog"
	"gopkg.in/oauth2.v3/errors"
)

type CaptchaSender func(mobile, captcha string) (err error)

func CaptchaLogSender(mobile, captcha string) (err error) {
	glog.Infof("captcha console sender:%v,%v", mobile, captcha)
	return
}

func SendCaptchaHandler(w http.ResponseWriter, r *http.Request) {
	err := SendCaptcha(w, r)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	response(w, "ok", http.StatusOK)
}

func SendCaptcha(w http.ResponseWriter, r *http.Request) (err error) {
	mobile := r.FormValue("mobile")
	if mobile == "" {
		err = o2x.ErrValueRequired
		return
	}

	_, err = oauth2UserStore.FindMobile(mobile)
	if err != nil {
		return
	}

	clientID, err := ClientBasicAuth(r)
	if err != nil {
		return
	}

	if fn := oauth2Svr.ClientAuthorizedHandler; fn != nil {
		allowed, verr := fn(clientID, o2x.Captcha)
		if verr != nil {
			err = verr
			return
		} else if !allowed {
			err = errors.ErrUnauthorizedClient
			return
		}
	}

	captcha := "123456"

	err = oauth2CaptchaSender(mobile, captcha)
	if err != nil {
		return
	}

	err = oauth2CaptchaStore.Save(mobile, captcha)
	if err != nil {
		return
	}

	return
}
