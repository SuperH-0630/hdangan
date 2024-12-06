package excelio

import (
	"fmt"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strings"
)

var OutputTitle = map[model.FileSetType][]string{
	model.QianRu:               {"卷宗ID", "联合档案ID", "档案ID", "组内ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "旧地址", "新地址", "办理时间", "备考", "材料页数", "人数", "材料"},
	model.ChuSheng:             {"卷宗ID", "联合档案ID", "档案ID", "组内ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "地址", "办理时间", "备考", "材料页数", "人数", "材料"},
	model.QianChu:              {"卷宗ID", "联合档案ID", "档案ID", "组内ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "新地址", "办理时间", "备考", "材料页数", "人数", "材料"},
	model.SiWang:               {"卷宗ID", "联合档案ID", "档案ID", "组内ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "地址", "办理时间", "备考", "材料页数", "人数", "材料"},
	model.BianGeng:             {"卷宗ID", "联合档案ID", "档案ID", "组内ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "地址", "办理时间", "备考", "材料页数", "人数", "材料"},
	model.SuoNeiYiJu:           {"卷宗ID", "联合档案ID", "档案ID", "组内ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "地址", "办理时间", "备考", "材料页数", "人数", "材料"},
	model.SuoJianYiJu:          {"卷宗ID", "联合档案ID", "档案ID", "组内ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "地址", "办理时间", "备考", "材料页数", "人数", "材料"},
	model.NongZiZhuanFei:       {"卷宗ID", "联合档案ID", "档案ID", "组内ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "地址", "办理时间", "备考", "材料页数", "人数", "材料"},
	model.YiZhanShiQianYiZheng: {"卷宗ID", "联合档案ID", "档案ID", "组内ID", "姓名", "曾用名", "身份证", "性别", "出生日期", "是否同上", "备注", "类型", "地址", "办理时间", "备考", "材料页数", "人数", "材料"},
}

func OutputFile(rt runtime.RunTime, fst model.FileSetType, savepath string, files []model.File, s *model.SearchWhere) error {
	var err error

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	tit, ok := InputTitle[fst]
	if !ok {
		return err
	}

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

	timeStyleID, err := f.NewStyle(&excelize.Style{})
	if err != nil {
		return err
	}

	for i, t := range tit {
		if strings.Contains(t, "日期") || strings.Contains(t, "时间") {
			err = f.SetColStyle(sheetName, Header[i], timeStyleID)
			if err != nil {
				return err
			}
		}
	}

	for i, k := range tit {
		err = f.SetCellStr(sheetName, fmt.Sprintf("%s1", Header[i]), k)
		if err != nil {
			return err
		}
	}

	if files == nil || len(files) == 0 {
		maker, ok := model.FileSetTypeMaker[fst]
		if !ok {
			return fmt.Errorf("file set type not found")
		}

		tptr := reflect.TypeOf(maker())
		if tptr.Kind() != reflect.Ptr {
			return fmt.Errorf("file set type not found")
		}

		t := tptr.Elem()
		if t.Kind() != reflect.Struct {
			return fmt.Errorf("file set type not found")
		}

		if strings.HasPrefix(t.Name(), "File") {
			return fmt.Errorf("file set type not found")
		}

		resValue := reflect.MakeSlice(reflect.SliceOf(t), 0, 20)
		res := resValue.Interface()

		err = model.GetAllFile(rt, fst, s, res)
		if err != nil {
			return err
		}

		for i := 0; i < resValue.Len(); i++ {
			file := resValue.Index(i).Interface().(model.File)
			res, err := fileToStringLst(rt, file)
			if err != nil {
				continue
			}
			for d, v := range res {
				err = f.SetCellStr(sheetName, fmt.Sprintf("%s%d", Header[d], i+2), v)
				if err != nil {
					continue
				}
			}
		}
	} else {
		for i, file := range files {
			res, err := fileToStringLst(rt, file)
			if err != nil {
				continue
			}
			for d, v := range res {
				err = f.SetCellStr(sheetName, fmt.Sprintf("%s%d", Header[d], i+2), v)
				if err != nil {
					continue
				}
			}
		}
	}

	err = f.SaveAs(savepath)
	if err != nil {
		return err
	}

	return nil
}

func fileToStringLst(rt runtime.RunTime, file model.File) ([]string, error) {
	f := file.GetFile()

	res := make([]string, 0, 15)

	res = append(res, fmt.Sprintf("%d", f.FileSetID))
	res = append(res, fmt.Sprintf("%d", f.FileUnionID))
	res = append(res, fmt.Sprintf("%d", f.FileID))
	res = append(res, fmt.Sprintf("%d", f.FileGroupID))

	res = append(res, fmt.Sprintf("%s", f.Name))
	res = append(res, fmt.Sprintf("%s", f.OldName.String))
	res = append(res, fmt.Sprintf("%s", f.IDCard.String))

	if f.IsMan {
		res = append(res, "男性")
	} else {
		res = append(res, "女性")
	}

	res = append(res, fmt.Sprintf("%s", f.Birthday.Format("2006-01-02")))
	if f.SameAsAbove {
		res = append(res, "同上")
	} else {
		res = append(res, "非同上")
	}

	switch ff := file.(type) {
	case *model.FileQianRu:
		res = append(res, model.FileSetTypeName[ff.FileSetType])
		res = append(res, ff.OldLocation)
		res = append(res, ff.NewLocation)
	case *model.FileChuSheng:
		res = append(res, model.FileSetTypeName[ff.FileSetType])
		res = append(res, ff.NewLocation)
	case *model.FileSiWang:
		res = append(res, model.FileSetTypeName[ff.FileSetType])
		res = append(res, ff.Location)
	case *model.FileBianGeng:
		res = append(res, model.FileSetTypeName[ff.FileSetType])
		res = append(res, ff.Location)
	case *model.FileSuoNeiYiJu:
		res = append(res, model.FileSetTypeName[ff.FileSetType])
		res = append(res, ff.Location)
	case *model.FileSuoJianYiJu:
		res = append(res, model.FileSetTypeName[ff.FileSetType])
		res = append(res, ff.Location)
	case *model.FileNongZiZhuanFei:
		res = append(res, model.FileSetTypeName[ff.FileSetType])
		res = append(res, ff.Location)
	case *model.FileYiZhanShiQianYiZheng:
		res = append(res, model.FileSetTypeName[ff.FileSetType])
		res = append(res, ff.Location)
	default:
		return nil, fmt.Errorf("file set type error")
	}

	res = append(res, fmt.Sprintf("%s", f.Time.Format("2006-01-02 15:04:05")))
	res = append(res, fmt.Sprintf("%s", f.BeiKao.String))
	res = append(res, fmt.Sprintf("%d", f.PageCount))
	res = append(res, fmt.Sprintf("%d", f.PeopleCount))
	res = append(res, fmt.Sprintf("%s", f.Material.String))

	return res, nil
}
