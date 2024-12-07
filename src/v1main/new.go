package v1main

import (
	"database/sql"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"github.com/SuperH-0630/hdangan/src/systeminit"
	datepicker "github.com/sdassow/fyne-datepicker"
	"strconv"
	"strings"
	"time"
)

func ShowNew(rt runtime.RunTime, w *CtrlWindow, refresh func(rt runtime.RunTime)) {
	config, err := systeminit.GetInit()
	if errors.Is(err, systeminit.LuckyError) {
		rt.DBConnectError(err)
		return
	} else if err != nil {
		rt.DBConnectError(fmt.Errorf("配置文件错误，请检查配置文件状态。"))
		return
	}

	newWindow := rt.App().NewWindow("创建记录")

	newWindow.SetOnClosed(func() {
		rt.Action()
		newWindow = nil
	})
	newWindow.SetCloseIntercept(func() {
		rt.Action()
		WinClose(newWindow)
		newWindow = nil
	})

	maker, ok := model.FileSetTypeMaker[w.fileSetType]
	if !ok {
		rt.DBConnectError(fmt.Errorf("配置文件错误，请检查配置文件状态。"))
		return
	}

	fm := maker()
	record := &model.FileMoveRecord{
		MoveStatus: config.Yaml.Move.MoveInStatus,
	}
	f := fm.GetFile()

	f.PeopleCount = 1

	comment := widget.NewMultiLineEntry()
	comment.Text = strToStr(f.Comment, "")
	comment.Wrapping = fyne.TextWrapWord
	comment.OnChanged = func(s string) {
		rt.Action()
		comment.Text = strToStr(f.Comment, "")
	}

	material := widget.NewMultiLineEntry()
	material.Text = strToStr(f.Material, "")
	material.Wrapping = fyne.TextWrapWord
	material.OnChanged = func(s string) {
		rt.Action()
		material.Text = strToStr(f.Material, "")
	}

	if !f.OldName.Valid {
		f.OldName.String = ""
	}

	if !f.IDCard.Valid {
		f.IDCard.String = ""
	}

	if !f.Comment.Valid {
		f.Comment.String = ""
	}

	sameAboveDisableList := make([]fyne.Disableable, 0, 10)

	leftLayout := layout.NewFormLayout()
	left := container.New(leftLayout,
		widget.NewLabel("姓名："),
		newEntry4(fmt.Sprintf("%s", f.Name), &f.Name, nil),

		widget.NewLabel("曾用名："),
		newEntryValid4(fmt.Sprintf("%s", f.OldName.String), &f.OldName, nil),

		widget.NewLabel("身份证号："),
		newEntryValid4(fmt.Sprintf("%s", f.IDCard.String), &f.IDCard, nil),

		widget.NewLabel("性别："),
		newSexSelect4(f.IsMan, &f.IsMan),

		widget.NewLabel("出生日期："),
		newDatePicker4(f.Birthday, &f.Birthday, newWindow, nil),

		widget.NewLabel("同上："),
		newSameAboveCheck4("", &f.SameAsAbove, &sameAboveDisableList),

		widget.NewLabel("备注："),
		comment,
	)

	rightDigit := make([]fyne.CanvasObject, 0, 10)

	switch ff := fm.(type) {
	case *model.FileQianRu:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect4(ff.Type, ff.FileSetType, &ff.Type, &sameAboveDisableList),

			widget.NewLabel("旧地址："),
			newEntry4(ff.OldLocation, &ff.OldLocation, &sameAboveDisableList),

			widget.NewLabel("新地址："),
			newEntry4(ff.NewLocation, &ff.NewLocation, &sameAboveDisableList),
		)
	case *model.FileChuSheng:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect4(ff.Type, ff.FileSetType, &ff.Type, &sameAboveDisableList),

			widget.NewLabel("地址："),
			newEntry4(ff.NewLocation, &ff.NewLocation, &sameAboveDisableList),
		)
	case *model.FileQianChu:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect4(ff.Type, ff.FileSetType, &ff.Type, &sameAboveDisableList),

			widget.NewLabel("新地址："),
			newEntry4(ff.NewLocation, &ff.NewLocation, &sameAboveDisableList),
		)
	case *model.FileSiWang:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect4(ff.Type, ff.FileSetType, &ff.Type, &sameAboveDisableList),

			widget.NewLabel("地址："),
			newEntry4(ff.Location, &ff.Location, &sameAboveDisableList),
		)
	case *model.FileBianGeng:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect4(ff.Type, ff.FileSetType, &ff.Type, &sameAboveDisableList),

			widget.NewLabel("地址："),
			newEntry4(ff.Location, &ff.Location, &sameAboveDisableList),
		)
	case *model.FileSuoNeiYiJu:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect4(ff.Type, ff.FileSetType, &ff.Type, &sameAboveDisableList),

			widget.NewLabel("地址："),
			newEntry4(ff.Location, &ff.Location, &sameAboveDisableList),
		)
	case *model.FileSuoJianYiJu:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect4(ff.Type, ff.FileSetType, &ff.Type, &sameAboveDisableList),

			widget.NewLabel("地址："),
			newEntry4(ff.Location, &ff.Location, &sameAboveDisableList),
		)
	case *model.FileNongZiZhuanFei:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect4(ff.Type, ff.FileSetType, &ff.Type, &sameAboveDisableList),

			widget.NewLabel("地址："),
			newEntry4(ff.Location, &ff.Location, &sameAboveDisableList),
		)
	case *model.FileYiZhanShiQianYiZheng:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect4(ff.Type, ff.FileSetType, &ff.Type, &sameAboveDisableList),

			widget.NewLabel("地址："),
			newEntry4(ff.Location, &ff.Location, &sameAboveDisableList),
		)
	}

	rightDigit = append(rightDigit,
		widget.NewLabel("办理时间："),
		newDatePicker4(f.Time, &f.Time, newWindow, &sameAboveDisableList),

		widget.NewLabel("备考："),
		newFileBeiKaoSelect4(f.BeiKao.String, f.FileSetType, &f.BeiKao, &f.Material, newWindow, nil),

		widget.NewLabel("材料页数："),
		newEntryPage4(f.PageCount, &f.PageCount, &sameAboveDisableList),

		widget.NewLabel("材料："),
		material,

		widget.NewLabel("记录人："),
		newEntryWithNULL4(config.Yaml.Move.MoveInPeopleDefault, &record.MoveInPeopleName, &sameAboveDisableList),

		widget.NewLabel("记录单位："),
		newMoveUnitSelectWithNULL4(config.Yaml.Move.MoveInUnitDefault, config.Yaml.Move.MoveUnit, &record.MoveInPeopleUnit, &sameAboveDisableList),
	)

	rightLayout := layout.NewFormLayout()
	right := container.New(rightLayout, rightDigit...)

	upBox := container.NewHBox(left, right)

	save := widget.NewButton("保存", func() {
		rt.Action()
		err := checkAllInputRight4()
		if err != nil {
			dialog.ShowError(fmt.Errorf("请检查错误：%s", err.Error()), newWindow)
			return
		}

		dialog.ShowConfirm("创建？", "你确定要新增档案嘛？", func(b bool) {
			rt.Action()
			if b {
				err := model.CreateFile(rt, w.fileSetType, fm, record)
				if err != nil {
					dialog.ShowError(fmt.Errorf("数据库错误: %s", err.Error()), newWindow)
				}
				WinClose(newWindow)
				newWindow = nil
				refresh(rt)
			}
		}, newWindow)
	})

	cancle := widget.NewButton("丢弃", func() {
		rt.Action()
		dialog.ShowConfirm("放弃？", "你确定要放弃你的操作码？", func(b bool) {
			rt.Action()
			if b {
				WinClose(newWindow)
				newWindow = nil
				refresh(rt)
			}
		}, newWindow)
	})

	downBox := container.NewHBox(save, cancle)
	downCenterBox := container.NewCenter(downBox)

	gg := NewBg(5, 20)

	box := container.NewVBox(upBox, gg, downCenterBox)
	cbox := container.NewCenter(box)

	bg := NewBg(fmax(cbox.MinSize().Width, cbox.Size().Width, 220),
		fmax(cbox.MinSize().Height, cbox.Size().Height, 360))

	lastContainer := container.NewStack(bg, cbox)

	newWindow.SetContent(lastContainer)
	newWindow.Show()
	newWindow.CenterOnScreen()
	newWindow.SetFixedSize(true)
}

