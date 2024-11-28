package model

import (
	"errors"
	"fmt"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"github.com/SuperH-0630/hdangan/src/systeminit"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path"
)

var db *gorm.DB

var gormConfig = &gorm.Config{}

func startDriver(rt runtime.RunTime) (*gorm.DB, error) {
	var err error

	if db != nil {
		return db, nil
	}

	initConfig, err := systeminit.GetInit()
	if errors.Is(err, systeminit.LuckyError) {
		rt.DBConnectError(err)
		return nil, fmt.Errorf("配置文件错误：%s", err.Error())
	} else if err != nil {
		rt.DBConnectError(fmt.Errorf("配置文件错误，请检查配置文件状态。"))
		return nil, fmt.Errorf("配置文件错误：%s", err.Error())
	}

	dbPath := path.Join(initConfig.HomeDir, "data.db")

	db, err = gorm.Open(sqlite.Open(dbPath), gormConfig)
	if err != nil {
		rt.DBConnectError(fmt.Errorf("配置文件错误，请检查配置文件状态。"))
		return nil, fmt.Errorf("数据库链接错误： %s", err.Error())
	}

	return db, nil
}

func GetDB(rt runtime.RunTime) (*gorm.DB, error) {
	return startDriver(rt)
}

func StopDB(rt runtime.RunTime) {
	if db != nil {
	}
}
