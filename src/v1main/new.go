package v1main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/SuperH-0630/hdangan/src/fail"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"github.com/SuperH-0630/hdangan/src/systeminit"
	"strconv"
	"time"
)

func ShowNew(rt runtime.RunTime, refresh func(rt runtime.RunTime)) {
	config, err := systeminit.GetInit()
	if errors.Is(err, systeminit.LuckyError) {
		fail.ToFail(err.Error())
		return
	}

	newWindow := rt.App().NewWindow("创建新卷宗")

	newWindow.SetOnClosed(func() {
		rt.Action()
		WinClose(newWindow)
		newWindow = nil
	})
	newWindow.SetCloseIntercept(func() {
		rt.Action()
		WinClose(newWindow)
		newWindow = nil
	})

	fid, err := model.GetNewFileID(rt)
	if err != nil {
		dialog.ShowError(fmt.Errorf("获取新的档案ID出错：%s", err.Error()), newWindow)
	}

	now := time.Now()

	f := model.File{
		FileID:      fid,
		FirstMoveIn: now,
		LastMoveIn:  now,
		MoveStatus:  config.Yaml.Move.MoveInStatus,
	}

	fileComment := widget.NewMultiLineEntry()
	fileComment.Text = strToStr(f.FileComment, "")
	fileComment.Wrapping = fyne.TextWrapWord

	formLayout := layout.NewFormLayout()
	form := container.New(formLayout,
		widget.NewLabel("卷宗号："),
		newFileIDEntry4(fmt.Sprintf("%d", f.FileID), &f.FileID),

		widget.NewLabel("姓名："),
		newEntry4(fmt.Sprintf("%s", f.Name), &f.Name),

		widget.NewLabel("身份证号："),
		newEntry4(fmt.Sprintf("%s", f.IDCard), &f.IDCard),

		widget.NewLabel("户籍地："),
		newEntry4(fmt.Sprintf("%s", f.Location), &f.Location),

		widget.NewLabel("卷宗标题："),
		newEntry4(fmt.Sprintf("%s", f.FileTitle), &f.FileTitle),

		widget.NewLabel("卷宗类型："),
		newFileTypeSelect4(fmt.Sprintf("%s", f.FileType), config.Yaml.File.FileType, &f.FileType),

		widget.NewLabel("卷宗备注："),
		fileComment,
	)

	upBox := container.NewHBox(form)

	save := widget.NewButton("保存", func() {
		rt.Action()
		err := checkAllInputRight4()
		if err != nil {
			dialog.ShowError(fmt.Errorf("请检查错误：%s", err.Error()), newWindow)
			return
		}

		err = checkNewFile(&f)
		if err != nil {
			dialog.ShowError(fmt.Errorf("请检查错误：%s", err.Error()), newWindow)
			return
		}

		dialog.ShowConfirm("创建？", "你确定要新增卷宗嘛？", func(b bool) {
			rt.Action()
			if b {
				err := model.CreateFile(&f)
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

	bg := NewBg(max(cbox.MinSize().Width, cbox.Size().Width, 220),
		max(cbox.MinSize().Height, cbox.Size().Height, 360))

	lastContainer := container.NewStack(bg, cbox)

	newWindow.SetContent(lastContainer)

	newWindow.Show()
	newWindow.CenterOnScreen()
}

func checkNewFile(f *model.File) error {
	if f.ID != 0 {
		return fmt.Errorf("系统错误")
	}

	if f.FileID <= 0 {
		return fmt.Errorf("卷宗号必须大于0")
	}

	if len(f.Name) <= 0 || len(f.Name) >= 45 {
		return fmt.Errorf("姓名必填，最大45字符")
	}

	if len(f.IDCard) <= 0 || len(f.IDCard) >= 20 {
		return fmt.Errorf("身份证必填，最大20字符")
	}

	if len(f.Location) <= 0 || len(f.Location) >= 20 {
		return fmt.Errorf("户籍地必填，最大145字符")
	}

	if len(f.FileTitle) <= 0 || len(f.FileTitle) >= 45 {
		return fmt.Errorf("卷宗标题必填，最大45字符")
	}

	if len(f.FileType) <= 0 || len(f.FileType) >= 45 {
		return fmt.Errorf("卷宗类型必填，最大15字符")
	}

	if f.FileComment.Valid {
		if len(f.FileComment.String) == 0 {
			f.FileComment.Valid = false
		}
	} else {
		f.FileComment.Valid = true
	}

	return nil
}

var entryList4 []*widget.Entry

func newEntry4(data string, input *string) *widget.Entry {
	entry := widget.NewEntry()
	entry.Text = data

	entry.OnChanged = func(s string) {
		if entry.Validate() == nil {
			*input = s
		}
	}

	entryList4 = append(entryList4, entry)
	return entry
}

func newFileIDEntry4(data string, input *int64) *widget.Entry {
	entry := widget.NewEntry()
	entry.Text = data
	entry.Validator = func(s string) error {
		n, err := strconv.ParseInt(s, 0, 64)
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
			n, err := strconv.ParseInt(s, 64, 10)
			if err == nil {
				*input = n
			}
		}
	}

	entryList4 = append(entryList4, entry)
	return entry
}

func newFileTypeSelect4(data string, options []string, input *string) *widget.Select {
	if data != "" {
		func() {
			for _, option := range options {
				if option == data {
					return
				}
			}
			options = append(options, data)
		}()
	}

	sel := widget.NewSelect(options, func(s string) {
		*input = s
	})

	sel.Selected = data
	return sel
}

func checkAllInputRight4() error {
	for _, e := range entryList4 {
		err := e.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
