package happ

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/SuperH-0630/hdangan/src/assest"
)

type App struct {
	fyne.App
}

func NewApp() *App {
	res := &App{
		App: app.NewWithID("com.song-zh.hdangan"),
	}

	res.App.SetIcon(assest.MainIco)
	return res
}

func (app *App) NewWindow(title string) fyne.Window {
	w := app.App.NewWindow(title)
	w.SetIcon(assest.MainIco)
	return w
}
