package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"hertz-starter-kit/utils/log"
)

func RecoveryHandler(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {
	log.Errorf(ctx, "[Recovery] err=%v\nstack=%s", err, stack)
	c.AbortWithStatus(consts.StatusInternalServerError)
}
