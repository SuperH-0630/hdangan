package v1main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
)

var TopHeaderDataRecord = []string{"发起时间", "迁入迁出状态", "借出人", "借出单位", "借入人", "借入单位", "详情"}
var xiangQingIndexRecord = -1

const defaultItemCountRecord = model.DefaultPageItemCount

func init() {
	xiangQingIndexRecord = -1
	for k, i := range TopHeaderDataRecord {
		if i == "详情" {
			xiangQingIndexRecord = k
			break
		}
	}
}

type RecordTable struct {
	Record  *RecordWindow
	Table   *widget.Table
	Width   []float32
	IdWidth float32
}

func CreateRecordTable(rt runtime.RunTime, w *RecordWindow) *RecordTable {
	w.Table = &RecordTable{
		Record: w,
		Table: widget.NewTableWithHeaders(
			func() (rows int, cols int) {
				return 0, 0
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("暂无数据")
			}, func(id widget.TableCellID, object fyne.CanvasObject) {

			}),
		Width:   make([]float32, 12),
		IdWidth: 0,
	}

	w.Table.Table.Length = func() (rows int, cols int) {
		return len(w.Table.Record.InfoRecord), len(TopHeaderDataRecord)
	}

	w.Table.Table.CreateCell = func() fyne.CanvasObject {
		return widget.NewLabel("暂无数据")
	}

	w.Table.Table.UpdateCell = func(id widget.TableCellID, object fyne.CanvasObject) {
		l := object.(*widget.Label)
		l.SetText(w.Table.Record.InfoDataRecord[id.Row][id.Col])
		w.Table.Width[id.Col] = fmax(w.Table.Width[id.Col], l.Size().Width, l.MinSize().Width)
		w.Table.Table.SetColumnWidth(id.Col, w.Table.Width[id.Col])
	}

	w.Table.Table.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("未知")
	}

	w.Table.Table.UpdateHeader = func(id widget.TableCellID, template fyne.CanvasObject) {
		if id.Row == -1 {
			l := template.(*widget.Label)
			l.SetText(TopHeaderDataRecord[id.Col])
			w.Table.Width[id.Col] = fmax(w.Table.Width[id.Col], l.Size().Width, l.MinSize().Width)
			w.Table.Table.SetColumnWidth(id.Col, w.Table.Width[id.Col])
		} else if id.Col == -1 {
			l := template.(*widget.Label)
			l.SetText(fmt.Sprintf("%02d", id.Row+1)) // 从1开始
			w.Table.IdWidth = fmax(w.Table.IdWidth, l.Size().Width, l.MinSize().Width)
			w.Table.Table.SetColumnWidth(-1, w.Table.IdWidth)
		}
	}

	w.Table.Table.OnSelected = func(id widget.TableCellID) {
		rt.Action()
		if id.Col == xiangQingIndexRecord {
			if id.Row >= 0 && id.Row < len(w.Table.Record.InfoRecord) {
				record := w.Table.Record.InfoRecord[id.Row]
				ShowRecordInfo(rt, w.Table.Record.Window, w.Table.Record.FileWindow, w.Table.Record.File, &record, func(rt runtime.RunTime) {
					w.Table.UpdateTableRecord(rt, 0, w.NowPage)
				})
			}
		}
		w.Table.Table.UnselectAll()
	}

	w.Table.Table.OnUnselected = func(id widget.TableCellID) {
		rt.Action()
	}

	return w.Table
}

func (c *RecordTable) UpdateTableInfoRecord(rt runtime.RunTime, record []model.FileMoveRecord) {
	res := make([][]string, len(record))

	for i, f := range record {
		res[i] = make([]string, len(TopHeaderDataRecord))

		res[i][0] = f.MoveTime.Format("2006-01-02 15:04:05")
		res[i][1] = f.MoveStatus
		res[i][2] = strToStr(f.MoveInPeopleName, "暂无")
		res[i][3] = strToStr(f.MoveInPeopleUnit, "暂无")
		res[i][4] = strToStr(f.MoveOutPeopleName, "暂无")
		res[i][5] = strToStr(f.MoveOutPeopleUnit, "暂无")
		res[i][6] = "点击查看"
	}

	c.Record.InfoRecord = record
	c.Record.InfoDataRecord = res
}

func (c *RecordTable) UpdateTableRecord(rt runtime.RunTime, pageItemCount int, p int64) {
	if pageItemCount <= 0 {
		pageItemCount = defaultItemCountRecord
	}

	record, pageMax, err := model.GetPageDataRecord(rt, c.Record.File, pageItemCount, p, &c.Record.SearchRecord)
	if err != nil {
		dialog.ShowError(fmt.Errorf("获取数据库档案信息错误。"), c.Record.Window)
	}

	c.UpdateTableInfoRecord(rt, record)
	c.Record.Menu.ChangePageMenuItemRecord(rt, pageItemCount, p, pageMax, fmt.Sprintf("本页共显示数据：%d条。", len(record)))
	c.Table.Refresh()
}

func (c *RecordTable) FirstUpdateData(rt runtime.RunTime) {
	c.UpdateTableRecord(rt, 0, c.Record.NowPage)
}
