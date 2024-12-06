package model

import (
	"errors"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"gorm.io/gorm"
)

type CountAll struct {
	Res int64 `gorm:"column:ps"`
}

type CountAllWithFile struct {
	File string `gorm:"column:filetype"`
	Res  int64  `gorm:"column:ps"`
}

func CountFile(rt runtime.RunTime, model interface{}) (int64, error) {
	db, err := GetDB(rt)
	if err != nil {
		return 0, err
	}

	var res CountAll
	err = db.Model(model).Select("COUNT(*) AS ps").First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return res.Res, err
}
