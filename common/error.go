package common

import "fmt"

const (
	ErrNotFound int = 1000
	ErrUserExists  int = 1001
	ErrTODO int = 1003
	)

var ErrorCodes = map[int]string{
	ErrNotFound:            "does not exists error",
}
type RuntimeErr struct {
	ErrCode          int    // code, see error constants below
	ErrMsg           string // optional message (otherwise resolved by errorMap)
	HTTPResponseCode int    // optional custom HTTP response code. default response codes are in ErrorCodes
	WrappedError     error  // optional nested error if there was a different root cause
}

func (e RuntimeErr) Error() string {
	var msg = e.ErrMsg

	if msg == "" {
		msg = ErrorCodes[e.ErrCode]
	}

	if msg == "" {
		msg = fmt.Sprintf("RuntimeErr with ErrCode=%d", e.ErrCode)
	}

	if e.WrappedError != nil {
		msg = fmt.Sprintf("%s (cause: %v)", msg, e.WrappedError)
	}

	return msg
}


func NewRuntimeError(ErrCode int) RuntimeErr {
	return RuntimeErr{
		ErrCode: ErrCode,
	}
}

func (re RuntimeErr) WithHttpResponseCode(httpResponseCode int) RuntimeErr {
	re.HTTPResponseCode = httpResponseCode
	return re
}
func (re RuntimeErr) WithMessagef(format string, args ...interface{}) RuntimeErr {
	re.ErrMsg = fmt.Sprintf(format, args...)
	return re
}