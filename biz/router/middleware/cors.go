package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/cors"
	"hertz-starter-kit/config"
	"time"
)

func Cors() func(ctx context.Context, c *app.RequestContext) {
	return cors.New(cors.Config{
		AllowOrigins:  config.Cfg.AllowOrigins,
		AllowMethods:  []string{"POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Cache-Control", "Content-Language", "Content-Type"},
		//AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	})
}
