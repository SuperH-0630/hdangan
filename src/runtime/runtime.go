package runtime

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type RunTime interface {
	App() fyne.App
	Action()
	SetAction(func())
	DBConnectError(err error)
	SetDBConnectErrorWindow(window fyne.Window)
	SetGameStopFunc(stop func())
	StopGame()
}

type runTime struct {
	app                  fyne.App
	action               func()
	dbConnectErrorWindow fyne.Window
	stopFunc             func()
}

func NewRunTime(app fyne.App) RunTime {
	return &runTime{app: app, action: func() {}, dbConnectErrorWindow: nil}
}

func (rt *runTime) App() fyne.App {
	return rt.app
}

func (rt *runTime) Action() {
	if rt.action != nil {
		rt.action()
	}
}

func (rt *runTime) SetAction(action func()) {
	rt.action = action
}

func (rt *runTime) DBConnectError(err error) {
	if rt.dbConnectErrorWindow != nil {
		dialog.ShowError(err, rt.dbConnectErrorWindow)
	}
}

func (rt *runTime) SetDBConnectErrorWindow(window fyne.Window) {
	rt.dbConnectErrorWindow = window
}

func (rt *runTime) SetGameStopFunc(stop func()) {
	rt.stopFunc = stop
}

func (rt *runTime) StopGame() {
	if rt.stopFunc != nil {
		rt.stopFunc()
		rt.stopFunc = nil
	}
}
