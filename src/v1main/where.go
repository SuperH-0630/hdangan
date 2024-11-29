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

var whereWindow fyne.Window

func createWindow(rt runtime.RunTime, target *model.SearchWhere, refresh func(rt runtime.RunTime)) {
	config, err := systeminit.GetInit()
	if errors.Is(err, systeminit.LuckyError) {
		rt.DBConnectError(err)
		return
	} else if err != nil {
		rt.DBConnectError(fmt.Errorf("配置文件错误，请检查配置文件状态。"))
		return
	}

	var s = *target
	whereWindow = rt.App().NewWindow(fmt.Sprintf("搜索筛选"))

	whereWindow.SetOnClosed(func() {
		rt.Action()
		whereWindow = nil
	})
	whereWindow.SetCloseIntercept(func() {
		rt.Action()
		WinClose(whereWindow)
		whereWindow = nil
	})

	leftLayout := layout.NewFormLayout()
	left := container.New(leftLayout,
		widget.NewLabel("卷宗号："),
		newFileIDEntry3(&s.FileID),

		widget.NewLabel("姓名："),
		newEntry3(&s.Name),

		widget.NewLabel("身份证号："),
		newEntry3(&s.IDCard),

		widget.NewLabel("户籍地："),
		newEntry3(&s.Location),

		widget.NewLabel("卷宗标题："),
		newEntry3(&s.FileTitle),

		widget.NewLabel("卷宗类型："),
		newSelect3(config.Yaml.File.FileType, &s.FileType),
	)

	centerLayout := layout.NewFormLayout()
	center := container.New(centerLayout,
		widget.NewLabel("最早第一次迁入时间（含）："),
		newTimePicker3(&s.FirstMoveInStart, whereWindow),

		widget.NewLabel("最晚第一次迁入时间（含）："),
		newTimePicker3(&s.FirstMoveInStart, whereWindow),

		widget.NewLabel("最早最后一次迁入（归还）时间（含）："),
		newTimePicker3(&s.LastMoveInStart, whereWindow),

		widget.NewLabel("最晚最后一次迁入（归还）时间（含）："),
		newTimePicker3(&s.LastMoveInEnd, whereWindow),

		widget.NewLabel("最早最后一次迁出时间（含）："),
		newTimePicker3(&s.LastMoveOutStart, whereWindow),

		widget.NewLabel("最晚最后一次迁出时间（含）："),
		newTimePicker3(&s.LastMoveOutEnd, whereWindow),
	)

	rightLayout := layout.NewFormLayout()
	right := container.New(rightLayout,
		widget.NewLabel("迁入迁出状态："),
		newSelect3(config.Yaml.Move.MoveStatus, &s.MoveStatus),

		widget.NewLabel("最后迁出人："),
		newEntry3(&s.MoveOutPeopleName),

		widget.NewLabel("最后迁出单位："),
		newSelect3(config.Yaml.Move.MoveUnit, &s.MoveOutPeopleUnit),
	)

	upBox := container.NewHBox(left, center, right)

	save := widget.NewButton("保存条件", func() {
		rt.Action()
		err := checkAllInputRight3()
		if err != nil {
			dialog.ShowError(fmt.Errorf("请检查错误：%s", err.Error()), whereWindow)
			return
		}
		*target = s
		whereWindow.Hide()
		refresh(rt)
	})

	clearAll := widget.NewButton("清空条件", func() {
		rt.Action()
		dialog.ShowConfirm("确定？", "你是否确定要清空全部条件？？", func(b bool) {
			rt.Action()
			if b {
				s = model.SearchWhere{}
				*target = s
				whereWindow.Hide()
				refresh(rt)
			}
		}, whereWindow)
	})

	cancle := widget.NewButton("取消操作", func() {
		rt.Action()
		dialog.ShowConfirm("放弃？", "你确定要放弃你的操作码？", func(b bool) {
			rt.Action()
			if b {
				whereWindow.Hide()
			}
		}, whereWindow)
	})

	downBox := container.NewHBox(save, clearAll, cancle)
	downCenterBox := container.NewCenter(downBox)

	gg := NewBg(5, 30)

	box := container.NewVBox(upBox, gg, downCenterBox)
	cbox := container.NewCenter(box)

	bg := NewBg(fmax(cbox.MinSize().Width, cbox.Size().Width, 600),
		fmax(cbox.MinSize().Height, cbox.Size().Height, 350))

	lastContainer := container.NewStack(bg, cbox)
	whereWindow.SetContent(lastContainer)
}

func ShowWhereWindow(rt runtime.RunTime, s *model.SearchWhere, refresh func(rt runtime.RunTime)) {
	createWindow(rt, s, refresh)
	whereWindow.Show()
}

func HideWhereWindow() {
	if whereWindow != nil {
		whereWindow.Hide()
	}
}

var entryList3 []*widget.Entry

func newEntry3(input *string) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetText(*input)

	entry.OnChanged = func(s string) {
		if entry.Validate() == nil {
			*input = s
		}
	}

	entryList3 = append(entryList3, entry)
	return entry
}

func newFileIDEntry3(input *int64) *widget.Entry {
	entry := widget.NewEntry()
	if *input <= 0 {
		entry.Text = ""
	} else {
		entry.Text = fmt.Sprintf("%d", *input)
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

func newSelect3(options []string, input *string) *widget.Select {
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

func checkAllInputRight3() error {
	for _, e := range entryList3 {
		err := e.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
func newTimePicker3(input *sql.NullTime, w fyne.Window) *widget.Button {
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
