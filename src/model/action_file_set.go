package model

import (
	"errors"
	"fmt"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"github.com/SuperH-0630/hdangan/src/systeminit"
	"gorm.io/gorm"
)

func GetFileSet(rt runtime.RunTime, setType FileSetType, sameAbove bool) (*FileSet, error) {
	db, err := GetDB(rt)
	if err != nil {
		return nil, err
	}

	return GetFileSetTx(rt, db, setType, sameAbove)
}

func GetFileSetTx(rt runtime.RunTime, db *gorm.DB, setType FileSetType, sameAbove bool) (*FileSet, error) {
	c, err := systeminit.GetInit()
	if err != nil {
		return nil, err
	}

	db, err = GetDB(rt)
	if err != nil {
		return nil, err
	}

	var res FileSet
	err = db.Model(&FileSet{}).Where("filesettype = ?", setType).Order("filesetid desc").First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fs := newFileSet(rt, setType, 1)
		err := db.Create(fs).Error
		if err != nil {
			return nil, err
		}
		return fs, nil
	} else if err != nil {
		return nil, err
	}

	if res.FileSetID <= 0 {
		return nil, fmt.Errorf("file set error")
	} else if !sameAbove && res.PageCount >= c.Yaml.FileSet.MaxFilePage {
		fs := newFileSet(rt, setType, res.FileSetID+1)
		err := db.Create(fs).Error
		if err != nil {
			return nil, err
		}
		return fs, nil
	}

	return &res, err
}

func newFileSet(rt runtime.RunTime, setType FileSetType, fileSetID int64) *FileSet {
	if fileSetID <= 0 {
		fileSetID = 1
	}

	return &FileSet{
		FileSetID:   fileSetID,
		FileSetType: setType,
		PageCount:   0,
	}
}
