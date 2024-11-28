package model

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
)

func SaveFileRecord(file *File, record *FileMoveRecord, oldRecord *FileMoveRecord) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if file != nil {
			err := tx.Save(file).Error
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("not file")
		}

		if record != nil {
			record.FileSQLID = int64(file.ID)
			if oldRecord != nil {
				record.UpRecord = sql.NullInt64{Valid: true, Int64: int64(oldRecord.ID)}
			}
			err := tx.Save(record).Error
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("not human")
		}

		file.LastMoveRecordID = sql.NullInt64{Valid: true, Int64: int64(record.ID)}
		err := tx.Save(file).Error
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
	})
}
