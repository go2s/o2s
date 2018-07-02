// authors: wangoo
// created: 2018-06-29
// oauth2 err

package o2

type ErrorCoder interface {
	error
	ErrorCode() string
}

type CodeError struct {
	code    string
	message string
}

func (e *CodeError) ErrorCode() string {
	return e.code
}

func (e *CodeError) Error() string {
	return e.message
}

func NewCodeError(code, message string) *CodeError {
	return &CodeError{
		code:    code,
		message: message,
	}
}

var (
	ErrValueRequired = NewCodeError("E100", "value required")
	ErrNotFound      = NewCodeError("E101", "not found")
	ErrDuplicated    = NewCodeError("E102", "duplicated")
)