var entryList4 []*widget.Entry

func newEntry4(data string, input *string, disableLst *[]fyne.Disableable) *widget.Entry {
	entry := widget.NewEntry()
	entry.Text = data

	entry.OnChanged = func(s string) {
		if entry.Validate() == nil {
			*input = s
		}
	}

	entryList4 = append(entryList4, entry)

	if disableLst != nil {
		*disableLst = append(*disableLst, entry)
	}

	return entry
}

func newEntryWithNULL4(data string, input *sql.NullString, disableLst *[]fyne.Disableable) *widget.Entry {
	entry := widget.NewEntry()
	entry.Text = data

	entry.OnChanged = func(s string) {
		if entry.Validate() == nil {
			if len(s) == 0 {
				*input = sql.NullString{
					Valid:  false,
					String: "",
				}
			} else {
				*input = sql.NullString{
					Valid:  true,
					String: s,
				}
			}
		}
	}

	entryList4 = append(entryList4, entry)

	if disableLst != nil {
		*disableLst = append(*disableLst, entry)
	}

	return entry
}

func newMoveUnitSelectWithNULL4(data string, firstOptions []string, input *sql.NullString, disableLst *[]fyne.Disableable) *widget.Select {
	const emptySelectItem = "暂无"

	if data == "" {
		data = emptySelectItem
	}

	options := make([]string, 0, len(firstOptions)+1)
	options = append(options, emptySelectItem)

	for _, fo := range firstOptions {
		if fo != emptySelectItem && fo != "" {
			options = append(options, fo)
		}
	}

	sel := widget.NewSelect(options, func(s string) {
		if len(s) == 0 || s == emptySelectItem {
			*input = sql.NullString{
				Valid:  false,
				String: "",
			}
		} else {
			*input = sql.NullString{
				Valid:  true,
				String: s,
			}
		}
	})

	sel.PlaceHolder = ""
	sel.Selected = data

	if disableLst != nil {
		*disableLst = append(*disableLst, sel)
	}

	return sel
}

