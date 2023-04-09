package errors

import "hertz-starter-kit/biz/model/common/err_code"

type Err struct {
	ErrCode int
	ErrMsg  string
}

func (e Err) Error() string {
	return e.ErrMsg
}

func New(errCode int, errMsg string) *Err {
	return &Err{ErrCode: errCode, ErrMsg: errMsg}
}

func System(errMsg string) *Err {
	if errMsg == "" {
		errMsg = "系统错误"
	}

	return New(int(err_code.ErrCode_ErrSystem), errMsg)
}

func GetErrCode(err error) int {
	if err == nil {
		return 0
	}
	if p, ok := err.(*Err); ok {
		return p.ErrCode
	}
	return 1000
}
