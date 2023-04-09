package ws

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/getsentry/sentry-go"
	"github.com/hertz-contrib/hertzsentry"
	"github.com/hertz-contrib/websocket"
	"log"
)

var upgrader = websocket.HertzUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(ctx *app.RequestContext) bool {
		return true
	},
	Error: func(ctx *app.RequestContext, status int, reason error) {
		log.Println(status)
		log.Println(reason)
	},
}

func Index(ctx context.Context, c *app.RequestContext) {
	var err error

	if hub := hertzsentry.GetHubFromContext(c); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("hertz", "CloudWeGo Hertz")
			scope.SetLevel(sentry.LevelDebug)
			hub.CaptureMessage("Just for debug")
		})
	}

	err = upgrader.Upgrade(c, func(conn *websocket.Conn) {
		//ws.Pool.AddConn(cs.Id, conn)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				//ws.Pool.RemoveConn(cs.Id, conn)
				log.Println("read:", err)
				break
			}

			log.Printf("recv: %s", message)

		}
	})
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
}
