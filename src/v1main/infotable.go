package v1main

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
)

var InfoFile []model.File
var InfoData = make([][]string, 0, 0)
var whereInfo = model.SearchWhere{}

var TopHeaderData = []string{"卷宗号", "姓名", "身份证", "户籍地", "卷宗标题", "卷宗类型", "最早迁入时间", "最后迁入（归还）时间", "迁入迁出状态", "迁出人姓名", "迁出人工作单位", "详情"}
var xiangQingIndex = -1

const defaultItemCount = model.DefaultPageItemCount

func init() {
	xiangQingIndex = -1
	for k, i := range TopHeaderData {
		if i == "详情" {
			xiangQingIndex = k
			break
		}
	}
}

var fileTable *widget.Table

func CreateInfoTable(rt runtime.RunTime, window fyne.Window) {
	if fileTable != nil {
		return
	}

	var width = make([]float32, 12)
	var idWidth float32 = 0

	fileTable = widget.NewTableWithHeaders(
		func() (rows int, cols int) {
			return 0, 0
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("暂无数据")
		}, func(id widget.TableCellID, object fyne.CanvasObject) {

		})

	fileTable.Length = func() (rows int, cols int) {
		return len(InfoData), len(TopHeaderData)
	}

	fileTable.CreateCell = func() fyne.CanvasObject {
		return widget.NewLabel("暂无数据")
	}

	fileTable.UpdateCell = func(id widget.TableCellID, object fyne.CanvasObject) {
		l := object.(*widget.Label)
		l.SetText(InfoData[id.Row][id.Col])
		width[id.Col] = max(width[id.Col], l.Size().Width, l.MinSize().Width)
		fileTable.SetColumnWidth(id.Col, width[id.Col])
	}

	fileTable.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("未知")
	}

	fileTable.UpdateHeader = func(id widget.TableCellID, template fyne.CanvasObject) {
		if id.Row == -1 {
			l := template.(*widget.Label)
			l.SetText(TopHeaderData[id.Col])
			width[id.Col] = max(width[id.Col], l.Size().Width, l.MinSize().Width)
			fileTable.SetColumnWidth(id.Col, width[id.Col])
		} else if id.Col == -1 {
			l := template.(*widget.Label)
			l.SetText(fmt.Sprintf("%02d", id.Row+1)) // 从1开始
			idWidth = max(idWidth, l.Size().Width, l.MinSize().Width)
			fileTable.SetColumnWidth(-1, idWidth)
		}
	}

	fileTable.OnSelected = func(id widget.TableCellID) {
		rt.Action()
		if id.Col == xiangQingIndex {
			if id.Row >= 0 && id.Row < len(InfoFile) {
				file := InfoFile[id.Row]
				ShowInfo(rt, &file, func(rt runtime.RunTime) {
					UpdateTable(rt, window, fileTable, 0, NowPage)
				})
			}
		}
		fileTable.UnselectAll()
	}

	fileTable.OnUnselected = func(id widget.TableCellID) {
		rt.Action()
	}

	UpdateTable(rt, window, fileTable, 0, 1)
}

func UpdateTableInfo(rt runtime.RunTime, files []model.File) {
	res := make([][]string, len(files))

	for i, f := range files {
		res[i] = make([]string, len(TopHeaderData))

		res[i][0] = fmt.Sprintf("%03d", f.FileID)
		res[i][1] = f.Name
		res[i][2] = f.IDCard
		res[i][3] = f.Location
		res[i][4] = f.FileTitle
		res[i][5] = f.FileType
		res[i][6] = f.FirstMoveIn.Format("2006-01-02 15:04:05")
		res[i][7] = f.LastMoveIn.Format("2006-01-02 15:04:05")
		res[i][8] = f.MoveStatus
		res[i][9] = strToStr(f.MoveOutPeopleName)
		res[i][10] = strToStr(f.MoveOutPeopleUnit)
		res[i][11] = "点击查看"
	}

	InfoData = res
}

func UpdateTable(rt runtime.RunTime, window fyne.Window, table *widget.Table, pageItemCount int, p int64) {
	if pageItemCount <= 0 {
		pageItemCount = defaultItemCount
	}

	files, pageMax, err := model.GetPageData(rt, pageItemCount, p, &whereInfo)
	if err != nil {
		dialog.ShowError(fmt.Errorf("获取数据库档案信息错误。"), window)
		return
	}

	InfoFile = files

	UpdateTableInfo(rt, files)
	ChangePageMenuItem(rt, window, table, pageItemCount, p, pageMax, fmt.Sprintf("本页共显示数据：%d条。", len(files)))
	table.Refresh()
}

func timeToStr(time sql.NullTime, NWord ...string) string {
	if time.Valid {
		return time.Time.Format("2006-01-02 15:04:05")
	}

	if len(NWord) > 0 {
		return NWord[0]
	}

	return "无"
}

func strToStr(str sql.NullString, NWord ...string) string {
	if str.Valid {
		return str.String
	}

	if len(NWord) > 0 {
		return NWord[0]
	}

	return "无"
}