package v1main

import (
	"database/sql"
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

func ShowEdit(rt runtime.RunTime, fc model.File, refresh func(rt runtime.RunTime)) {
	f := fc.GetFile()

	infoWindow := rt.App().NewWindow(fmt.Sprintf("编辑信息-%s-%d", f.Name, f.FileID))

	infoWindow.SetOnClosed(func() {
		rt.Action()
		infoWindow = nil
	})
	infoWindow.SetCloseIntercept(func() {
		rt.Action()
		WinClose(infoWindow)
		infoWindow = nil
	})

	comment := widget.NewMultiLineEntry()
	comment.Text = strToStr(f.Comment, "")
	comment.Wrapping = fyne.TextWrapWord
	comment.OnChanged = func(s string) {
		rt.Action()
		comment.Text = strToStr(f.Comment, "")
	}

	material := widget.NewMultiLineEntry()
	material.Text = strToStr(f.Comment, "")
	material.Wrapping = fyne.TextWrapWord
	material.OnChanged = func(s string) {
		rt.Action()
		material.Text = strToStr(f.Comment, "")
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

	leftLayout := layout.NewFormLayout()
	left := container.New(leftLayout,
		widget.NewLabel("姓名："),
		newEntry(fmt.Sprintf("%s", f.Name), &f.Name),

		widget.NewLabel("曾用名："),
		newEntryValid(fmt.Sprintf("%s", f.OldName.String), &f.OldName),

		widget.NewLabel("身份证号："),
		newEntryValid(fmt.Sprintf("%s", f.IDCard.String), &f.IDCard),

		widget.NewLabel("性别："),
		newSexSelect(f.IsMan, &f.IsMan),

		widget.NewLabel("出生日期："),
		newDatePicker(f.Birthday, &f.Birthday, infoWindow),

		widget.NewLabel("卷宗备注："),
		comment,
	)

	rightDigit := make([]fyne.CanvasObject, 0, 10)

	switch ff := fc.(type) {
	case *model.FileQianRu:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect(ff.Type, ff.FileSetType, &ff.Type),

			widget.NewLabel("旧地址："),
			newEntry(ff.OldLocation, &ff.OldLocation),

			widget.NewLabel("新地址："),
			newEntry(ff.NewLocation, &ff.NewLocation),
		)
	case *model.FileChuSheng:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect(ff.Type, ff.FileSetType, &ff.Type),

			widget.NewLabel("地址："),
			newEntry(ff.NewLocation, &ff.NewLocation),
		)
	case *model.FileQianChu:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect(ff.Type, ff.FileSetType, &ff.Type),

			widget.NewLabel("新地址："),
			newEntry(ff.NewLocation, &ff.NewLocation),
		)
	case *model.FileSiWang:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect(ff.Type, ff.FileSetType, &ff.Type),

			widget.NewLabel("地址："),
			newEntry(ff.Location, &ff.Location),
		)
	case *model.FileBianGeng:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect(ff.Type, ff.FileSetType, &ff.Type),

			widget.NewLabel("地址："),
			newEntry(ff.Location, &ff.Location),
		)
	case *model.FileSuoNeiYiJu:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect(ff.Type, ff.FileSetType, &ff.Type),

			widget.NewLabel("地址："),
			newEntry(ff.Location, &ff.Location),
		)
	case *model.FileSuoJianYiJu:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect(ff.Type, ff.FileSetType, &ff.Type),

			widget.NewLabel("地址："),
			newEntry(ff.Location, &ff.Location),
		)
	case *model.FileNongZiZhuanFei:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect(ff.Type, ff.FileSetType, &ff.Type),

			widget.NewLabel("地址："),
			newEntry(ff.Location, &ff.Location),
		)
	case *model.FileYiZhanShiQianYiZheng:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			newFileTypeSelect(ff.Type, ff.FileSetType, &ff.Type),

			widget.NewLabel("地址："),
			newEntry(ff.Location, &ff.Location),
		)
	}

	rightDigit = append(rightDigit,
		widget.NewLabel("办理时间："),
		newDatePicker(f.Time, &f.Time, infoWindow),

		widget.NewLabel("备考："),
		newFileBeiKaoSelect(f.BeiKao.String, f.FileSetType, &f.BeiKao, &f.Material, infoWindow),

		widget.NewLabel("材料："),
		material,
	)

	rightLayout := layout.NewFormLayout()
	right := container.New(rightLayout, rightDigit...)

	upBox := container.NewHBox(left, right)

	save := widget.NewButton("保存", func() {
		rt.Action()
		err := checkAllInputRight()
		if err != nil {
			dialog.ShowError(fmt.Errorf("请检查错误：%s", err.Error()), infoWindow)
			return
		}
		dialog.ShowConfirm("更新？", "你确定要更新嘛？", func(b bool) {
			rt.Action()
			if b {
				err := model.SaveFile(rt, f)
				if err != nil {
					dialog.ShowError(fmt.Errorf("数据库错误: %s", err.Error()), infoWindow)
				}
				refresh(rt)
				WinClose(infoWindow)
				infoWindow = nil
			}
		}, infoWindow)
	})

	cancle := widget.NewButton("丢弃", func() {
		rt.Action()
		dialog.ShowConfirm("放弃？", "你确定要放弃你的操作码？", func(b bool) {
			rt.Action()
			if b {
				WinClose(infoWindow)
				infoWindow = nil
				refresh(rt)
			}
		}, infoWindow)
	})

	downBox := container.NewHBox(save, cancle)
	downCenterBox := container.NewCenter(downBox)

	gg := NewBg(5, 20)

	box := container.NewVBox(upBox, gg, downCenterBox)
	cbox := container.NewCenter(box)

	bg := NewBg(fmax(cbox.MinSize().Width, cbox.Size().Width, 220),
		fmax(cbox.MinSize().Height, cbox.Size().Height, 360))

	lastContainer := container.NewStack(bg, cbox)

	infoWindow.SetContent(lastContainer)

	infoWindow.SetFixedSize(true)
	infoWindow.Show()
	infoWindow.CenterOnScreen()
}

