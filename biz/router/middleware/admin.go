package middleware

import (
	"context"
	"errors"
	"hertz-starter-kit/biz/model/public/auth"
	"hertz-starter-kit/config"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/jwt"
)

var (
	JwtAdminMiddleware *jwt.HertzJWTMiddleware
	IdentityKey        = "identity"
)

func InitAdminJwt() {

	type SuperAdmin struct {
		Username string `json:"username" validate:"required"`
	}
	var err error
	JwtAdminMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:         "admin",
		Key:           []byte("wtsdf3w4tsdrg34"),
		Timeout:       time.Hour * 24 * 15,
		MaxRefresh:    time.Hour,
		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, utils.H{
				"code":    code,
				"token":   token,
				"expire":  expire.Format(time.RFC3339),
				"message": "success",
			})
		},
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var req auth.AdminLoginReq
			if err := c.BindAndValidate(&req); err != nil {
				return nil, err
			}

			if req.Username == config.Cfg.AdminUsername && req.Password == config.Cfg.AdminPassword {
				return SuperAdmin{Username: req.Username}, nil
			}

			return nil, errors.New("账号或密码不正确")
		},
		IdentityKey: "admin",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			return jwt.MapClaims{
				"role": "administrator",
			}
		},
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)

			if claims["role"] != "administrator" {
				return nil
			}

			return &SuperAdmin{Username: config.Cfg.AdminUsername}
		},
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			hlog.CtxErrorf(ctx, "jwt biz err = %+v", e.Error())
			return e.Error()
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(http.StatusOK, utils.H{
				"code":    code,
				"message": message,
			})
		},
	})
	if err != nil {
		panic(err)
	}
}
