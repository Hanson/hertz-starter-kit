package utils

import (
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"hertz-starter-kit/biz/model/admin/customer_service"
)

func GetCs(c *app.RequestContext) (*customer_service.ModelCustomerService, error) {
	if v, ok := c.Get("cs"); !ok {
		return nil, errors.New("unknown cs")
	} else {
		return v.(*customer_service.ModelCustomerService), nil
	}
}
