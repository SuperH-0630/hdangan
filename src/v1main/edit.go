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
	"strconv"
)

func ShowEdit(rt runtime.RunTime, f *model.File, refresh func(rt runtime.RunTime)) {
	config, err := systeminit.GetInit()
	if errors.Is(err, systeminit.LuckyError) {
		rt.DBConnectError(err)
		return
	} else if err != nil {
		rt.DBConnectError(fmt.Errorf("配置文件错误，请检查配置文件状态。"))
		return
	}

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

	fileComment := widget.NewMultiLineEntry()
	fileComment.Text = strToStr(f.FileComment, "")
	fileComment.Wrapping = fyne.TextWrapWord
	fileComment.OnChanged = func(text string) {
		f.FileComment = sql.NullString{String: text, Valid: len(text) != 0}
	}

	leftLayout := layout.NewFormLayout()
	left := container.New(leftLayout,
		widget.NewLabel("卷宗号："),
		newFileIDEntry(f.FileID, &f.FileID),

		widget.NewLabel("姓名："),
		newEntry(fmt.Sprintf("%s", f.Name), &f.Name),

		widget.NewLabel("身份证号："),
		newEntry(fmt.Sprintf("%s", f.IDCard), &f.IDCard),

		widget.NewLabel("户籍地："),
		newEntry(fmt.Sprintf("%s", f.Location), &f.Location),

		widget.NewLabel("卷宗标题："),
		newEntry(fmt.Sprintf("%s", f.FileTitle), &f.FileTitle),

		widget.NewLabel("卷宗类型："),
		newFileTypeSelect(fmt.Sprintf("%s", f.FileType), config.Yaml.File.FileType, &f.FileType),

		widget.NewLabel("卷宗备注："),
		fileComment,
	)

	upBox := container.NewHBox(left)

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
				err := model.SaveFile(f)
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

func newFileIDEntry(data int64, input *int64) *widget.Entry {
	entry := widget.NewEntry()
	if data <= 0 {
		entry.Text = ""
	} else {
		entry.Text = fmt.Sprintf("%d", data)
	}

	entry.Validator = func(s string) error {
		if len(s) == 0 {
			return nil
		}

		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}

		if n <= 0 {
			return fmt.Errorf("must bigger than zero")
		}

		return nil
	}

	entry.OnChanged = func(s string) {
		if entry.Validate() == nil {
			if len(s) == 0 {
				*input = 0
			} else {
				n, err := strconv.ParseInt(s, 10, 64)
				if err == nil {
					*input = n
				} else {
					*input = 0
				}
			}
		}
	}

	entryList = append(entryList, entry)
	return entry
}

func newFileTypeSelect(data string, options []string, input *string) *widget.Select {
	func() {
		for _, option := range options {
			if option == data {
				return
			}
		}
		options = append(options, data)
	}()

	sel := widget.NewSelect(options, func(s string) {
		*input = s
	})

	sel.Selected = data
	return sel
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
