package model

import (
	"errors"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"gorm.io/gorm"
)

func SetRecord(tx *gorm.DB, s *SearchRecord) *gorm.DB {
	if s == nil {
		return tx
	}

	if s.MoveOutStart.Valid {
		tx = tx.Where("movetime >= ?", s.MoveOutStart)
	}

	if s.MoveOutEnd.Valid {
		tx = tx.Where("movetime <= ?", s.MoveOutEnd)
	}

	if s.MoveStatus != "" {
		tx = tx.Where("movestatus LIKE ?", likeValue(s.MoveStatus))
	}

	if s.MoveOutPeopleName != "" {
		tx = tx.Where("moveoutpeoplename LIKE ?", likeValue(s.MoveOutPeopleName))
	}

	if s.MoveOutPeopleUnit != "" {
		tx = tx.Where("moveoutpeopleunit LIKE ?", likeValue(s.MoveOutPeopleUnit))
	}

	return tx
}

func CountAllFileRecord(rt runtime.RunTime, fid int64, s *SearchRecord) (int64, error) {
	db, err := GetDB(rt)
	if err != nil {
		return 0, err
	}

	var res PageSizeResult
	err = SetRecord(db.Model(&FileMoveRecord{}), s).Select("COUNT(*) AS ps").Where("filesqlid = ?", fid).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res.AllCount = 0
	} else if err != nil {
		return 0, err
	}

	return res.AllCount, err
}

func GetPageMaxRecord(rt runtime.RunTime, fid int64, pageItemCount int, s *SearchRecord) (int64, error) {
	allCount, err := CountAllFileRecord(rt, fid, s)
	if err != nil {
		return 0, err
	}

	if allCount%int64(pageItemCount) != 0 {
		return allCount/int64(pageItemCount) + 1, nil
	}

	return allCount / int64(pageItemCount), nil
}

func PageChoiceOffsetRecord(rt runtime.RunTime, fid int64, pageItemCount int, page int64, s *SearchRecord) (int64, int, int64, error) {
	pageMax, err := GetPageMaxRecord(rt, fid, pageItemCount, s)
	if err != nil {
		return 0, 0, 0, err
	}

	if page > pageMax {
		page = pageMax
	} else if page <= 0 {
		page = 1
	}

	return (page - 1) * int64(pageItemCount), pageItemCount, pageMax, nil
}

func GetPageDataRecord(rt runtime.RunTime, file *File, pageItemCount int, page int64, s *SearchRecord) ([]FileMoveRecord, int64, error) {
	if pageItemCount <= 0 {
		pageItemCount = DefaultPageItemCount
	}

	db, err := GetDB(rt)
	if err != nil {
		return []FileMoveRecord{}, 0, err
	}

	offset, limit, pageMax, err := PageChoiceOffsetRecord(rt, int64(file.ID), pageItemCount, page, s)
	res := make([]FileMoveRecord, 0, pageItemCount)

	err = SetRecord(db.Model(FileMoveRecord{}), s).Preload("File").Where("filesqlid = ?", file.ID).Limit(limit).Offset(int(offset)).Order("movetime desc, id desc").Find(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return make([]FileMoveRecord, 0, pageItemCount), 0, nil
	} else if err != nil {
		return nil, 0, err
	}

	return res, pageMax, nil
}

func FindFileRecord(recordID int64) *FileMoveRecord {
	var res FileMoveRecord
	err := db.Model(FileMoveRecord{}).Preload("File").Select("id = ?", recordID).Order("movetime desc, id desc").First(&res).Error
	if err != nil {
		return nil
	}
	return &res
}

func GetAllFileRecord(rt runtime.RunTime, f *File, s *SearchRecord) ([]FileMoveRecord, error) {
	db, err := GetDB(rt)
	if err != nil {
		return nil, err
	}

	var res []FileMoveRecord
	tx := SetRecord(db.Model(&FileMoveRecord{}), s).Preload("File")
	if f != nil {
		tx = tx.Where("filesqlid = ?", f.ID)
	}
	err = tx.Order("movetime desc, id desc").Find(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []FileMoveRecord{}, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func CheckFileMoveOut(f *File) (*FileMoveRecord, bool) {
	if !f.LastMoveRecordID.Valid {
		return nil, false
	}

	record, err := FindMoveRecord(f)
	if err != nil {
		return nil, false
	}

	return record, true
}
