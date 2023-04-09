package model

import (
	"gorm.io/gorm"
	"hertz-starter-kit/biz/model/common/paginate"
	"log"
	"strconv"
	"strings"
)

func HandleListOption(listOption *paginate.ListOption, m *gorm.DB, handlers map[int32]func(value string) error) error {
	m = handleDefaultListOption(listOption, m)

	for _, option := range listOption.Options {
		if fun, ok := handlers[option.Type]; ok {
			err := fun(option.Value)
			if err != nil {
				log.Printf("err: %+v", err)
				return err
			}
		}
	}

	return nil
}

func handleDefaultListOption(listOption *paginate.ListOption, m *gorm.DB) *gorm.DB {
	for _, option := range listOption.Options {
		switch option.Type {
		case int32(paginate.DefaultListOption_ListOptionIdList):
			if option.Value == "" {
				m.Where("1 = 0")
			} else {
				v := strings.Split(option.Value, ",")
				m.Where("id IN ?", v)
			}
		case int32(paginate.DefaultListOption_ListOptionOrderBy):
			if option.Value == strconv.Itoa(int(paginate.DefaultOrderBy_OrderByIdAsc)) {
				m.Order("id ASC")
			} else if option.Value == strconv.Itoa(int(paginate.DefaultOrderBy_OrderByIdDesc)) {
				m.Order("id DESC")
			}
		}
	}

	return m
}
