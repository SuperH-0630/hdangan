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
	"time"
)

type SaveWhereWindow struct {
	Window       fyne.Window
	Download     func(rt runtime.RunTime, w *SaveWhereWindow)
	SearchRecord model.SearchRecord
}

func NewSaveWhereWindow(rt runtime.RunTime, download func(rt runtime.RunTime, w *SaveWhereWindow)) *SaveWhereWindow {
	w := &SaveWhereWindow{
		Window:   rt.App().NewWindow(fmt.Sprintf("搜索筛选")),
		Download: download,
	}
	return w
}

func (w *SaveWhereWindow) create(rt runtime.RunTime) {
	config, err := systeminit.GetInit()
	if errors.Is(err, systeminit.LuckyError) {
		rt.DBConnectError(err)
		return
	} else if err != nil {
		rt.DBConnectError(fmt.Errorf("配置文件错误，请检查配置文件状态。"))
		return
	}

	var s = w.SearchRecord

	w.Window.SetOnClosed(func() {
		rt.Action()
		w.Window = nil
	})
	w.Window.SetCloseIntercept(func() {
		rt.Action()
		WinClose(w.Window)
		w.Window = nil
	})

	leftLayout := layout.NewFormLayout()
	left := container.New(leftLayout,
		widget.NewLabel("最早迁出时间（含）："),
		newTimePicker6(&s.MoveOutStart, w.Window),

		widget.NewLabel("最晚迁出时间（含）："),
		newTimePicker6(&s.MoveOutEnd, w.Window),
	)

	rightLayout := layout.NewFormLayout()
	right := container.New(rightLayout,
		widget.NewLabel("迁入迁出状态："),
		newSelect6(config.Yaml.Move.MoveStatus, &s.MoveStatus),

		widget.NewLabel("最后迁出人："),
		newEntry6(&s.MoveOutPeopleName),

		widget.NewLabel("最后迁出单位："),
		newSelect6(config.Yaml.Move.MoveUnit, &s.MoveOutPeopleUnit),
	)

	upBox := container.NewHBox(left, right)

	save := widget.NewButton("保存下载", func() {
		rt.Action()
		err := checkAllInputRight6()
		if err != nil {
			dialog.ShowError(fmt.Errorf("请检查错误：%s", err.Error()), w.Window)
			return
		}
		w.SearchRecord = s
		w.Download(rt, w)
		WinClose(w.Window)
		w.Window = nil
	})

	clearAll := widget.NewButton("清空条件", func() {
		rt.Action()
		dialog.ShowConfirm("确定？", "你是否确定要清空全部条件？？", func(b bool) {
			rt.Action()
			if b {
				s = model.SearchRecord{}
				w.SearchRecord = s
				w.Window.Hide()
			}
		}, w.Window)
	})

	cancle := widget.NewButton("取消操作", func() {
		rt.Action()
		dialog.ShowConfirm("放弃？", "你确定要放弃你的操作码？", func(b bool) {
			rt.Action()
			if b {
				WinClose(w.Window)
				w.Window = nil
			}
		}, w.Window)
	})

	downBox := container.NewHBox(save, clearAll, cancle)
	downCenterBox := container.NewCenter(downBox)

	gg := NewBg(5, 30)

	box := container.NewVBox(upBox, gg, downCenterBox)
	cbox := container.NewCenter(box)

	bg := NewBg(fmax(cbox.MinSize().Width, cbox.Size().Width, 600),
		fmax(cbox.MinSize().Height, cbox.Size().Height, 350))

	lastContainer := container.NewStack(bg, cbox)
	w.Window.SetContent(lastContainer)
}

func (w *SaveWhereWindow) Show(rt runtime.RunTime) {
	w.create(rt)
	w.Window.Show()
}

func (w *SaveWhereWindow) Close() {
	WinClose(w.Window)
	w.Window = nil
}

var entryList6 []*widget.Entry

func newEntry6(input *string) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetText(*input)

	entry.OnChanged = func(s string) {
		if entry.Validate() == nil {
			*input = s
		}
	}

	entryList6 = append(entryList6, entry)
	return entry
}

func newFileIDEntry6(input *int64) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetText(fmt.Sprintf("%d", *input))

	entry.Validator = func(s string) error {
		if s == "" {
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

	if entry.Validate() != nil {
		entry.SetText("")
		*input = 0
	}

	entry.OnChanged = func(s string) {
		if entry.Validate() == nil {
			n, err := strconv.ParseInt(s, 10, 64)
			if err == nil {
				*input = n
			}
		}
	}

	entryList6 = append(entryList6, entry)
	return entry
}

func newSelect6(options []string, input *string) *widget.Select {
	const defaultStatus = "不约束"

	func() {
		for _, option := range options {
			if option == defaultStatus {
				return
			}
		}
		options = append([]string{defaultStatus}, options...)
	}()

	if *input != "" {
		func() {
			for _, option := range options {
				if option == *input {
					return
				}
			}
			options = append(options, *input)
		}()
	} else {
		*input = defaultStatus
	}

	sel := widget.NewSelect(options, func(s string) {
		if s == defaultStatus {
			*input = ""
		} else {
			*input = s
		}
	})
	sel.SetSelected(*input)
	return sel
}

func checkAllInputRight6() error {
	for _, e := range entryList6 {
		err := e.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
func newTimePicker6(input *sql.NullTime, w fyne.Window) *widget.Button {
	btn := widget.NewButton("选择时间", func() {})

	t := time.Now()
	if input.Valid {
		t = input.Time
	}

	d := datepicker.NewDateTimePicker(t, time.Monday, func(t time.Time, b bool) {
		if b {
			input.Valid = true
			input.Time = t
			btn.SetText(t.Format("2006-01-02 15:04:05"))
		} else {
			input.Valid = false
			btn.SetText("选择时间")
		}
	})

	btn.OnTapped = func() {
		dialog.ShowCustomConfirm("选择时间", "确认", "放弃", d, d.OnActioned, w)
	}

	return btn
}
func newBoolCheck6(label string, defaultChoice bool, input *bool) *widget.Check {
	c := widget.NewCheck(label, func(b bool) {
		*input = b
	})
	c.Checked = defaultChoice
	return c
}
