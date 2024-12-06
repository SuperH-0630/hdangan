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
	"gorm.io/gorm"
	"time"
)

func ShowNewMove(rt runtime.RunTime, fc model.File, refresh func(rt runtime.RunTime)) {
	config, err := systeminit.GetInit()
	if errors.Is(err, systeminit.LuckyError) {
		rt.DBConnectError(err)
		return
	} else if err != nil {
		rt.DBConnectError(fmt.Errorf("配置文件错误，请检查配置文件状态。"))
		return
	}

	f := fc.GetFile()

	infoWindow := rt.App().NewWindow(fmt.Sprintf("新增出借记录-%s-%d", f.Name, f.FileID))

	infoWindow.SetOnClosed(func() {
		rt.Action()
		infoWindow = nil
	})
	infoWindow.SetCloseIntercept(func() {
		rt.Action()
		WinClose(infoWindow)
		infoWindow = nil
	})

	oldRecord, err := model.FindMoveRecord(rt, fc)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		infoWindow.Resize(fyne.NewSize(300, 200))
		d := dialog.NewInformation("提示", "当前档案不存在出借记录", infoWindow)
		d.SetOnClosed(func() {
			infoWindow.Close()
		})
		d.Show()
		infoWindow.Show()
		return
	} else if err != nil {
		infoWindow.Resize(fyne.NewSize(300, 200))
		d := dialog.NewError(fmt.Errorf("数据错误: %s", err.Error()), infoWindow)
		d.SetOnClosed(func() {
			infoWindow.Close()
		})
		d.Show()
		infoWindow.Show()
		return
	}

	var newRecord model.FileMoveRecord

	moveComment := widget.NewMultiLineEntry()
	moveComment.Text = ""
	moveComment.Wrapping = fyne.TextWrapWord
	moveComment.OnChanged = func(text string) {
		newRecord.MoveComment = sql.NullString{String: text, Valid: len(text) != 0}
	}

	leftLayout := layout.NewFormLayout()
	left := container.New(leftLayout,
		widget.NewLabel("出借状态："),
		newMoveStatusSelect7(config.Yaml.Move.MoveStatus, &newRecord.MoveStatus),

		widget.NewLabel("借入人："),
		newEntryWithNULL7(&newRecord.MoveInPeopleName),

		widget.NewLabel("借入单位："),
		newMoveUnitSelectWithNULL7(config.Yaml.Move.MoveUnit, &newRecord.MoveInPeopleUnit),

		widget.NewLabel("借出人："),
		newEntryWithNULL7(&newRecord.MoveOutPeopleName),

		widget.NewLabel("借出单位："),
		newMoveUnitSelectWithNULL7(config.Yaml.Move.MoveUnit, &newRecord.MoveOutPeopleUnit),

		widget.NewLabel("借出单位："),
		newTimePicker7(time.Now(), &newRecord.MoveTime, infoWindow),

		widget.NewLabel("最后迁出备注："),
		moveComment,
	)

	upBox := container.NewHBox(left)

	save := widget.NewButton("保存", func() {
		rt.Action()
		err := checkAllInputRight7()
		if err != nil {
			dialog.ShowError(fmt.Errorf("请检查错误：%s", err.Error()), infoWindow)
			return
		}

		if len(newRecord.MoveStatus) == 0 {
			dialog.ShowError(fmt.Errorf("必须设置状态"), infoWindow)
			return
		}

		if newRecord.MoveTime.Before(oldRecord.MoveTime) {
			dialog.ShowError(fmt.Errorf("借出时间不得早于上次借出的事件"), infoWindow)
			return
		}

		if newRecord.MoveTime.After(time.Now().Add(24 * time.Hour)) {
			dialog.ShowError(fmt.Errorf("借出时间不得晚于当下的24小时之后"), infoWindow)
			return
		}

		err = model.CreateFileRecord(rt, f, &newRecord)
		if err != nil {
			dialog.ShowError(fmt.Errorf("数据库错误：%s", err.Error()), infoWindow)
			return
		}

		WinClose(infoWindow)
		infoWindow = nil
		refresh(rt)
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
		fmax(cbox.MinSize().Height, cbox.Size().Height, 280))

	lastContainer := container.NewStack(bg, cbox)
	infoWindow.SetContent(lastContainer)

	infoWindow.Show()
	infoWindow.CenterOnScreen()
}

var entryList7 []*widget.Entry

func newEntryWithNULL7(input *sql.NullString) *widget.Entry {
	entry := widget.NewEntry()

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

	entryList7 = append(entryList7, entry)
	return entry
}

func newMoveStatusSelect7(options []string, input *string) *widget.Select {
	sel := widget.NewSelect(options, func(s string) {
		*input = s
	})
	return sel
}

func newMoveUnitSelectWithNULL7(firstOptions []string, input *sql.NullString) *widget.Select {
	const emptySelectItem = "暂无"

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
	return sel
}

func newTimePicker7(t time.Time, input *time.Time, w fyne.Window) *widget.Button {
	btn := widget.NewButton("选择时间", func() {})

	d := datepicker.NewDateTimePicker(t, time.Monday, func(t time.Time, b bool) {
		if b {
			*input = t
			btn.SetText(t.Format("2006-01-02 15:04:05"))
		}
	})

	btn.OnTapped = func() {
		dialog.ShowCustomConfirm("选择时间", "确认", "放弃", d, d.OnActioned, w)
	}

	return btn
}

func checkAllInputRight7() error {
	for _, e := range entryList7 {
		err := e.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
