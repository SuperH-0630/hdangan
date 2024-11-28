package excelreader

import (
	"database/sql"
	"fmt"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"github.com/SuperH-0630/hdangan/src/systeminit"
	"github.com/xuri/excelize/v2"
	"io"
	"strconv"
	"time"
)

var Title = []string{
	"档案号", "姓名", "身份证", "户籍地", "卷宗标题", "卷宗类型", "卷宗备注", "第一次迁移时间", "最后一次迁移时间", "迁移状态", "迁移人", "迁移单位", "迁移备注",
}

var RecordTitle = []string{
	"档案号", "姓名", "身份证", "户籍地", "卷宗标题", "卷宗类型", "卷宗备注", "迁入迁出状态", "时间", "迁出单位", "迁出备注",
}

var Header = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N"}

func indexToDataRecord(f *model.FileMoveRecord, i int) string {
	if f == nil {
		return ""
	}

	switch i {
	case 0:
		return fmt.Sprintf("%d", f.File.FileID)
	case 1:
		return fmt.Sprintf("%s", f.File.Name)
	case 2:
		return fmt.Sprintf("%s", f.File.IDCard)
	case 3:
		return fmt.Sprintf("%s", f.File.Location)

	case 4:
		return fmt.Sprintf("%s", f.File.FileType)
	case 5:
		return fmt.Sprintf("%s", f.File.FileType)
	case 6:
		return fmt.Sprintf("%s", f.File.FileComment.String)

	case 7:
		return fmt.Sprintf("%s", f.MoveStatus)
	case 8:
		return fmt.Sprintf("%s", f.MoveTime.Format("2006-01-02 15:04:05"))

	case 9:
		return fmt.Sprintf("%s", getString(f.MoveOutPeopleName))
	case 10:
		return fmt.Sprintf("%s", getString(f.MoveOutPeopleUnit))
	case 11:
		return fmt.Sprintf("%s", getString(f.MoveComment))
	}

	return ""
}

func indexToData(f *model.File, i int) string {
	if f == nil {
		return ""
	}

	switch i {
	case 0:
		return fmt.Sprintf("%d", f.FileID)
	case 1:
		return fmt.Sprintf("%s", f.Name)
	case 2:
		return fmt.Sprintf("%s", f.IDCard)
	case 3:
		return fmt.Sprintf("%s", f.Location)

	case 4:
		return fmt.Sprintf("%s", f.FileTitle)
	case 5:
		return fmt.Sprintf("%s", f.FileType)
	case 6:
		return fmt.Sprintf("%s", f.FileComment.String)
	case 7:
		return fmt.Sprintf("%s", f.FirstMoveIn)
	case 8:
		return fmt.Sprintf("%s", f.LastMoveIn)
	case 9:
		return fmt.Sprintf("%s", f.MoveStatus)
	case 10:
		return fmt.Sprintf("%s", getString(f.MoveOutPeopleName))
	case 11:
		return fmt.Sprintf("%s", getString(f.MoveOutPeopleUnit))
	case 12:
		return fmt.Sprintf("%s", getString(f.MoveComment))
	}

	return ""
}

func getString(str sql.NullString) string {
	if str.Valid {
		return str.String
	}
	return ""
}

func getTime(t sql.NullTime) string {
	if t.Valid {
		return t.Time.Format("2006-01-02 15:04:05")
	}
	return ""
}

var BadTitle = fmt.Errorf("表格首行的标题错误")

const (
	successAdd = iota
	successUpdate
	fail
)

