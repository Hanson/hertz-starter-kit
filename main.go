// Code generated by hertz generator.

package main

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/logger/accesslog"
	"hertz-starter-kit/biz/router/middleware"
	"hertz-starter-kit/config"
	"hertz-starter-kit/db"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	// 配置
	log.SetFlags(log.Llongfile | log.LstdFlags)
	log.SetOutput(getMultiWriter())

	config.LoadConfig()

	// db 操作
	err := db.InitDb()
	if err != nil {
		log.Printf("err: %+v", err)
		return
	}

	err = db.AutoMigrate()
	if err != nil {
		log.Printf("err: %+v", err)
		return
	}

	h := server.Default(
		server.WithAutoReloadRender(true, 0),
		server.WithHostPorts(config.Cfg.HostPort),
	)

	corsMiddleware(h)
	h.Use(accesslog.New(accesslog.WithFormat("[${latency}] method: ${method} path: ${path} req_body: ${body} rsp_body: ${resBody}")))
	middleware.InitAdminJwt()
	register(h)
	h.Spin()
}

func corsMiddleware(h *server.Hertz) {
	h.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Cache-Control", "Content-Language", "Content-Type"},
		//AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))
}

func getMultiWriter() io.Writer {
	f, err := os.OpenFile(fmt.Sprintf("logs\\%s.txt", time.Now().Format("2006-01-02")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	return io.MultiWriter(os.Stdout, f)
}