var entryList []*widget.Entry

func newEntry(data string, input *string) *widget.Entry {
	entry := widget.NewEntry()
	entry.Text = data

	entry.OnChanged = func(s string) {
		if entry.Validate() == nil {
			*input = s
		}
	}

	entryList = append(entryList, entry)

	return entry
}

func newEntryNumber(data int64, input *int64) *widget.Entry {
	entry := widget.NewEntry()
	entry.Text = fmt.Sprintf("%d", data)

	entry.Validator = func(s string) error {
		_, err := strconv.ParseInt(s, 10, 64)
		return err
	}

	entry.OnChanged = func(s string) {
		if entry.Validate() == nil {
			n, err := strconv.ParseInt(s, 10, 64)
			if err == nil {
				*input = n
			}
		}
	}

	entryList = append(entryList, entry)
	return entry
}

func newEntryValid(data string, input *sql.NullString) *widget.Entry {
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

	entryList = append(entryList, entry)
	return entry
}

func newFileTypeSelect(data string, fst model.FileSetType, input *string) fyne.CanvasObject {
	c, err := systeminit.GetInit()
	if err != nil {
		return newEntry(data, input)
	}

	fstName, ok := model.FileSetTypeName[fst]
	if !ok {
		return newEntry(data, input)
	}

	fileType, ok := c.Yaml.File.FileType[fstName]
	if !ok {
		return newEntry(data, input)
	}

	if !func() bool {
		for _, t := range fileType {
			if t == data {
				return true
			}
		}
		return false
	}() {
		return newEntry(data, input)
	}

	sel := widget.NewSelect(fileType, func(s string) {
		*input = s
	})

	sel.Selected = data
	return sel
}

func newFileBeiKaoSelect(data string, fst model.FileSetType, input *sql.NullString, materialInput *sql.NullString, w fyne.Window) fyne.CanvasObject {
	c, err := systeminit.GetInit()
	if err != nil {
		return newEntryValid(data, input)
	}

	fstName, ok := model.FileSetTypeName[fst]
	if !ok {
		return newEntryValid(data, input)
	}

	beiKao, ok := c.Yaml.File.BeiKao[fstName]
	if !ok {
		return newEntryValid(data, input)
	}

	if !func() bool {
		for _, t := range beiKao {
			if t.Name == data {
				return true
			}
		}
		return false
	}() {
		return newEntryValid(data, input)
	}

	for _, i := range beiKao {
		for _, j := range beiKao {
			if i.Name == j.Name {
				return newEntryValid(data, input)
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

	sel.Selected = data
	return sel
}

func newSexSelect(data bool, isMan *bool) *widget.Select {
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

func newDatePicker(data time.Time, input *time.Time, w fyne.Window) *widget.Button {
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

	return btn
}

func checkAllInputRight() error {
	for _, e := range entryList {
		err := e.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
