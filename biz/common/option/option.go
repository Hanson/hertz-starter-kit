package option

import (
	"hertz-starter-kit/biz/model/common/paginate"
	"log"
)

func HandleOptions(options []*paginate.Options, fun func(typ int32, value string) error) error {
	for _, option := range options {
		err := fun(option.Type, option.Value)
		if err != nil {
			log.Printf("err: %+v", err)
			return err
		}
	}

	return nil
}