func CreateTemplate(rt runtime.RunTime, savepath string) error {
	var err error

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	sheetIndex := 0
	sheetName := "sheet1"
	slts := f.GetSheetList()
	if len(slts) == 0 {
		sheetIndex, err = f.NewSheet(sheetName)
		if err != nil {
			return err
		}
		f.SetActiveSheet(sheetIndex)
	}

	f.SetActiveSheet(sheetIndex)

	for i, k := range Title {
		err = f.SetCellStr(sheetName, fmt.Sprintf("%s1", Header[i]), k)
		if err != nil {
			return err
		}
	}

	styleId, err := f.NewStyle(&excelize.Style{})
	if err != nil {
		return err
	}

	_ = f.SetColStyle(sheetName, Header[7], styleId)
	_ = f.SetColStyle(sheetName, Header[8], styleId)
	_ = f.SetColStyle(sheetName, Header[10], styleId)

	err = f.SaveAs(savepath)
	if err != nil {
		return err
	}

	return nil
}

func OutputFile(rt runtime.RunTime, savepath string, files []model.File, s *model.SearchWhere) error {
	var err error

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	sheetIndex := 0
	sheetName := "Sheet1"
	slts := f.GetSheetList()
	if len(slts) == 0 {
		sheetIndex, err = f.NewSheet(sheetName)
		if err != nil {
			return err
		}
	} else {
		sheetName = slts[0]
	}

	f.SetActiveSheet(sheetIndex)

	styleId, err := f.NewStyle(&excelize.Style{})
	if err != nil {
		return err
	}

	_ = f.SetColStyle(sheetName, Header[7], styleId)
	_ = f.SetColStyle(sheetName, Header[8], styleId)
	_ = f.SetColStyle(sheetName, Header[10], styleId)

	if files == nil || len(files) == 0 {
		files, err = model.GetAllFile(rt, s)
		if err != nil {
			return err
		}
	}

	for i, k := range Title {
		err = f.SetCellStr(sheetName, fmt.Sprintf("%s1", Header[i]), k)
		if err != nil {
			return err
		}
	}

	for h, _ := range Title {
		for j, y := range files {
			err = f.SetCellStr(sheetName, fmt.Sprintf("%s%d", Header[h], j+2), indexToData(&y, h))
			if err != nil {
				return err
			}
		}
	}

	err = f.SaveAs(savepath)
	if err != nil {
		return err
	}

	return nil
}

func OutputFileRecord(rt runtime.RunTime, savepath string, file *model.File, record []model.FileMoveRecord, s *model.SearchRecord) error {
	var err error

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	sheetIndex := 0
	sheetName := "Sheet1"
	slts := f.GetSheetList()
	if len(slts) == 0 {
		sheetIndex, err = f.NewSheet(sheetName)
		if err != nil {
			return err
		}
	} else {
		sheetName = slts[0]
	}

	f.SetActiveSheet(sheetIndex)

	styleId, err := f.NewStyle(&excelize.Style{})
	if err != nil {
		return err
	}

	_ = f.SetColStyle(sheetName, Header[7], styleId)
	_ = f.SetColStyle(sheetName, Header[8], styleId)

	if record == nil || len(record) == 0 {
		record, err = model.GetAllFileRecord(rt, file, s)
		if err != nil {
			return err
		}
	}

	for i, k := range RecordTitle {
		err = f.SetCellStr(sheetName, fmt.Sprintf("%s1", Header[i]), k)
		if err != nil {
			return err
		}
	}

	for h, _ := range RecordTitle {
		for j, y := range record {
			err = f.SetCellStr(sheetName, fmt.Sprintf("%s%d", Header[h], j+2), indexToDataRecord(&y, h))
			if err != nil {
				return err
			}
		}
	}

	err = f.SaveAs(savepath)
	if err != nil {
		return err
	}

	return nil
}

