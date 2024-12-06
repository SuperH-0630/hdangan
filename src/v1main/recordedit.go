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
	"gorm.io/gorm"
)

func ShowMoveEdit(rt runtime.RunTime, fc model.File, refresh func(rt runtime.RunTime)) {
	config, err := systeminit.GetInit()
	if errors.Is(err, systeminit.LuckyError) {
		rt.DBConnectError(err)
		return
	} else if err != nil {
		rt.DBConnectError(fmt.Errorf("配置文件错误，请检查配置文件状态。"))
		return
	}

	f := fc.GetFile()

	infoWindow := rt.App().NewWindow(fmt.Sprintf("出借记录-%s-%d", f.Name, f.FileID))

	infoWindow.SetOnClosed(func() {
		rt.Action()
		infoWindow = nil
	})
	infoWindow.SetCloseIntercept(func() {
		rt.Action()
		WinClose(infoWindow)
		infoWindow = nil
	})

	record, err := model.FindMoveRecord(rt, fc)
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

	moveComment := widget.NewMultiLineEntry()
	moveComment.Text = strToStr(record.MoveComment, "")
	moveComment.Wrapping = fyne.TextWrapWord
	moveComment.OnChanged = func(text string) {
		record.MoveComment = sql.NullString{String: text, Valid: len(text) != 0}
	}

	leftLayout := layout.NewFormLayout()
	left := container.New(leftLayout,
		widget.NewLabel("出借状态："),
		newMoveStatusSelect2(fmt.Sprintf("%s", record.MoveStatus), config.Yaml.Move.MoveStatus, &record.MoveStatus),

		widget.NewLabel("借入人："),
		newEntryWithNULL2(fmt.Sprintf("%s", strToStr(record.MoveInPeopleName, "")), &record.MoveInPeopleName),

		widget.NewLabel("借入单位："),
		newMoveUnitSelectWithNULL2(fmt.Sprintf("%s", strToStr(record.MoveInPeopleUnit, "")), config.Yaml.Move.MoveUnit, &record.MoveInPeopleUnit),

		widget.NewLabel("借出人："),
		newEntryWithNULL2(fmt.Sprintf("%s", strToStr(record.MoveOutPeopleName, "")), &record.MoveOutPeopleName),

		widget.NewLabel("借出单位："),
		newMoveUnitSelectWithNULL2(fmt.Sprintf("%s", strToStr(record.MoveOutPeopleUnit, "")), config.Yaml.Move.MoveUnit, &record.MoveOutPeopleUnit),

		widget.NewLabel("最后迁出备注："),
		moveComment,
	)

	upBox := container.NewHBox(left)

	save := widget.NewButton("保存", func() {
		rt.Action()
		err := checkAllInputRight2()
		if err != nil {
			dialog.ShowError(fmt.Errorf("请检查错误：%s", err.Error()), infoWindow)
			return
		}

		err = model.SaveRecord(rt, record)
		if err != nil {
			dialog.ShowError(fmt.Errorf("数据库错误: %s", err.Error()), infoWindow)
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

var entryList2 []*widget.Entry

func newEntryWithNULL2(data string, input *sql.NullString) *widget.Entry {
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

	entryList2 = append(entryList2, entry)
	return entry
}

func newMoveStatusSelect2(data string, options []string, input *string) *widget.Select {
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

func newMoveUnitSelectWithNULL2(data string, firstOptions []string, input *sql.NullString) *widget.Select {
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
	return sel
}

func checkAllInputRight2() error {
	for _, e := range entryList2 {
		err := e.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
