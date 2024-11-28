package v1main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/SuperH-0630/hdangan/src/fail"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"time"
)

type window fyne.Window

type CtrlWindow struct {
	window
	lastUse  time.Time
	killTime time.Duration
	rt       runtime.RunTime
}

func NewWindow(rt runtime.RunTime, title string, killSecond time.Duration) *CtrlWindow {
	ks := killSecond * time.Second
	cw := &CtrlWindow{
		window:   rt.App().NewWindow(title),
		lastUse:  time.Now(),
		killTime: ks,
		rt:       rt,
	}

	cw.window.SetOnClosed(func() {
		rt.Action()
		WinClose(cw.window)
	})

	cw.window.SetCloseIntercept(func() {
		rt.Action()
		WinClose(cw.window)
	})

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
	if w.window == nil {
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
				ShowHelloWindowTimeout(rt)
				HideCtrlWindow(rt)
				return
			}
		}
	}(w.rt, w)
}

var ctrlWindow *CtrlWindow

func onClose(rt runtime.RunTime) {
	if ctrlWindow == nil {
		return
	}

	WinClose(ctrlWindow)
	ctrlWindow = nil
	ShowHelloWindow(rt)
}

func createCtrlWindow(rt runtime.RunTime) error {
	if ctrlWindow != nil {
		return nil
	}

	ctrlWindow = NewWindow(rt, "桓档案-控制中心", 15*60)
	rt.SetDBConnectErrorWindow(ctrlWindow)
	rt.SetAction(func() {
		ctrlWindow.Action()
	})

	ctrlWindow.SetMainMenu(getMainMenu(rt, ctrlWindow, func(rt runtime.RunTime) {
		UpdateTable(rt, ctrlWindow, fileTable, 0, NowPage)
	}))

	CreateInfoTable(rt, ctrlWindow)

	bg := NewBg(max(fileTable.MinSize().Width, fileTable.Size().Width, 800),
		max(fileTable.MinSize().Height, fileTable.Size().Height, 500))

	lastContainer := container.NewStack(bg, fileTable)
	ctrlWindow.SetContent(lastContainer)

	ctrlWindow.SetOnClosed(func() {
		onClose(rt)
	})
	ctrlWindow.SetCloseIntercept(func() {
		onClose(rt)
	})
	return nil
}

func ShowCtrlWindow(rt runtime.RunTime) {
	err := createCtrlWindow(rt)
	if err != nil {
		fail.ToFail("非常抱歉，数据加载失败。")
		return
	}
	ctrlWindow.Show()
	ctrlWindow.CenterOnScreen()
}

func HideCtrlWindow(rt runtime.RunTime) {
	err := createCtrlWindow(rt)
	if err != nil {
		fail.ToFail("非常抱歉，数据加载失败。")
		return
	}
	ctrlWindow.Hide()
}
