// authors: wangoo
// created: 2018-07-10

package o2x

import (
	goErr "errors"
	"net/http"
)

type CodeError interface {
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
	ErrCodeInvalidCaptcha    = 102
	ErrCodeValueRequired     = 200
	ErrCodeNotFound          = 201
	ErrCodeDuplicated        = 202
)

var (
	ErrInternalError     = NewOauthError(http.StatusInternalServerError, ErrCodeInternalError, "internal error")
	ErrInvalidCredential = NewOauthError(http.StatusUnauthorized, ErrCodeInvalidCredential, "invalid credential")
	ErrInvalidCaptcha    = NewOauthError(http.StatusUnauthorized, ErrCodeInvalidCaptcha, "invalid captcha")
	ErrValueRequired     = NewOauthError(http.StatusBadRequest, ErrCodeValueRequired, "value required")
	ErrNotFound          = NewOauthError(http.StatusNotFound, ErrCodeNotFound, "not found")
	ErrDuplicated        = NewOauthError(http.StatusConflict, ErrCodeDuplicated, "duplicated")
)
