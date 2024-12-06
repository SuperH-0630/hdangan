package v1main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
)

type RecordWindow struct {
	Window         fyne.Window
	Menu           *Controller
	Table          *RecordTable
	FileWindow     fyne.Window
	File           model.File
	NowPage        int64
	InfoRecord     []model.FileMoveRecord
	InfoDataRecord [][]string
	SearchRecord   model.SearchRecord
}

func CreateRecordWindow(rt runtime.RunTime, fc model.File, fileWindow fyne.Window) *RecordWindow {
	f := fc.GetFile()

	recordWindow := rt.App().NewWindow(fmt.Sprintf("档案出借记录-%d", f.FileID))
	w := &RecordWindow{
		Window: recordWindow,

		Menu:  nil,
		Table: nil,

		FileWindow:     fileWindow,
		File:           f,
		NowPage:        1,
		InfoRecord:     make([]model.FileMoveRecord, 0, 0),
		InfoDataRecord: make([][]string, 0, 0),
	}

	CreateRecordTable(rt, w)
	GetMainMenuRecord(rt, w, func(rt runtime.RunTime) {
		w.Table.UpdateTableRecord(rt, 0, w.NowPage)
	})

	bg := NewBg(fmax(w.Table.Table.MinSize().Width, w.Table.Table.Size().Width, 600),
		fmax(w.Table.Table.MinSize().Height, w.Table.Table.Size().Height, 400))

	lastContainer := container.NewStack(bg, w.Table.Table)
	w.Window.SetContent(lastContainer)

	w.Window.SetOnClosed(func() {
		rt.Action()
		w.Window = nil
	})

	w.Window.SetCloseIntercept(func() {
		rt.Action()
		WinClose(w.Window)
		w.Window = nil
	})

	w.Table.FirstUpdateData(rt)
	w.Window.SetFixedSize(true)
	return w
}

func (w *RecordWindow) Show() {
	w.Window.Show()
}

func (w *RecordWindow) CenterOnScreen() {
	w.Window.CenterOnScreen()
}
