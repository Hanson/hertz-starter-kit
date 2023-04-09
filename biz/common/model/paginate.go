package model

import (
	"gorm.io/gorm"
	"hertz-starter-kit/biz/model/common/paginate"
	"log"
)

var defaultOffset int = 0
var defaultLimit int = 20

func Paginate(m *gorm.DB, listOption *paginate.ListOption, list interface{}) (*paginate.Paginate, error) {
	var offset, limit int

	p := paginate.Paginate{}

	if listOption.Offset == 0 {
		offset = defaultOffset
	} else {
		offset = int(listOption.Offset)
	}

	if listOption.Limit == 0 {
		limit = defaultLimit
	} else {
		limit = int(listOption.Limit)
	}

	err := m.Offset(offset).Limit(limit).Find(list).Error
	if err != nil {
		log.Printf("err: %+v", err)
		return nil, err
	}

	if listOption.NeedCount && offset == 0 {
		err = m.Model(list).Count(&p.Total).Error
		if err != nil {
			log.Printf("err: %+v", err)
			return nil, err
		}
	}

	p.Limit = uint32(limit)
	p.Offset = uint32(offset)

	return &p, err
}
