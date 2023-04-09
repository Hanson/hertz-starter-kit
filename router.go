// Code generated by hertz generator.

package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
	"hertz-starter-kit/biz/handler/ws"
	"hertz-starter-kit/biz/router/middleware"
	_ "hertz-starter-kit/docs"
)

// customizeRegister registers customize routers.
func customizedRegister(r *server.Hertz) {
	// your code ...
	url := swagger.URL("http://localhost:8888/swagger/doc.json")
	r.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))

	// 管理后台授权
	adminRoute := r.Group("/admin")
	adminRoute.POST("/login", middleware.JwtAdminMiddleware.LoginHandler)

	r.GET("/echo", ws.Index)
}
