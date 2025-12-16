package middleware

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hanson/go-toolbox/utils"
	"gorm.io/gorm/logger"
	"hertz-starter-kit/utils/log"
	"strconv"
	"strings"
	"sync"
	"time"
)

func LogTrace() func(ctx context.Context, c *app.RequestContext) {
	var (
		wsEchoPrefix = []byte("/ws/echo")
		weworkPrefix = []byte("WeworkThirdCallback")
	)

	return func(ctx context.Context, c *app.RequestContext) {
		//uri := c.Request.URI().RequestURI()
		uri := c.Request.URI().Path()
		if bytes.Contains(uri, wsEchoPrefix) || bytes.Contains(uri, weworkPrefix) {
			c.Next(ctx)
			return
		}

		traceId := genTraceID()
		ctx = context.WithValue(ctx, "trace_id", traceId)

		body := c.Request.Body()
		if len(body) > 1024 {
			body = []byte(fmt.Sprintf("too large [%d bytes]", len(body)))
		}

		buf := bytes.NewBuffer(make([]byte, 0, 256))
		buf.WriteString(logger.Green)
		buf.WriteString(" path:")
		buf.Write(c.Request.Path())
		buf.WriteString(logger.Reset)
		buf.WriteString(" body:")
		buf.Write(body)
		log.Infof(ctx, buf.String())

		//log.Infof(ctx, "%s path:%s%s body:%s", logger.Green, c.Request.Path(), logger.Reset, body)

		c.Next(ctx)

		c.Response.Header.Set("X-Trace-ID", traceId)
	}
}

var builderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

func genTraceID() string { // 1. 从 Pool 获取 Builder
	builder := builderPool.Get().(*strings.Builder)
	defer builderPool.Put(builder) // 用完后放回 Pool

	// 2. 重置 Builder（避免残留数据）
	builder.Reset()

	builder.WriteString(strconv.FormatInt(time.Now().Unix(), 10))
	//sb.WriteString("-")
	builder.WriteString(utils.RandStr(4, utils.RandomStringModNumber))
	return builder.String()
}