func newEntryPage4(data int64, input *int64, disableLst *[]fyne.Disableable) *widget.Entry {
	entry := widget.NewEntry()
	entry.Text = fmt.Sprintf("%d", data)

	entry.Validator = func(s string) error {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}

		if n <= 0 {
			return fmt.Errorf("page musr bigger than zero")
		}

		return nil
	}

	entry.OnChanged = func(s string) {
		if entry.Validate() == nil {
			n, err := strconv.ParseInt(s, 10, 64)
			if err == nil {
				*input = n
			}
		}
	}

	if disableLst != nil {
		*disableLst = append(*disableLst, entry)
	}

	entryList4 = append(entryList4, entry)
	return entry
}

func newEntryValid4(data string, input *sql.NullString, disableLst *[]fyne.Disableable) *widget.Entry {
	entry := widget.NewEntry()
	entry.Text = data

	entry.OnChanged = func(s string) {
		if entry.Validate() == nil {
			if len(s) == 0 {
				*input = sql.NullString{Valid: false}
			} else {
				*input = sql.NullString{Valid: true, String: s}
			}
		}
	}

	if disableLst != nil {
		*disableLst = append(*disableLst, entry)
	}

	entryList4 = append(entryList4, entry)
	return entry
}

func newFileTypeSelect4(data string, fst model.FileSetType, input *string, disableLst *[]fyne.Disableable) fyne.CanvasObject {
	c, err := systeminit.GetInit()
	if err != nil {
		return newEntry4(data, input, disableLst)
	}

	fstName, ok := model.FileSetTypeName[fst]
	if !ok {
		return newEntry4(data, input, disableLst)
	}

	fileType, ok := c.Yaml.File.FileType[fstName]
	if !ok {
		return newEntry4(data, input, disableLst)
	}

	if !func() bool {
		for _, t := range fileType {
			if t == data {
				return true
			}
		}
		return false
	}() {
		return newEntry4(data, input, disableLst)
	}

	sel := widget.NewSelect(fileType, func(s string) {
		*input = s
	})

	if disableLst != nil {
		*disableLst = append(*disableLst, sel)
	}

	sel.Selected = data
	return sel
}

