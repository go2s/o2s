// authors: wangoo
// created: 2018-06-29
// oauth2 err

package o2

import "errors"

var (
	ErrValueRequired = errors.New("value required")
	ErrNotFound      = errors.New("not found")
	ErrDuplicated    = errors.New("duplicated")
)
