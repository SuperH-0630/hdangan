package model

import (
	"errors"
	"fmt"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"github.com/SuperH-0630/hdangan/src/systeminit"
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

func CreateFile(rt runtime.RunTime, fst FileSetType, fc File, record *FileMoveRecord) error {
	config, err := systeminit.GetInit()
	if err != nil {
		return err
	}

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

		fs.PageCount += nf.PageEnd

		if nf.SameAsAbove {
			lf.PeopleCount += 1
			nf.PeopleCount = lf.PeopleCount

			nf.Time = lf.Time
			nf.PageCount = lf.PageCount
			nf.PageStart = lf.PageStart
			nf.PageEnd = lf.PageEnd
			nf.PeopleCount = lf.PeopleCount

			nf.LastMoveRecordID = lf.LastMoveRecordID

			switch ff := fc.(type) {
			case *FileQianRu:
				lff, ok := lastf.(*FileQianRu)
				if !ok {
					return fmt.Errorf("file set type error")
				}

				ff.Type = lff.Type
				ff.OldLocation = lff.OldLocation
				ff.NewLocation = lff.NewLocation
			case *FileChuSheng:
				lff, ok := lastf.(*FileChuSheng)
				if !ok {
					return fmt.Errorf("file set type error")
				}

				ff.Type = lff.Type
				ff.NewLocation = lff.NewLocation
			case *FileQianChu:
				lff, ok := lastf.(*FileQianChu)
				if !ok {
					return fmt.Errorf("file set type error")
				}

				ff.Type = lff.Type
				ff.NewLocation = lff.NewLocation
			case *FileSiWang:
				lff, ok := lastf.(*FileSiWang)
				if !ok {
					return fmt.Errorf("file set type error")
				}

				ff.Type = lff.Type
				ff.Location = lff.Location
			case *FileBianGeng:
				lff, ok := lastf.(*FileBianGeng)
				if !ok {
					return fmt.Errorf("file set type error")
				}

				ff.Type = lff.Type
				ff.Location = lff.Location
			case *FileSuoNeiYiJu:
				lff, ok := lastf.(*FileSuoNeiYiJu)
				if !ok {
					return fmt.Errorf("file set type error")
				}

				ff.Type = lff.Type
				ff.Location = lff.Location
			case *FileSuoJianYiJu:
				lff, ok := lastf.(*FileSuoJianYiJu)
				if !ok {
					return fmt.Errorf("file set type error")
				}

				ff.Type = lff.Type
				ff.Location = lff.Location
			case *FileNongZiZhuanFei:
				lff, ok := lastf.(*FileNongZiZhuanFei)
				if !ok {
					return fmt.Errorf("file set type error")
				}

				ff.Type = lff.Type
				ff.Location = lff.Location
			case *FileYiZhanShiQianYiZheng:
				lff, ok := lastf.(*FileYiZhanShiQianYiZheng)
				if !ok {
					return fmt.Errorf("file set type error")
				}

				ff.Type = lff.Type
				ff.Location = lff.Location
			}

		} else {
			nf.PeopleCount = 1
			nf.PageStart = fs.PageCount + 1
			nf.PageEnd = fs.PageCount + nf.PageCount
		}

		err = tx.Create(fc).Error
		if err != nil {
			return err
		}

		modelMaker, ok := FileSetTypeMaker[fst]
		if !ok {
			return fmt.Errorf("bad file set type")
		}
		err = tx.Model(modelMaker()).Update("peoplecount", lf.PeopleCount).Where("fileunionid = ?", lf.FileUnionID).Error
		if err != nil {
			return err
		}

		err = tx.Save(fs).Error
		if err != nil {
			return err
		}

		if record != nil && !nf.SameAsAbove {
			record.MoveStatus = config.Yaml.Move.MoveInStatus
			record.MoveTime = nf.Time

			if !record.MoveInPeopleName.Valid || len(record.MoveInPeopleName.String) == 0 {
				record.MoveInPeopleName.Valid = true
				record.MoveInPeopleName.String = config.Yaml.Move.MoveInPeopleDefault
			}

			if !record.MoveInPeopleUnit.Valid || len(record.MoveInPeopleUnit.String) == 0 {
				record.MoveInPeopleUnit.Valid = true
				record.MoveInPeopleUnit.String = config.Yaml.Move.MoveInUnitDefault
			}

			err := createFileRecord(rt, tx, fc, record)
			if err != nil {
				return err
			}
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
