// authors: wangoo
// created: 2018-06-29
// oauth2 err

package o2

import (
	"github.com/go2s/o2s/o2x"
	"github.com/golang/glog"
	"gopkg.in/oauth2.v3/errors"
)

func InternalErrorHandler(err error) (re *errors.Response) {
	if herr, ok := err.(o2x.CodeError); ok {
		re = &errors.Response{
			StatusCode: herr.Status(),
			ErrorCode:  herr.Code(),
			Error:      herr,
		}
		return
	}

	re = &errors.Response{
		StatusCode: o2x.ErrInternalError.Status(),
		ErrorCode:  o2x.ErrInternalError.Code(),
		Error:      err,
	}
	return
}

func ResponseErrorHandler(re *errors.Response) {
	glog.Errorf("Error:%v", re.Error)
}
