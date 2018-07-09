// authors: wangoo
// created: 2018-06-29
// oauth2 err

package o2

import (
	goErr "errors"
	"net/http"
	"gopkg.in/oauth2.v3/errors"
	"log"
)

type httpError interface {
	error
	Code() int
	Status() int
}

type OauthError struct {
	status int
	code   int
	err    error
}

func (e *OauthError) Status() int {
	return e.status
}

func (e *OauthError) Code() int {
	return e.code
}

func (e *OauthError) Error() string {
	return e.err.Error()
}

func NewOauthError(status, code int, err string) *OauthError {
	return &OauthError{
		status: status,
		code:   code,
		err:    goErr.New(err),
	}
}

const (
	ErrCodeInternalError     = 100
	ErrCodeInvalidCredential = 101
	ErrCodeValueRequired     = 200
	ErrCodeNotFound          = 201
	ErrCodeDuplicated        = 202
)

var (
	ErrInternalError     = NewOauthError(http.StatusInternalServerError, ErrCodeInternalError, "internal error")
	ErrInvalidCredential = NewOauthError(http.StatusUnauthorized, ErrCodeInvalidCredential, "invalid credential")
	ErrValueRequired     = NewOauthError(http.StatusBadRequest, ErrCodeValueRequired, "value required")
	ErrNotFound          = NewOauthError(http.StatusNotFound, ErrCodeNotFound, "not found")
	ErrDuplicated        = NewOauthError(http.StatusConflict, ErrCodeDuplicated, "duplicated")
)

func InternalErrorHandler(err error) (re *errors.Response) {
	if herr, ok := err.(httpError); ok {
		re = &errors.Response{
			StatusCode: herr.Status(),
			ErrorCode:  herr.Code(),
			Error:      herr,
		}
		return
	}

	re = &errors.Response{
		StatusCode: ErrInternalError.status,
		ErrorCode:  ErrInternalError.code,
		Error:      err,
	}
	return
}

func ResponseErrorHandler(re *errors.Response) {
	log.Println("Internal Error:", re.Error)
}