func ReadFile(rt runtime.RunTime, reader io.ReadCloser) (int64, int64, int64, error) {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return 0, 0, 0, err
	}

	defer func() {
		_ = f.Close()
	}()

	slst := f.GetSheetList()
	if len(slst) != 1 {
		return 0, 0, 0, fmt.Errorf("sheet 过多")
	}

	sheet := slst[0]

	styleId, err := f.NewStyle(&excelize.Style{})
	if err != nil {
		return 0, 0, 0, err
	}

	err = f.SetColStyle(sheet, Header[7], styleId)
	if err != nil {
		return 0, 0, 0, err
	}

	err = f.SetColStyle(sheet, Header[8], styleId)
	if err != nil {
		return 0, 0, 0, err
	}

	err = f.SetColStyle(sheet, Header[10], styleId)
	if err != nil {
		return 0, 0, 0, err
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return 0, 0, 0, err
	}

	var sa, su, fu int64

	for i, j := range rows {
		if i == 0 {
			if !checkTitle(j) {
				return 0, 0, 0, BadTitle
			}
		} else {
			t := make([]string, 0, len(Title))
			t = append(t, j...)
			t = append(t, make([]string, len(Title)-len(t))...)
			s := makeFile(rt, t)
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

func checkTitle(t []string) bool {
	for i, s := range t {
		if Title[i] != s {
			return false
		}
	}
	return true
}

func makeFile(rt runtime.RunTime, t []string) int {
	file, record, oldRecord, err := _makeFile(rt, t)
	if err != nil {
		fmt.Printf("err: %s", err)
	}

	if err == nil && file != nil && record != nil {
		isNew := file.ID == 0
		err := model.SaveFileRecord(file, record, oldRecord)
		if err == nil {
			if isNew {
				return successAdd
			} else {
				return successUpdate
			}
		}
	}
	return fail
}

func _makeFile(rt runtime.RunTime, t []string) (*model.File, *model.FileMoveRecord, *model.FileMoveRecord, error) {
	if t[0] == "" {
		fileID, err := model.GetNewFileID(rt)
		if err != nil {
			return nil, nil, nil, err
		}
		return makeNewFile(rt, fileID, t)
	} else {
		tmp, err := strconv.ParseInt(t[0], 10, 64)
		if err != nil {
		}

		if err != nil || tmp <= 0 {
			fileID, err := model.GetNewFileID(rt)
			if err != nil {
				return nil, nil, nil, err
			}
			return makeNewFile(rt, fileID, t)
		} else {
			f := model.FindFile(tmp)
			if f == nil {
				return makeNewFile(rt, tmp, t)
			} else {
				return makeUpdateFile(rt, f, t)
			}
		}
	}
}

func makeNewFile(rt runtime.RunTime, fileID int64, t []string) (*model.File, *model.FileMoveRecord, *model.FileMoveRecord, error) {
	c, err := systeminit.GetInit()
	if err != nil {
		return nil, nil, nil, err
	}

	name := t[1]
	if len(name) == 0 {
		return nil, nil, nil, fmt.Errorf("must has name")
	}

	idcard := t[2]
	if len(idcard) == 0 {
		return nil, nil, nil, fmt.Errorf("must has id")
	} else if idcard[0] == '\'' {
		idcard = idcard[1:]
	}

	loc := t[3]
	if len(loc) == 0 {
		return nil, nil, nil, fmt.Errorf("must has loc")
	}

	fileTitle := t[4]
	if len(fileTitle) == 0 {
		return nil, nil, nil, fmt.Errorf("must has title")
	}

	fileType := t[5]
	if len(fileType) == 0 {
		return nil, nil, nil, fmt.Errorf("must has filetype")
	}

	fileComment := t[6]

	firstMoveInTime := time.Now()
	if len(t[7]) != 0 {
		firstMoveInTime, err = timeReader(time.Now(), t[7])
		if err != nil {
			return nil, nil, nil, err
		}
	}

	lastMoveInTime := firstMoveInTime // t[8]忽略

	moveStatus := t[9]
	isMoveIn := moveStatus == c.Yaml.Move.MoveInStatus
	if len(moveStatus) == 0 {
		return nil, nil, nil, fmt.Errorf("must has move status")
	}

	lastMoveTime := sql.NullTime{
		Valid: false,
	}

	lasttMove := t[8]
	if isMoveIn {
		if len(lasttMove) != 0 {
			tmp, err := timeReader(firstMoveInTime, lasttMove)
			if err == nil {
				lastMoveTime.Valid = true
				lastMoveTime.Time = tmp
			} else {
				return nil, nil, nil, err
			}
		} else {
			lastMoveTime.Valid = true
			lastMoveTime.Time = firstMoveInTime
		}
	} else {
		if len(lasttMove) != 0 {
			tmp, err := timeReader(time.Now(), lasttMove)
			if err == nil {
				lastMoveTime.Valid = true
				lastMoveTime.Time = tmp
			} else {
				return nil, nil, nil, err
			}
		} else {
			lastMoveTime.Valid = true
			lastMoveTime.Time = time.Now()
		}
	}

	movePeople := t[10]
	if isMoveIn {
		movePeople = ""
	} else {
		if len(movePeople) == 0 {
			return nil, nil, nil, fmt.Errorf("must has move people")
		}
	}

	moveUnit := t[11]
	if isMoveIn {
		moveUnit = ""
	} else {
		if len(moveUnit) == 0 {
			return nil, nil, nil, fmt.Errorf("must has move unit")
		}
	}

	moveComment := t[12]
	if isMoveIn && len(moveComment) != 0 {
		moveComment = ""
	}

	file := model.File{
		FileID:    fileID,
		Name:      name,
		IDCard:    idcard,
		Location:  loc,
		FileTitle: fileTitle,
		FileType:  fileType,
		FileComment: sql.NullString{
			Valid:  len(fileComment) != 0,
			String: fileComment,
		},

		FirstMoveIn: firstMoveInTime,
		LastMoveIn:  lastMoveInTime,

		MoveStatus: moveStatus,
		MoveOutPeopleName: sql.NullString{
			Valid:  !isMoveIn,
			String: movePeople,
		},
		MoveOutPeopleUnit: sql.NullString{
			Valid:  !isMoveIn,
			String: moveUnit,
		},
		MoveComment: sql.NullString{
			Valid:  len(moveComment) != 0 && !isMoveIn,
			String: moveComment,
		},
	}

	record := model.FileMoveRecord{
		// 不设置File参数
		MoveStatus: moveStatus,
		MoveTime:   file.LastMoveIn,
		MoveOutPeopleName: sql.NullString{
			Valid:  !isMoveIn,
			String: movePeople,
		},
		MoveOutPeopleUnit: sql.NullString{
			Valid:  !isMoveIn,
			String: moveUnit,
		},
		MoveComment: sql.NullString{
			Valid:  len(moveComment) != 0 && !isMoveIn,
			String: moveComment,
		},
	}

	return &file, &record, nil, nil
}

func makeUpdateFile(rt runtime.RunTime, file *model.File, t []string) (*model.File, *model.FileMoveRecord, *model.FileMoveRecord, error) {
	c, err := systeminit.GetInit()
	if err != nil {
		return nil, nil, nil, err
	}

	name := t[1]
	if len(name) != 0 {
		file.Name = name
	}

	idcard := t[2]
	if len(idcard) > 0 && idcard[0] == '\'' {
		idcard = idcard[1:]
	}
	if len(idcard) != 0 {
		file.IDCard = idcard
	}

	loc := t[3]
	if len(loc) != 0 {
		file.Location = loc
	}

	fileTitle := t[4]
	if len(fileTitle) != 0 {
		file.FileTitle = fileTitle
	}

	fileType := t[5]
	if len(fileType) != 0 {
		file.FileType = fileType
	}

	fileComment := t[6]
	if len(fileComment) != 0 {
		file.FileComment = sql.NullString{
			Valid:  true,
			String: fileComment,
		}
	}

	firstMoveInTime := time.Now()
	if len(t[7]) != 0 {
		firstMoveInTime, err = timeReader(time.Now(), t[7])
		if err != nil {
			return nil, nil, nil, err
		}
	}

	file.FirstMoveIn = firstMoveInTime

	lastMoveTime := time.Time{}

	moveStatus := t[9]
	isMoveIn := moveStatus == c.Yaml.Move.MoveInStatus
	if len(moveStatus) != 0 {
		file.MoveStatus = moveStatus
	}

	lastMove := t[8]
	if isMoveIn {
		if len(lastMove) != 0 {
			tmp, err := timeReader(file.LastMoveIn, lastMove)
			if err == nil {
				lastMoveTime = tmp
			} else {
				return nil, nil, nil, err
			}
		} else {
			lastMoveTime = file.LastMoveIn
		}
	} else {
		if len(lastMove) != 0 {
			tmp, err := timeReader(time.Now(), lastMove)
			if err == nil {
				lastMoveTime = tmp
			} else {
				return nil, nil, nil, err
			}
		} else {
			lastMoveTime = time.Now()
		}
	}

	file.LastMoveIn = lastMoveTime

	movePeople := t[10]
	file.MoveOutPeopleName = sql.NullString{
		Valid:  true,
		String: movePeople,
	}

	moveUnit := t[11]
	file.MoveOutPeopleUnit = sql.NullString{
		Valid:  true,
		String: moveUnit,
	}

	moveComment := t[12]
	file.MoveComment = sql.NullString{
		Valid:  len(moveComment) != 0,
		String: moveComment,
	}

	if moveStatus == c.Yaml.Move.MoveInStatus {
		file.MoveOutPeopleName.Valid = false
		file.MoveOutPeopleUnit.Valid = false
		file.MoveComment.Valid = false
	} else {
		if !file.MoveOutPeopleName.Valid || !file.MoveOutPeopleUnit.Valid {
			return nil, nil, nil, fmt.Errorf("迁出状态需要填写迁出人和单位")
		}
	}

	if moveStatus == c.Yaml.Move.MoveInStatus {
		file.MoveOutPeopleName.Valid = false
		file.MoveOutPeopleUnit.Valid = false
		file.MoveComment.Valid = false

		record := &model.FileMoveRecord{
			// FileID在Save时会自动增加
			MoveStatus:        file.MoveStatus,
			MoveTime:          file.LastMoveIn,
			MoveOutPeopleName: file.MoveOutPeopleName,
			MoveOutPeopleUnit: file.MoveOutPeopleUnit,
			MoveComment:       file.MoveComment,
		}

		oldRecord, _ := model.CheckFileMoveOut(file)

		return file, record, oldRecord, nil
	} else {
		if !file.MoveOutPeopleName.Valid || !file.MoveOutPeopleUnit.Valid {
			return nil, nil, nil, fmt.Errorf("迁出状态需要填写迁出人和单位")
		} else {
			record := &model.FileMoveRecord{
				// FileID在Save时会自动增加
				MoveStatus:        file.MoveStatus,
				MoveTime:          file.LastMoveIn,
				MoveOutPeopleName: file.MoveOutPeopleName,
				MoveOutPeopleUnit: file.MoveOutPeopleUnit,
				MoveComment:       file.MoveComment,
			}

			oldRecord, _ := model.CheckFileMoveOut(file)

			return file, record, oldRecord, nil
		}
	}
}

func timeReader(baseTime time.Time, data string) (time.Time, error) {
	if len(data) == 0 {
		return baseTime, nil
	}

	timeFloat, err := strconv.ParseFloat(data, 64)
	if err != nil {
		//err 说明无法转化成float64 那么有可能本身是字符串时间进行返回
		timeTime, err := time.Parse("2006-01-02 15:04:05", data)
		if err != nil {
			return time.Time{}, fmt.Errorf("未知的时间类型")
		} else {
			return timeTime, nil
		}
	} else {
		timeTime, err := excelize.ExcelDateToTime(timeFloat, false)
		if err != nil {
			return time.Time{}, fmt.Errorf("未知的时间类型")
		} else {
			return timeTime, nil
		}
	}
}