func newFileBeiKaoSelect4(data string, fst model.FileSetType, input *sql.NullString, materialInput *sql.NullString, w fyne.Window, disableLst *[]fyne.Disableable) fyne.CanvasObject {
	c, err := systeminit.GetInit()
	if err != nil {
		return newEntryValid4(data, input, disableLst)
	}

	fstName, ok := model.FileSetTypeName[fst]
	if !ok {
		return newEntryValid4(data, input, disableLst)
	}

	beiKao, ok := c.Yaml.File.BeiKao[fstName]
	if !ok {
		return newEntryValid4(data, input, disableLst)
	}

	if !func() bool {
		for _, t := range beiKao {
			if t.Name == data {
				return true
			}
		}
		return false
	}() {
		return newEntryValid4(data, input, disableLst)
	}

	for _, i := range beiKao {
		for _, j := range beiKao {
			if i.Name == j.Name {
				return newEntryValid4(data, input, disableLst)
			}
		}
	}

	beiKaoLst := make([]string, 0, len(beiKao))
	for _, i := range beiKao {
		beiKaoLst = append(beiKaoLst, i.Name)
	}

	sel := widget.NewSelect(beiKaoLst, func(s string) {
		input.Valid = len(s) != 0
		input.String = s

		for _, i := range beiKao {
			if s == i.Name && len(i.Material) != 0 {
				dialog.ShowConfirm("确认", "此”备考“有预先录制的材料列表，是否立即使用？", func(b bool) {
					if !b {
						return
					}

					materialInput.Valid = true
					materialInput.String = strings.Join(i.Material, "、") + "。"
				}, w)
			}
		}

	})

	if disableLst != nil {
		*disableLst = append(*disableLst, sel)
	}

	sel.Selected = data
	return sel
}

func newSexSelect4(data bool, isMan *bool) *widget.Select {
	sel := widget.NewSelect([]string{"男性", "女性"}, func(s string) {
		*isMan = s == "男性"
	})

	if data {
		sel.Selected = "男性"
	} else {
		sel.Selected = "女性"
	}

	return sel
}

func newDatePicker4(data time.Time, input *time.Time, w fyne.Window, disableLst *[]fyne.Disableable) *widget.Button {
	btn := widget.NewButton("选择时间", func() {})

	d := datepicker.NewDatePicker(data, time.Monday, func(t time.Time, b bool) {
		if b {
			*input = t
			btn.SetText(t.Format("2006-01-02"))
		}
	})

	btn.OnTapped = func() {
		dialog.ShowCustomConfirm("选择时间", "确认", "放弃", d, d.OnActioned, w)
	}

	if disableLst != nil {
		*disableLst = append(*disableLst, btn)
	}
	return btn
}

func newSameAboveCheck4(data string, input *bool, disableLst *[]fyne.Disableable) *widget.Check {
	btn := widget.NewCheck(data, func(b bool) {
		*input = b
		for _, d := range *disableLst {
			if b {
				d.Disable()
			} else {
				d.Enable()
			}
		}
	})
	return btn
}

func checkAllInputRight4() error {
	for _, e := range entryList4 {
		if e.Disabled() {
			continue
		}

		err := e.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
