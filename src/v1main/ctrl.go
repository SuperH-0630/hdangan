package v1main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"time"
)

type CtrlWindow struct {
	window   fyne.Window
	lastUse  time.Time
	killTime time.Duration
	rt       runtime.RunTime
	menu     *MainMenu
	table    *MainTable
}

var ctrlWindow *CtrlWindow

func newCtrlWindow(rt runtime.RunTime, title string, killSecond time.Duration) *CtrlWindow {
	ks := killSecond * time.Second
	cw := &CtrlWindow{
		window:   rt.App().NewWindow(title),
		lastUse:  time.Now(),
		killTime: ks,
		rt:       rt,
	}

	rt.SetDBConnectErrorWindow(cw.window)
	rt.SetAction(func() {
		cw.Action()
	})

	cw.menu = getMainMenu(rt, cw, func(rt runtime.RunTime) {
		cw.table.UpdateTable(rt, cw.table.fileSetType, 0, cw.menu.NowPage)
	})
	cw.window.SetMainMenu(cw.menu.Main)

	cw.table = CreateInfoTable(rt, cw)

	return cw
}

func (w *CtrlWindow) Show() {
	w.window.Show()
	w.lastUse = time.Now()
	w.cc()
}

func (w *CtrlWindow) Close() {
	WinClose(w.window)
	w.window = nil
}

func (w *CtrlWindow) Hide() {
	w.window.Hide()
}

func (w *CtrlWindow) Action() {
	if w == nil || w.window == nil {
		return
	}
	w.lastUse = time.Now()
}

func (w *CtrlWindow) cc() {
	go func(rt runtime.RunTime, w *CtrlWindow) {
		for range time.Tick(time.Second) {
			if w.window != nil {
				return
			}

			if time.Now().Sub(w.lastUse) > w.killTime {
				err := HideCtrlWindow(rt) // 强行关闭
				if err == nil {
					ShowHelloWindowTimeout(rt)
				}
				return
			}
		}
	}(w.rt, w)
}

func createCtrlWindow(rt runtime.RunTime) error {
	if ctrlWindow != nil {
		return nil
	}

	ctrlWindow = newCtrlWindow(rt, "桓档案-控制中心", 15*60)

	ctrlWindow.window.SetOnClosed(func() {
		ctrlWindow = nil
		ShowHelloWindow(rt)
	})

	ctrlWindow.window.SetCloseIntercept(func() {
		WinClose(ctrlWindow.window)
		ctrlWindow = nil
		ShowHelloWindow(rt)
	})

	bg := NewBg(fmax(ctrlWindow.table.fileTable.MinSize().Width, ctrlWindow.table.fileTable.Size().Width, 800),
		fmax(ctrlWindow.table.fileTable.MinSize().Height, ctrlWindow.table.fileTable.Size().Height, 500))

	lastContainer := container.NewStack(bg, ctrlWindow.table.fileTable)
	ctrlWindow.window.SetContent(lastContainer)

	return nil
}

func ShowCtrlWindow(rt runtime.RunTime) error {
	err := createCtrlWindow(rt)
	if err != nil {
		return err
	}
	ctrlWindow.Show()
	ctrlWindow.window.CenterOnScreen()
	return err
}

func HideCtrlWindow(rt runtime.RunTime) error {
	err := createCtrlWindow(rt)
	if err != nil {
		return err
	}
	ctrlWindow.Hide()
	return nil
}
