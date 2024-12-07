package excelio

import (
	"database/sql"
	"fmt"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"github.com/SuperH-0630/hdangan/src/systeminit"
	"github.com/xuri/excelize/v2"
	"io"
	"strconv"
	"strings"
	"time"
)

var InputTitle = map[model.FileSetType][]string{
	model.QianRu:               {"档案ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "旧地址", "新地址", "办理时间", "备考", "材料页数", "材料", "录入人", "录入单位"},
	model.ChuSheng:             {"档案ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "地址", "办理时间", "备考", "材料页数", "材料", "录入人", "录入单位"},
	model.QianChu:              {"档案ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "新地址", "办理时间", "备考", "材料页数", "材料", "录入人", "录入单位"},
	model.SiWang:               {"档案ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "地址", "办理时间", "备考", "材料页数", "材料", "录入人", "录入单位"},
	model.BianGeng:             {"档案ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "地址", "办理时间", "备考", "材料页数", "材料", "录入人", "录入单位"},
	model.SuoNeiYiJu:           {"档案ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "地址", "办理时间", "备考", "材料页数", "材料", "录入人", "录入单位"},
	model.SuoJianYiJu:          {"档案ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "地址", "办理时间", "备考", "材料页数", "材料", "录入人", "录入单位"},
	model.NongZiZhuanFei:       {"档案ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "地址", "办理时间", "备考", "材料页数", "材料", "录入人", "录入单位"},
	model.YiZhanShiQianYiZheng: {"档案ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "地址", "办理时间", "备考", "材料页数", "材料", "录入人", "录入单位"},
}

var Header = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S"}

var BadTitle = fmt.Errorf("表格首行的标题错误")

const (
	successAdd = iota
	successUpdate
	fail
)

func ReadFile(rt runtime.RunTime, fst model.FileSetType, reader io.ReadCloser) (int64, int64, int64, error) {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return 0, 0, 0, err
	}

	tit, ok := InputTitle[fst]
	if !ok {
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

			s := makeFile(rt, fst, t, tit)
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

func checkTitle(t []string, tit []string) bool {
	for i, s := range t {
		if tit[i] != s {
			return false
		}
	}
	return true
}

func makeFile(rt runtime.RunTime, fst model.FileSetType, t []string, tit []string) int {
	if t[0] == "" {
		file, record, err := makeNewFile(rt, fst, t, tit)
		if err != nil {
			return fail
		}

		err = model.CreateFile(rt, file.GetFile().FileSetType, file, record)
		if err != nil {
			return fail
		}
		return successAdd
	} else {
		fileID, err := strconv.ParseInt(t[0], 10, 64)
		if err != nil {
			return fail
		}

		file, err := makeUpdateFile(rt, fileID, fst, t, tit)
		if err != nil {
			return fail
		}

		err = model.SaveFile(rt, file)
		if err != nil {
			return fail
		}
		return successUpdate
	}
}

func makeNewFile(rt runtime.RunTime, fst model.FileSetType, t []string, tit []string) (model.File, *model.FileMoveRecord, error) {
	config, err := systeminit.GetInit()
	if err != nil {
		return nil, nil, err
	}

	name := t[1]
	if len(name) == 0 {
		return nil, nil, fmt.Errorf("must has name")
	}

	oldName := sql.NullString{
		Valid:  len(t[2]) != 0,
		String: t[2],
	}

	idcard := sql.NullString{
		Valid:  len(t[3]) != 0,
		String: t[3],
	}

	sex := t[4] == "女性" || t[4] == "女" || t[4] == "W" || t[4] == "WOMAN"

	birthday := time.Now()
	if len(t[5]) != 0 {
		var err error
		birthday, err = timeReader(time.Now(), t[5])
		if err != nil {
			return nil, nil, err
		}
	}

	sameAbove := t[6] == "是" || t[6] == "是的" || t[6] == "同上" || t[6] == "T" || t[6] == "True"

	comment := sql.NullString{
		Valid:  len(t[7]) != 0,
		String: t[7],
	}

	fileTimeIndex := 0
	for i, j := range tit {
		if j == "办理时间" {
			fileTimeIndex = i
			break
		}
	}

	fileTime := time.Now()
	if len(t[fileTimeIndex]) != 0 {
		var err error
		fileTime, err = timeReader(time.Now(), t[fileTimeIndex])
		if err != nil {
			return nil, nil, err
		}
	}

	beiKao := sql.NullString{
		Valid:  len(t[fileTimeIndex+1]) != 0,
		String: t[fileTimeIndex+1],
	}

	page, err := strconv.ParseInt(t[fileTimeIndex+2], 10, 64)
	if err != nil {
		return nil, nil, err
	}

	material := sql.NullString{
		Valid:  len(t[fileTimeIndex+3]) != 0,
		String: t[fileTimeIndex+3],
	}

	moveInPeople := sql.NullString{
		Valid:  len(t[fileTimeIndex+4]) != 0,
		String: t[fileTimeIndex+4],
	}

	moveInUnit := sql.NullString{
		Valid:  len(t[fileTimeIndex+5]) != 0,
		String: t[fileTimeIndex+5],
	}

	var res model.File

	abs := model.FileAbs{
		Name:     name,
		OldName:  oldName,
		IDCard:   idcard,
		IsMan:    sex,
		Birthday: birthday,
		Comment:  comment,

		SameAsAbove: sameAbove,
		PeopleCount: 1,

		Time: fileTime,

		PageCount: page,

		BeiKao:   beiKao,
		Material: material,
	}

	switch fst {
	case model.QianRu:
		fileType := t[8]
		if len(fileType) == 0 {
			return nil, nil, fmt.Errorf("must has file type")
		}

		oldLoc := t[9]
		if len(oldLoc) == 0 {
			return nil, nil, fmt.Errorf("must has old location")
		}

		newLoc := t[10]
		if len(newLoc) == 0 {
			return nil, nil, fmt.Errorf("must has new location")
		}

		res = &model.FileQianRu{
			FileAbs:     abs,
			Type:        fileType,
			OldLocation: oldLoc,
			NewLocation: newLoc,
		}
	case model.ChuSheng:
		fileType := t[8]
		if len(fileType) == 0 {
			return nil, nil, fmt.Errorf("must has file type")
		}

		loc := t[9]
		if len(loc) == 0 {
			return nil, nil, fmt.Errorf("must has location")
		}

		res = &model.FileChuSheng{
			FileAbs:     abs,
			Type:        fileType,
			NewLocation: loc,
		}
	case model.SiWang:
		fileType := t[8]
		if len(fileType) == 0 {
			return nil, nil, fmt.Errorf("must has file type")
		}

		loc := t[9]
		if len(loc) == 0 {
			return nil, nil, fmt.Errorf("must has new location")
		}

		res = &model.FileSiWang{
			FileAbs:  abs,
			Type:     fileType,
			Location: loc,
		}
	case model.BianGeng:
		fileType := t[8]
		if len(fileType) == 0 {
			return nil, nil, fmt.Errorf("must has file type")
		}

		loc := t[9]
		if len(loc) == 0 {
			return nil, nil, fmt.Errorf("must has location")
		}

		res = &model.FileBianGeng{
			FileAbs:  abs,
			Type:     fileType,
			Location: loc,
		}
	case model.SuoNeiYiJu:
		fileType := t[8]
		if len(fileType) == 0 {
			return nil, nil, fmt.Errorf("must has file type")
		}

		loc := t[9]
		if len(loc) == 0 {
			return nil, nil, fmt.Errorf("must has location")
		}

		res = &model.FileSuoNeiYiJu{
			FileAbs:  abs,
			Type:     fileType,
			Location: loc,
		}
	case model.SuoJianYiJu:
		fileType := t[8]
		if len(fileType) == 0 {
			return nil, nil, fmt.Errorf("must has file type")
		}

		loc := t[9]
		if len(loc) == 0 {
			return nil, nil, fmt.Errorf("must has new location")
		}

		res = &model.FileSuoJianYiJu{
			FileAbs:  abs,
			Type:     fileType,
			Location: loc,
		}
	case model.NongZiZhuanFei:
		fileType := t[8]
		if len(fileType) == 0 {
			return nil, nil, fmt.Errorf("must has file type")
		}

		loc := t[9]
		if len(loc) == 0 {
			return nil, nil, fmt.Errorf("must has location")
		}

		res = &model.FileNongZiZhuanFei{
			FileAbs:  abs,
			Type:     fileType,
			Location: loc,
		}
	case model.YiZhanShiQianYiZheng:
		fileType := t[8]
		if len(fileType) == 0 {
			return nil, nil, fmt.Errorf("must has file type")
		}

		loc := t[9]
		if len(loc) == 0 {
			return nil, nil, fmt.Errorf("must has location")
		}

		res = &model.FileYiZhanShiQianYiZheng{
			FileAbs:  abs,
			Type:     fileType,
			Location: loc,
		}
	default:
		return nil, nil, fmt.Errorf("file set type error")
	}

	record := &model.FileMoveRecord{
		MoveStatus:       config.Yaml.Move.MoveInStatus,
		MoveInPeopleName: moveInPeople,
		MoveInPeopleUnit: moveInUnit,
	}

	return res, record, nil
}

func makeUpdateFile(rt runtime.RunTime, fileID int64, fst model.FileSetType, t []string, tit []string) (model.File, error) {
	maker, ok := model.FileSetTypeMaker[fst]
	if !ok {
		return nil, fmt.Errorf("bad file set type")
	}

	f := maker()
	err := model.FindFile(rt, fileID, f) // f 已经是指针
	if err != nil {
		return nil, err
	}

	file := f.GetFile()
	if file.FileSetType != fst {
		return nil, fmt.Errorf("bad file set type")
	}

	name := t[1]
	if len(name) != 0 {
		file.Name = name
	}

	oldName := t[2]
	if len(oldName) != 0 {
		file.OldName = sql.NullString{
			Valid:  len(oldName) != 0,
			String: oldName,
		}
	}

	idcard := t[3]
	if len(idcard) != 0 {
		file.IDCard = sql.NullString{
			Valid:  len(idcard) != 0,
			String: idcard,
		}
	}

	if len(t[4]) != 0 {
		file.IsMan = t[4] == "女性" || t[4] == "女" || t[4] == "W" || t[4] == "WOMAN"
	}

	if len(t[5]) != 0 {
		var err error
		file.Birthday, err = timeReader(time.Now(), t[5])
		if err != nil {
			return nil, err
		}
	}

	if len(t[7]) != 0 {
		file.Comment = sql.NullString{
			Valid:  true,
			String: t[7],
		}
	}

	fileTimeIndex := 0
	for i, j := range tit {
		if j == "办理时间" {
			fileTimeIndex = i
			break
		}
	}

	if len(t[fileTimeIndex]) != 0 {
		var err error
		file.Time, err = timeReader(time.Now(), t[fileTimeIndex])
		if err != nil {
			return nil, err
		}
	}

	if len(t[fileTimeIndex+1]) != 0 {
		file.BeiKao = sql.NullString{
			Valid:  true,
			String: t[fileTimeIndex+1],
		}
	}

	if len(t[fileTimeIndex+2]) != 0 {
		file.Material = sql.NullString{
			Valid:  true,
			String: t[fileTimeIndex+2],
		}
	}

	switch ff := f.(type) {
	case *model.FileQianRu:
		fileType := t[8]
		if len(fileType) != 0 {
			ff.Type = fileType
		}

		oldLoc := t[9]
		if len(oldLoc) != 0 {
			ff.OldLocation = oldLoc
		}

		newLoc := t[10]
		if len(newLoc) != 0 {
			ff.NewLocation = newLoc
		}
	case *model.FileChuSheng:
		fileType := t[8]
		if len(fileType) != 0 {
			ff.Type = fileType
		}

		loc := t[9]
		if len(loc) != 0 {
			ff.NewLocation = loc
		}
	case *model.FileSiWang:
		fileType := t[8]
		if len(fileType) != 0 {
			ff.Type = fileType
		}

		loc := t[9]
		if len(loc) != 0 {
			ff.Location = loc
		}
	case *model.FileBianGeng:
		fileType := t[8]
		if len(fileType) != 0 {
			ff.Type = fileType
		}

		loc := t[9]
		if len(loc) != 0 {
			ff.Location = loc
		}
	case *model.FileSuoNeiYiJu:
		fileType := t[8]
		if len(fileType) != 0 {
			ff.Type = fileType
		}

		loc := t[9]
		if len(loc) != 0 {
			ff.Location = loc
		}
	case *model.FileSuoJianYiJu:
		fileType := t[8]
		if len(fileType) != 0 {
			ff.Type = fileType
		}

		loc := t[9]
		if len(loc) != 0 {
			ff.Location = loc
		}
	case *model.FileNongZiZhuanFei:
		fileType := t[8]
		if len(fileType) != 0 {
			ff.Type = fileType
		}

		loc := t[9]
		if len(loc) != 0 {
			ff.Location = loc
		}
	case *model.FileYiZhanShiQianYiZheng:
		fileType := t[8]
		if len(fileType) != 0 {
			ff.Type = fileType
		}

		loc := t[9]
		if len(loc) != 0 {
			ff.Location = loc
		}
	default:
		return nil, fmt.Errorf("file set type error")
	}

	return f, nil
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
