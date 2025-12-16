package cache

import (
	"context"
	"fmt"
	"log"
)

type logger struct {
	log *log.Logger
}

func (l *logger) Printf(ctx context.Context, format string, v ...interface{}) {
	var traceId string
	if ctx != nil {
		if ctx.Value("trace_id") != nil {
			traceId = ctx.Value("trace_id").(string)
		}
	}
	_ = l.log.Output(2, fmt.Sprintf("<"+traceId+">"+format, v...))
}
