package middleware

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hanson/go-toolbox/utils"
	"github.com/tidwall/sjson"
	"gorm.io/gorm/logger"
	"hertz-starter-kit/utils/log"
	"strings"
	"time"
)

func LogTrace() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		if strings.Contains(c.Request.URI().String(), "/ws/echo") {
			c.Next(ctx)
			return
		}
		reqHeaders := make([]string, 0)
		c.Request.Header.VisitAll(func(k, v []byte) {
			reqHeaders = append(reqHeaders, string(k)+"="+string(v))
		})
		traceId := fmt.Sprintf("%d%s", time.Now().UnixMicro(), utils.RandStr(8, utils.RandomStringModNumber))
		ctx = context.WithValue(ctx, "trace_id", traceId)
		log.Infof(ctx, "req %strace_id:%s path:%s%s header:%+v, body:%s", logger.Green, traceId, c.Request.Path(), logger.Reset, strings.Join(reqHeaders, "&"), c.Request.Body())
		beginUnix := time.Now().UnixMilli()
		c.Next(ctx)
		resetBody, _ := sjson.Set(string(c.Response.Body()), "hint", traceId)
		c.Response.SetBody([]byte(resetBody))
		endUnix := time.Now().UnixMilli()
		log.Infof(ctx, "rsp %strace_id:%s%s body:%s [%dms]", logger.Green, traceId, logger.Reset, c.Response.Body(), endUnix-beginUnix)
	}
}
