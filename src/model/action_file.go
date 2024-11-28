package model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"gorm.io/gorm"
)

var FileMoveRecordNotFound = fmt.Errorf("file move record not found")

func likeValue(t string) string {
	return fmt.Sprintf("%%%s%%", t)
}

func SetWhere(tx *gorm.DB, s *SearchWhere) *gorm.DB {
	if s == nil {
		return tx
	}

	if s.Name != "" {
		tx = tx.Where("name LIKE ?", likeValue(s.Name))
	}

	if s.IDCard != "" {
		tx = tx.Where("idcard LIKE ?", likeValue(s.IDCard))
	}

	if s.Location != "" {
		tx = tx.Where("loc LIKE ?", likeValue(s.Location))
	}

	if s.FileID != 0 {
		tx = tx.Where("fileid = ?", s.FileID)
	}

	if s.FileTitle != "" {
		tx = tx.Where("filetitle LIKE ?", likeValue(s.FileTitle))
	}

	if s.FileType != "" {
		tx = tx.Where("filetype LIKE ?", likeValue(s.FileType))
	}

	if s.FirstMoveInStart.Valid {
		tx = tx.Where("firstmovein >= ?", s.FirstMoveInStart)
	}

	if s.FirstMoveInEnd.Valid {
		tx = tx.Where("firstmovein <= ?", s.FirstMoveInEnd)
	}

	if s.LastMoveInStart.Valid {
		tx = tx.Where("lastmovein >= ?", s.LastMoveInStart)
	}

	if s.LastMoveInStart.Valid {
		tx = tx.Where("lastmovein <= ?", s.LastMoveInStart)
	}

	if s.LastMoveOutStart.Valid {
		tx = tx.Where("lastmoveout >= ?", s.LastMoveOutStart)
	}

	if s.LastMoveOutStart.Valid {
		tx = tx.Where("lastmoveout <= ?", s.LastMoveOutStart)
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

func GetAllFile(rt runtime.RunTime, s *SearchWhere) ([]File, error) {
	db, err := GetDB(rt)
	if err != nil {
		return nil, err
	}

	var res []File
	err = SetWhere(db.Model(&File{}), s).Order("fileid desc, id desc").Find(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []File{}, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func GetNewFileID(rt runtime.RunTime) (int64, error) {
	db, err := GetDB(rt)
	if err != nil {
		return 0, err
	}

	var res FileIDMax
	err = db.Model(&File{}).Select("MAX(fileid) AS fid").First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res.FileID = 1
	} else if err != nil {
		return 0, err
	}

	if res.FileID <= 0 {
		res.FileID = 1
	}

	return res.FileID + 1, err
}

func CountAllFile(rt runtime.RunTime, s *SearchWhere) (int64, error) {
	db, err := GetDB(rt)
	if err != nil {
		return 0, err
	}

	var res PageSizeResult
	err = SetWhere(db.Model(&File{}), s).Select("COUNT(*) AS ps").First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res.AllCount = 0
	} else if err != nil {
		return 0, err
	}

	return res.AllCount, err
}

func GetPageMax(rt runtime.RunTime, pageItemCount int, s *SearchWhere) (int64, error) {
	allCount, err := CountAllFile(rt, s)
	if err != nil {
		return 0, err
	}

	if allCount%int64(pageItemCount) != 0 {
		return allCount/int64(pageItemCount) + 1, nil
	}

	return allCount / int64(pageItemCount), nil
}

func PageChoiceOffset(rt runtime.RunTime, pageItemCount int, page int64, s *SearchWhere) (int64, int, int64, error) {
	pageMax, err := GetPageMax(rt, pageItemCount, s)
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

func GetPageData(rt runtime.RunTime, pageItemCount int, page int64, s *SearchWhere) ([]File, int64, error) {
	if pageItemCount <= 0 {
		pageItemCount = DefaultPageItemCount
	}

	db, err := GetDB(rt)
	if err != nil {
		return []File{}, 0, err
	}

	offset, limit, pageMax, err := PageChoiceOffset(rt, pageItemCount, page, s)
	res := make([]File, 0, pageItemCount)

	err = SetWhere(db.Model(File{}), s).Limit(limit).Offset(int(offset)).Order("fileid desc, id desc").Find(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return make([]File, 0, pageItemCount), 0, nil
	} else if err != nil {
		return nil, 0, err
	}

	return res, pageMax, nil
}

func DeleteFile(f *File) error {
	return db.Delete(f).Error
}

func CreateFile(f *File) error {
	return db.Create(f).Error
}

func SaveFile(f *File) error {
	return db.Save(f).Error
}

func SaveNewMoveOut(f *File, record *FileMoveRecord) error {
	err := db.Create(record).Error
	if err != nil {
		return err
	}

	f.LastMoveRecordID = sql.NullInt64{
		Valid: true,
		Int64: int64(record.ID),
	}

	return db.Save(f).Error
}

func UpdateRecord(f *File, record *FileMoveRecord) error {
	err := db.Save(record).Error
	if err != nil {
		return err
	}

	return db.Save(f).Error
}

func FindMoveRecord(f *File) (*FileMoveRecord, error) {
	if !f.LastMoveRecordID.Valid {
		return nil, FileMoveRecordNotFound
	}

	var res FileMoveRecord

	err := db.Model(&FileMoveRecord{}).Where("id = ?", f.LastMoveRecordID.Int64).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, FileMoveRecordNotFound
	} else if err != nil {
		return nil, err
	}

	return &res, nil
}

func FindFile(fileID int64) *File {
	var res File
	err := db.Model(File{}).Where("fileid = ?", fileID).Order("fileid desc, id desc").First(&res).Error
	if err != nil {
		return nil
	}
	return &res
}
