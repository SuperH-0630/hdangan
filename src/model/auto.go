package model

import (
	"github.com/SuperH-0630/hdangan/src/runtime"
)

func AutoCreateModel(rt runtime.RunTime) error {
	db, err := GetDB(rt)
	if err != nil {
		return err
	}

	return db.AutoMigrate(&FileSet{}, &FileQianRu{}, &FileChuSheng{}, &FileQianChu{}, &FileSiWang{},
		&FileBianGeng{}, &FileSuoNeiYiJu{}, &FileSuoJianYiJu{}, &FileNongZiZhuanFei{},
		&FileYiZhanShiQianYiZheng{}, &FileMoveRecord{})
}
