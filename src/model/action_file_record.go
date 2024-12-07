package model

import (
	"database/sql"
	"errors"
	"fmt"
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

	if s.MoveInPeopleName != "" {
		tx = tx.Where("moveinpeoplename LIKE ?", likeValue(s.MoveInPeopleName))
	}

	if s.MoveInPeopleUnit != "" {
		tx = tx.Where("moveinpeopleunit LIKE ?", likeValue(s.MoveInPeopleUnit))
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

func GetPageDataRecord(rt runtime.RunTime, f File, pageItemCount int, page int64, s *SearchRecord) ([]FileMoveRecord, int64, error) {
	if pageItemCount <= 0 {
		pageItemCount = DefaultPageItemCount
	}

	file := f.GetFile()

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

func GetAllRecord(rt runtime.RunTime, fc File, s *SearchRecord) ([]FileMoveRecord, error) {
	db, err := GetDB(rt)
	if err != nil {
		return nil, err
	}

	var res []FileMoveRecord
	tx := SetRecord(db.Model(&FileMoveRecord{}), s)
	if fc != nil {
		tx = tx.Where("fileunionid = ?", fc.GetFile().FileUnionID)
	}
	err = tx.Order("movetime desc, id desc").Find(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []FileMoveRecord{}, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func CheckFileMoveOut(rt runtime.RunTime, fc File) (*FileMoveRecord, bool) {
	f := fc.GetFile()

	if !f.LastMoveRecordID.Valid {
		return nil, false
	}

	record, err := FindMoveRecord(rt, fc)
	if err != nil {
		return nil, false
	}

	return record, true
}

func FindMoveRecord(rt runtime.RunTime, fc File) (*FileMoveRecord, error) {
	db, err := GetDB(rt)
	if err != nil {
		return nil, err
	}

	f := fc.GetFile()

	if !f.LastMoveRecordID.Valid || f.LastMoveRecordID.Int64 <= 0 {
		return nil, FileMoveRecordNotFound
	}

	var res FileMoveRecord
	err = db.Model(&FileMoveRecord{}).Where("id = ?", f.LastMoveRecordID).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, FileMoveRecordNotFound
	} else if err != nil {
		return nil, err
	}

	return &res, nil
}

func FindRecord(rt runtime.RunTime, recordID int64) (*FileMoveRecord, error) {
	db, err := GetDB(rt)
	if err != nil {
		return nil, err
	}

	var res FileMoveRecord
	err = db.Model(&FileMoveRecord{}).Where("id = ?", recordID).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, FileMoveRecordNotFound
	} else if err != nil {
		return nil, err
	}

	return &res, nil
}

func SaveRecord(rt runtime.RunTime, r *FileMoveRecord) error {
	db, err := GetDB(rt)
	if err != nil {
		return err
	}

	return db.Save(r).Error
}

func CreateFileRecord(rt runtime.RunTime, fc File, record *FileMoveRecord) error {
	db, err := GetDB(rt)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		return createFileRecord(rt, tx, fc, record)
	})
}

func createFileRecord(rt runtime.RunTime, tx *gorm.DB, fc File, record *FileMoveRecord) error {
	f := fc.GetFile()
	if f.ID <= 0 {
		return fmt.Errorf("file not save")
	}

	var file1 FileAbs
	var oldRecord *FileMoveRecord

	oldRecord, err := FindMoveRecord(rt, fc)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		oldRecord = nil
	} else {
		return err
	}

	err = tx.Where("fileunionid = ?", f.FileUnionID).Order("filegroupid asc").First(&file1).Error
	err = tx.Save(record).Error
	if err != nil {
		return err
	}

	record.FileSetSQLID = file1.FileSetSQLID
	record.FileSetID = file1.FileSetID
	record.FileSetType = file1.FileSetType

	record.FileSQLID = int64(file1.ID)
	record.FileUnionID = int64(file1.ID)

	if oldRecord != nil {
		record.UpRecord = sql.NullInt64{Valid: true, Int64: int64(oldRecord.ID)}
	}
	err = tx.Save(record).Error
	if err != nil {
		return err
	}

	err = tx.Update("lastmoverecordid", sql.NullInt64{Valid: true, Int64: int64(record.ID)}).Where("fileunionid = ?", file1.FileUnionID).Error
	if err != nil {
		return err
	}

	if oldRecord != nil {
		oldRecord.NextRecord = sql.NullInt64{Valid: true, Int64: int64(record.ID)}
		err = tx.Save(oldRecord).Error
		if err != nil {
			return err
		}
	}
	return nil
}
