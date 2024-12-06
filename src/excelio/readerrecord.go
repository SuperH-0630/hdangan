package excelio

import (
	"database/sql"
	"fmt"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"github.com/xuri/excelize/v2"
	"io"
	"strconv"
	"strings"
	"time"
)

var RecordInputTitle = []string{
	"档案ID", "出借记录ID", "操作目的", "状态", "借出时间", "借出人", "借出单位", "借入人", "借入单位", "备注",
}

func ReadRecord(rt runtime.RunTime, fst model.FileSetType, reader io.ReadCloser) (int64, int64, int64, error) {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return 0, 0, 0, err
	}

	tit := RecordInputTitle

	defer func() {
		_ = f.Close()
	}()

	slst := f.GetSheetList()
	if len(slst) != 1 {
		return 0, 0, 0, fmt.Errorf("sheet 过多")
	}

	sheet := slst[0]

	timeStyleID, err := f.NewStyle(&excelize.Style{})
	if err != nil {
		return 0, 0, 0, err
	}

	for i, t := range tit {
		if strings.Contains(t, "日期") || strings.Contains(t, "时间") {
			err = f.SetColStyle(sheet, Header[i], timeStyleID)
			if err != nil {
				return 0, 0, 0, err
			}
		}
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return 0, 0, 0, err
	}

	var sa, su, fu int64

	for i, j := range rows {
		if i == 0 {
			if !checkTitle(j, tit) {
				return 0, 0, 0, BadTitle
			}
		} else {
			t := make([]string, 0, len(tit))
			t = append(t, j...)
			if len(t) < len(tit) {
				t = append(t, make([]string, len(tit)-len(t))...)
			}

			s := makeRecord(rt, fst, t, tit)
			if s == successAdd {
				sa += 1
			} else if s == successUpdate {
				su += 1
			} else {
				fu += 1
			}
		}
	}

	return sa, su, fu, nil
}

func makeRecord(rt runtime.RunTime, fst model.FileSetType, t []string, tit []string) int {
	if t[1] != "更新" {
		var f model.File

		if len(t[0]) != 0 {
			fileID, err := strconv.ParseInt(t[0], 10, 64)
			if err != nil {
				return fail
			}

			maker, ok := model.FileSetTypeMaker[fst]
			if !ok {
				return fail
			}

			f = maker()
			err = model.FindFile(rt, fileID, f) // f 已经是指针
			if err != nil {
				return fail
			}
		} else {
			return fail
		}

		file, record, err := makeNewRecord(rt, f, fst, t, tit)
		if err != nil {
			return fail
		}

		err = model.CreateFileRecord(rt, file, record)
		if err != nil {
			return fail
		}
		return successAdd
	} else {
		var r *model.FileMoveRecord

		if len(t[1]) != 0 {
			recordID, err := strconv.ParseInt(t[1], 10, 64)
			if err != nil {
				return fail
			}

			r, err = model.FindRecord(rt, recordID)
			if err != nil {
				return fail
			}
		} else {
			return fail
		}

		r, err := makeUpdateRecord(rt, r, t, tit)
		if err != nil {
			return fail
		}

		err = model.SaveRecord(rt, r)
		if err != nil {
			return fail
		}
		return successUpdate
	}
}

func makeNewRecord(rt runtime.RunTime, f model.File, fst model.FileSetType, t []string, tit []string) (model.File, *model.FileMoveRecord, error) {
	file := f.GetFile()

	moveStatus := t[2]
	if len(moveStatus) == 0 {
		return nil, nil, fmt.Errorf("move status should be given")
	}

	moveTime := time.Now()
	if len(t[3]) != 0 {
		var err error
		moveTime, err = timeReader(time.Now(), t[3])
		if err != nil {
			return nil, nil, err
		}
	}

	moveOutPeople := sql.NullString{Valid: len(t[4]) != 0, String: t[4]}
	moveOutUnit := sql.NullString{Valid: len(t[5]) != 0, String: t[5]}
	moveInPeople := sql.NullString{Valid: len(t[6]) != 0, String: t[6]}
	moveInUnit := sql.NullString{Valid: len(t[7]) != 0, String: t[7]}
	moveComment := sql.NullString{Valid: len(t[8]) != 0, String: t[8]}

	res := &model.FileMoveRecord{
		MoveStatus:        moveStatus,
		MoveTime:          moveTime,
		MoveOutPeopleName: moveOutPeople,
		MoveOutPeopleUnit: moveOutUnit,
		MoveInPeopleName:  moveInPeople,
		MoveInPeopleUnit:  moveInUnit,
		MoveComment:       moveComment,
	}

	return file, res, nil
}

func makeUpdateRecord(rt runtime.RunTime, record *model.FileMoveRecord, t []string, tit []string) (*model.FileMoveRecord, error) {
	moveStatus := t[2]
	if len(moveStatus) != 0 {
		record.MoveStatus = moveStatus
	}

	if len(t[3]) != 0 {
		var err error
		record.MoveTime, err = timeReader(time.Now(), t[3])
		if err != nil {
			return nil, err
		}
	}

	if len(t[4]) != 0 {
		record.MoveOutPeopleName = sql.NullString{
			Valid:  true,
			String: t[4],
		}
	}

	if len(t[5]) != 0 {
		record.MoveOutPeopleUnit = sql.NullString{
			Valid:  true,
			String: t[5],
		}
	}

	if len(t[6]) != 0 {
		record.MoveInPeopleName = sql.NullString{
			Valid:  true,
			String: t[6],
		}
	}

	if len(t[7]) != 0 {
		record.MoveInPeopleUnit = sql.NullString{
			Valid:  true,
			String: t[7],
		}
	}

	if len(t[8]) != 0 {
		record.MoveComment = sql.NullString{
			Valid:  true,
			String: t[8],
		}
	}

	return record, nil
}
