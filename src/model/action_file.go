package model

import (
	"errors"
	"fmt"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"gorm.io/gorm"
)

const FileOrder = "filesetid desc, fileunionid desc, fileid desc, filegroupid asc, id desc"

var FileMoveRecordNotFound = fmt.Errorf("file move record not found")

func likeValue(t string) string {
	return fmt.Sprintf("%%%s%%", t)
}

func SetWhere(tx *gorm.DB, s *SearchWhere, fst FileSetType) *gorm.DB {
	if s == nil {
		return tx
	}

	tx = tx.Where("filesettype = ?", fst)

	if s.Name != "" {
		tx = tx.Where("name LIKE ?", likeValue(s.Name))
	}

	if s.OldName != "" {
		tx = tx.Where("oldname LIKE ?", likeValue(s.OldName))
	}

	if s.IDCard != "" {
		tx = tx.Where("idcard LIKE ?", likeValue(s.IDCard))
	}

	if s.IsMan == "男性" {
		tx = tx.Where("isman == 1")
	} else if s.IsMan == "女性" {
		tx = tx.Where("isman == 0")
	}

	if s.BirthdayStart.Valid {
		tx = tx.Where("birthday >= ?", s.BirthdayStart.Time)
	}

	if s.BirthdayEnd.Valid {
		tx = tx.Where("birthday <= ?", s.BirthdayEnd.Time)
	}

	if s.Comment != "" {
		tx = tx.Where("comment LIKE ?", likeValue(s.Comment))
	}

	if s.FileSetID != 0 {
		tx = tx.Where("filesetid = ?", s.FileSetID)
	}

	if s.FileUnionID != 0 {
		tx = tx.Where("fileunionid = ?", s.FileUnionID)
	}

	if s.FileID != 0 {
		tx = tx.Where("fileid = ?", s.FileID)
	}

	if s.FileGroupID != 0 {
		tx = tx.Where("filegroupid = ?", s.FileGroupID)
	}

	return tx
}

func GetAllFile(rt runtime.RunTime, fst FileSetType, s *SearchWhere, res interface{}) error {
	db, err := GetDB(rt)
	if err != nil {
		return err
	}

	err = SetWhere(db.Model(&FileAbs{}), s, fst).Order(FileOrder).Find(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	} else if err != nil {
		return err
	}

	return nil
}

func GetNewFileID(rt runtime.RunTime, fst FileSetType, sameAbove bool) (fs *FileSet, lst File, unionID int64, fileID int64, groupID int64, err error) {
	db, err := GetDB(rt)
	if err != nil {
		return nil, nil, 0, 0, 0, err
	}

	fs, err = GetFileSet(rt, fst, sameAbove)
	if err != nil {
		return nil, nil, 0, 0, 0, err
	}

	var last FileAbs
	err = db.Model(&FileAbs{}).Where("filesettype = ?", fst).Order("fileid desc, fileunionid desc, filegroupid desc").First(&last).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, 0, 0, 0, fmt.Errorf("not found")
	} else if err != nil {
		return nil, nil, 0, 0, 0, err
	}

	if sameAbove {
		unionID = last.FileUnionID
		groupID = last.FileGroupID + 1
		fileID = last.FileID + 1
	} else {
		unionID = last.FileUnionID + 1
		groupID = 1
		fileID = last.FileID + 1
	}

	return fs, &last, unionID, fileID, groupID, nil
}

func CountAllFile(rt runtime.RunTime, fst FileSetType, s *SearchWhere) (int64, error) {
	db, err := GetDB(rt)
	if err != nil {
		return 0, err
	}

	var res PageSizeResult
	err = SetWhere(db.Model(&FileAbs{}), s, fst).Select("COUNT(*) AS ps").First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res.AllCount = 0
	} else if err != nil {
		return 0, err
	}

	return res.AllCount, err
}

func GetPageMax(rt runtime.RunTime, fileSetType FileSetType, pageItemCount int, s *SearchWhere) (int64, error) {
	allCount, err := CountAllFile(rt, fileSetType, s)
	if err != nil {
		return 0, err
	}

	if allCount%int64(pageItemCount) != 0 {
		return allCount/int64(pageItemCount) + 1, nil
	}

	return allCount / int64(pageItemCount), nil
}

func PageChoiceOffset(rt runtime.RunTime, fileSetType FileSetType, pageItemCount int, page int64, s *SearchWhere) (int64, int, int64, error) {
	pageMax, err := GetPageMax(rt, fileSetType, pageItemCount, s)
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

func GetPageData(rt runtime.RunTime, fst FileSetType, pageItemCount int, page int64, s *SearchWhere, res interface{}) (int64, error) {
	if pageItemCount <= 0 {
		pageItemCount = DefaultPageItemCount
	}

	db, err := GetDB(rt)
	if err != nil {
		return 0, err
	}

	offset, limit, pageMax, err := PageChoiceOffset(rt, fst, pageItemCount, page, s)
	err = SetWhere(db.Model(FileAbs{}), s, fst).Limit(limit).Offset(int(offset)).Order(FileOrder).Find(res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return pageMax, nil
}

func DeleteFile(rt runtime.RunTime, f File) error {
	db, err := GetDB(rt)
	if err != nil {
		return err
	}

	return db.Delete(f).Error
}

func CreateFile(rt runtime.RunTime, fst FileSetType, fc File) error {
	db, err := GetDB(rt)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		fs, lastf, unionid, fildid, groupid, err := GetNewFileID(rt, fst, fc.GetFile().SameAsAbove)
		if err != nil {
			return err
		}

		nf := fc.GetFile()
		lf := lastf.GetFile()
		nf.FileUnionID = unionid
		nf.FileID = fildid
		nf.FileGroupID = groupid
		nf.FileSetID = fs.FileSetID
		nf.FileSetType = fs.FileSetType
		nf.FileSetSQLID = int64(fs.ID)
		nf.PageStart = fs.PageCount + 1
		nf.PageEnd = fs.PageCount + nf.PageCount

		fs.PageCount += nf.PageEnd

		if nf.SameAsAbove {
			lf.PeopleCount += 1
			nf.PeopleCount = lf.PeopleCount
		} else {
			nf.PeopleCount = 1
		}

		err = tx.Create(nf).Error
		if err != nil {
			return err
		}

		err = tx.Update("peoplecount", lf.PeopleCount).Error
		if err != nil {
			return err
		}

		err = tx.Save(fs).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func SaveFile(rt runtime.RunTime, f File) error {
	db, err := GetDB(rt)
	if err != nil {
		return err
	}

	return db.Save(f).Error
}

func FindFile(rt runtime.RunTime, fileID int64, res interface{}) error {
	db, err := GetDB(rt)
	if err != nil {
		return err
	}

	return db.Where("fileid = ?", fileID).Order("fileid desc, id desc").First(res).Error
}
