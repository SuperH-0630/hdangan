package excelio

import (
	"fmt"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"github.com/xuri/excelize/v2"
	"strings"
)

var OutputRecordTitle = []string{
	"出借记录ID", "联合档案ID", "状态", "借出时间", "借出人", "借出单位", "借入人", "借入单位", "备注",
}

func OutputFileRecord(rt runtime.RunTime, savepath string, file model.File, record []model.FileMoveRecord, s *model.SearchRecord) error {
	var err error

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	tit := OutputRecordTitle

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

	if record == nil || len(record) == 0 {
		record, err = model.GetAllRecord(rt, file, s)
		if err != nil {
			return err
		}
	}

	for i, r := range record {
		res, err := recordToStringLst(rt, r)
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

	err = f.SaveAs(savepath)
	if err != nil {
		return err
	}

	return nil
}

func recordToStringLst(rt runtime.RunTime, record model.FileMoveRecord) ([]string, error) {
	res := make([]string, 0, 15)

	res = append(res, fmt.Sprintf("%d", record.ID))
	res = append(res, fmt.Sprintf("%d", record.FileUnionID))
	res = append(res, fmt.Sprintf("%s", record.MoveStatus))
	res = append(res, fmt.Sprintf("%s", record.MoveTime.Format("2006-01-02 15:04:05")))

	res = append(res, fmt.Sprintf("%s", record.MoveOutPeopleName.String))
	res = append(res, fmt.Sprintf("%s", record.MoveOutPeopleUnit.String))

	res = append(res, fmt.Sprintf("%s", record.MoveInPeopleName.String))
	res = append(res, fmt.Sprintf("%s", record.MoveInPeopleUnit.String))

	res = append(res, fmt.Sprintf("%s", record.MoveComment.String))

	return res, nil
}
