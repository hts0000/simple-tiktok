package errno

import (
	"fmt"
	"net/http"
)

type ErrNo struct {
	ErrCode int64
	ErrMsg  string
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

func NewErrNo(code int64, msg string) ErrNo {
	return ErrNo{
		ErrCode: code,
		ErrMsg:  msg,
	}
}

func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrMsg = msg
	return e
}

var (
	Success                = NewErrNo(int64(0), "Success")
	ServiceErr             = NewErrNo(int64(http.StatusInternalServerError), "Service is unable to start successfully")
	ParamErr               = NewErrNo(int64(http.StatusBadRequest), "Wrong Parameter has been given")
	UserAlreadyExistErr    = NewErrNo(int64(http.StatusBadRequest), "User already exists")
	UserNotExistErr        = NewErrNo(int64(http.StatusBadRequest), "User not exists")
	VideoNotExistErr       = NewErrNo(int64(http.StatusBadRequest), "Video not exists")
	AuthorizationFailedErr = NewErrNo(int64(http.StatusBadRequest), "Authorization failed")
	PageNotFound           = NewErrNo(int64(http.StatusNotFound), "Page not found")
	MethodNotAllowed       = NewErrNo(int64(http.StatusMethodNotAllowed), "Method not allowed")
)
