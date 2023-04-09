package response

import (
	"github.com/cloudwego/hertz/pkg/app"
	"hertz-starter-kit/biz/common/errors"
	"hertz-starter-kit/biz/model/common/err_code"
	"hertz-starter-kit/feishu"
	"log"
)

type Response struct {
	ErrCode int         `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
	Data    interface{} `json:"data"`
}

func JsonOk(c *app.RequestContext, data interface{}) {
	resp := Response{
		ErrCode: 0,
		ErrMsg:  "success",
		Data:    data,
	}
	c.JSON(200, resp)
}

func JsonErr(c *app.RequestContext, err error) {
	go func() {
		feishu.SendText(err.Error())
	}()
	log.Printf("err: %s", err.Error())
	errMsg := err.Error()
	if errors.GetErrCode(err) == 1000 {
		errMsg = "系统错误"
	}
	resp := Response{
		ErrCode: errors.GetErrCode(err),
		ErrMsg:  errMsg,
	}
	c.JSON(200, resp)
}

func JsonSystemErr(c *app.RequestContext, err error) {
	JsonErr(c, errors.System(err.Error()))
}

func JsonValidateErr(c *app.RequestContext, err error) {
	resp := Response{
		ErrCode: int(err_code.ErrCode_ErrValidate),
		ErrMsg:  "参数校验不通过 err:" + err.Error(),
	}
	c.JSON(200, resp)
}
