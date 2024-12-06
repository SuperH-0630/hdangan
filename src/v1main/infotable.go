package v1main

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"reflect"
	"strings"
)

var TopHeaderData = []string{"卷宗号", "卷宗类型", "文件联合编号", "文件编号", "文件组内编号", "姓名", "曾用名", "身份证", "性别", "出生日期", "备注", "详情"}
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

type MainTable struct {
	fileTable   *widget.Table
	window      *CtrlWindow
	InfoFile    []model.File
	InfoData    [][]string
	whereInfo   model.SearchWhere
	fileSetType model.FileSetType
}

func CreateInfoTable(rt runtime.RunTime, window *CtrlWindow) *MainTable {
	var width = make([]float32, 12)
	var idWidth float32 = 0

	fileTable := widget.NewTableWithHeaders(
		func() (rows int, cols int) {
			return 0, 0
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("暂无数据")
		}, func(id widget.TableCellID, object fyne.CanvasObject) {

		})

	m := &MainTable{
		fileTable:   fileTable,
		window:      ctrlWindow,
		fileSetType: model.QianRu,
	}

	fileTable.Length = func() (rows int, cols int) {
		return len(m.InfoData), len(TopHeaderData)
	}

	fileTable.CreateCell = func() fyne.CanvasObject {
		return widget.NewLabel("暂无数据")
	}

	fileTable.UpdateCell = func(id widget.TableCellID, object fyne.CanvasObject) {
		l := object.(*widget.Label)
		l.SetText(m.InfoData[id.Row][id.Col])
		width[id.Col] = fmax(width[id.Col], l.Size().Width, l.MinSize().Width)
		fileTable.SetColumnWidth(id.Col, width[id.Col])
	}

	fileTable.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("未知")
	}

	fileTable.UpdateHeader = func(id widget.TableCellID, template fyne.CanvasObject) {
		if id.Row == -1 {
			l := template.(*widget.Label)
			l.SetText(TopHeaderData[id.Col])
			width[id.Col] = fmax(width[id.Col], l.Size().Width, l.MinSize().Width)
			fileTable.SetColumnWidth(id.Col, width[id.Col])
		} else if id.Col == -1 {
			l := template.(*widget.Label)
			l.SetText(fmt.Sprintf("%02d", id.Row+1)) // 从1开始
			idWidth = fmax(idWidth, l.Size().Width, l.MinSize().Width)
			fileTable.SetColumnWidth(-1, idWidth)
		}
	}

	fileTable.OnSelected = func(id widget.TableCellID) {
		rt.Action()
		if id.Col == xiangQingIndex {
			if id.Row >= 0 && id.Row < len(m.InfoData) {
				file := m.InfoFile[id.Row]
				ShowInfo(rt, file, func(rt runtime.RunTime) {
					m.UpdateTable(rt, m.fileSetType, 0, window.menu.NowPage)
				})
			}
		}
		fileTable.UnselectAll()
	}

	fileTable.OnUnselected = func(id widget.TableCellID) {
		rt.Action()
	}

	m.UpdateTable(rt, m.fileSetType, 0, window.menu.NowPage)
	return m
}

func (m *MainTable) UpdateTableInfo(rt runtime.RunTime) {
	res := make([][]string, len(m.InfoFile))

	for i, j := range m.InfoFile {
		res[i] = make([]string, len(TopHeaderData))
		f := j.GetFile()

		if !f.OldName.Valid {
			f.OldName.String = ""
		}

		if !f.IDCard.Valid {
			f.IDCard.String = ""
		}

		if !f.Comment.Valid {
			f.Comment.String = ""
		}

		sex := "男性"
		if !f.IsMan {
			sex = "女性"
		}

		res[i][0] = fmt.Sprintf("%03d", f.FileSetID)
		res[i][1] = fmt.Sprintf("%s", model.FileSetTypeName[m.fileSetType])
		res[i][2] = fmt.Sprintf("%03d", f.FileUnionID)
		res[i][3] = fmt.Sprintf("%03d", f.FileID)
		res[i][4] = fmt.Sprintf("%03d", f.FileGroupID)
		res[i][5] = f.Name
		res[i][6] = f.OldName.String
		res[i][7] = f.IDCard.String
		res[i][8] = sex
		res[i][9] = f.Birthday.Format("2006-01-02")
		res[i][10] = f.Comment.String
		res[i][11] = "点击查看"
	}

	m.InfoData = res
}

func (m *MainTable) UpdateTable(rt runtime.RunTime, fileSetType model.FileSetType, pageItemCount int, p int64) {
	if pageItemCount <= 0 {
		pageItemCount = defaultItemCount
	}

	maker, ok := model.FileSetTypeMaker[fileSetType]
	if !ok {
		return
	}

	tptr := reflect.TypeOf(maker())
	if tptr.Kind() != reflect.Ptr {
		return
	}

	t := tptr.Elem()
	if t.Kind() != reflect.Struct {
		return
	}

	if strings.HasPrefix(t.Name(), "File") {
		return
	}

	resValue := reflect.MakeSlice(reflect.SliceOf(t), 0, pageItemCount)
	res := resValue.Interface()
	pageMax, err := model.GetPageData(rt, fileSetType, pageItemCount, p, &m.whereInfo, &res)
	if err != nil {
		dialog.ShowError(fmt.Errorf("获取数据库档案信息错误。"), m.window.window)
		return
	}

	m.InfoFile = make([]model.File, 0, pageItemCount)
	for i := 0; i < resValue.Len(); i++ {
		elem := resValue.Index(i)
		etype := elem.Type()
		if !etype.Implements(reflect.TypeOf((*model.File)(nil)).Elem()) {
			continue
		}

		m.InfoFile = append(m.InfoFile, elem.Interface().(model.File))
	}

	m.UpdateTableInfo(rt)
	m.window.menu.ChangePageMenuItem(rt, pageItemCount, p, pageMax, fmt.Sprintf("本页共显示数据：%d条。", len(m.InfoFile)))
	m.window.menu.ChangeFileSetModelItem(rt, pageItemCount, fmt.Sprintf("本页共显示数据：%d条。", len(m.InfoFile)))
	m.fileTable.Refresh()
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
