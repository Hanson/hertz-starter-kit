// Code generated by hertz generator. DO NOT EDIT.

package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	auth "hertz-starter-kit/biz/router/auth"
	err_code "hertz-starter-kit/biz/router/err_code"
)

// GeneratedRegister registers routers generated by IDL.
func GeneratedRegister(r *server.Hertz) {
	//INSERT_POINT: DO NOT DELETE THIS LINE!
	err_code.Register(r)

	auth.Register(r)

}