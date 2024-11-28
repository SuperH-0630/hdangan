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
	"github.com/SuperH-0630/hdangan/src/fail"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"github.com/SuperH-0630/hdangan/src/systeminit"
	"time"
)

func ShowMove(rt runtime.RunTime, f *model.File, refresh func(rt runtime.RunTime)) {
	config, err := systeminit.GetInit()
	if errors.Is(err, systeminit.LuckyError) {
		fail.ToFail(err.Error())
		return
	}

	infoWindow := rt.App().NewWindow(fmt.Sprintf("迁入迁出-%s-%d", f.Name, f.FileID))

	infoWindow.SetOnClosed(func() {
		rt.Action()
		WinClose(infoWindow)
		infoWindow = nil
	})
	infoWindow.SetCloseIntercept(func() {
		rt.Action()
		WinClose(infoWindow)
		infoWindow = nil
	})

	moveComment := widget.NewMultiLineEntry()
	moveComment.Text = strToStr(f.MoveComment, "")
	moveComment.Wrapping = fyne.TextWrapWord
	moveComment.OnChanged = func(text string) {
		f.MoveComment = sql.NullString{String: text, Valid: len(text) != 0}
	}

	leftLayout := layout.NewFormLayout()
	left := container.New(leftLayout,
		widget.NewLabel("迁入迁出状态："),
		newMoveStatusSelect2(fmt.Sprintf("%s", f.MoveStatus), config.Yaml.Move.MoveStatus, &f.MoveStatus),

		widget.NewLabel("最后迁出人："),
		newEntryWithNULL2(fmt.Sprintf("%s", strToStr(f.MoveOutPeopleName, "")), &f.MoveOutPeopleName),

		widget.NewLabel("最后迁出单位："),
		newMoveUnitSelectWithNULL2(fmt.Sprintf("%s", strToStr(f.MoveOutPeopleUnit, "")), config.Yaml.Move.MoveUnit, &f.MoveOutPeopleUnit),

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

		if f.MoveStatus == config.Yaml.Move.MoveInStatus {
			// 正常非迁入之间更改
			if f.MoveComment.Valid || f.MoveOutPeopleName.Valid || f.MoveOutPeopleUnit.Valid {
				// 非借出状态，确有借出数据
				dialog.ShowError(fmt.Errorf("非迁出状态却有迁出人等数据"), infoWindow)
			} else {
				oldRecord, _ := model.CheckFileMoveOut(f)

				record := &model.FileMoveRecord{
					// FileID会在执行SaveFile时绑定
					MoveStatus:        f.MoveStatus,
					MoveTime:          time.Now(),
					MoveOutPeopleName: f.MoveOutPeopleName,
					MoveOutPeopleUnit: f.MoveOutPeopleUnit,
					MoveComment:       f.MoveComment,
				}

				f.LastMoveIn = record.MoveTime

				err := model.SaveFileRecord(f, record, oldRecord)
				if err != nil {
					dialog.ShowError(fmt.Errorf("数据库错误: %s", err.Error()), infoWindow)
				}
			}
		} else {
			if !f.MoveOutPeopleName.Valid || !f.MoveOutPeopleUnit.Valid {
				dialog.ShowError(fmt.Errorf("迁出状态需要填写迁出人和单位"), infoWindow)
			} else {
				oldRecord, _ := model.CheckFileMoveOut(f)

				record := &model.FileMoveRecord{
					MoveStatus:        f.MoveStatus,
					MoveTime:          time.Now(),
					MoveOutPeopleName: f.MoveOutPeopleName,
					MoveOutPeopleUnit: f.MoveOutPeopleUnit,
					MoveComment:       f.MoveComment,
				}

				f.LastMoveIn = record.MoveTime

				err := model.SaveFileRecord(f, record, oldRecord)
				if err != nil {
					dialog.ShowError(fmt.Errorf("数据库错误: %s", err.Error()), infoWindow)
				}
			}
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

	bg := NewBg(max(cbox.MinSize().Width, cbox.Size().Width, 220),
		max(cbox.MinSize().Height, cbox.Size().Height, 280))

	lastContainer := container.NewStack(bg, cbox)
	infoWindow.SetContent(lastContainer)

	infoWindow.Show()
	infoWindow.CenterOnScreen()
}

var entryList2 []*widget.Entry

func newEntry2(data string, input *string) *widget.Entry {
	entry := widget.NewEntry()
	entry.Text = data

	entry.OnChanged = func(s string) {
		if entry.Validate() == nil {
			*input = s
		}
	}

	entryList2 = append(entryList2, entry)
	return entry
}

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
