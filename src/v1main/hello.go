package v1main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/SuperH-0630/hdangan/src/assest"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"image/color"
)

var helloWindow fyne.Window

func startBtnClick(rt runtime.RunTime) error {
	err := ShowCtrlWindow(rt)
	if err != nil {
		dialog.ShowError(fmt.Errorf("无法启动: %s", err.Error()), helloWindow)
		return err
	}

	HideHelloWindow(rt)
	return nil
}

func createHelloWindow(rt runtime.RunTime) {
	fmt.Println("TAG A")
	if helloWindow != nil {
		return
	}

	helloWindow = rt.App().NewWindow("桓档案")

	pic := canvas.NewImageFromResource(assest.StartPic)

	welcomeStr := canvas.NewText("欢迎使用桓档案管理器！", color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	welcomeStr.TextSize = 20
	welcomeStr.TextStyle.Bold = true
	welcomeStr.FontSource = assest.XiaoheRegularFont
	welcomeStrContainer := container.NewCenter(welcomeStr)

	wbBottom := widget.NewLabel("")
	wbBottom.Resize(fyne.NewSize(1, 20))

	wboader := layout.NewBorderLayout(wbBottom, wbBottom, nil, nil)
	welcomeStrContainerWithSpace := container.New(wboader, welcomeStrContainer)

	startBtn := widget.NewButton("立刻进入", func() {
		_ = startBtnClick(rt)
	})
	startBtn.SetIcon(assest.MainIco)
	startBtn.Resize(fyne.NewSize(200, startBtn.MinSize().Height))
	startBtnContainer := container.NewCenter(startBtn)

	vbox := layout.NewVBoxLayout()
	welcomeStrWithBotton := container.New(vbox, welcomeStrContainerWithSpace, startBtnContainer)

	welcomeStrWithBottonAsCenter := container.NewCenter(welcomeStrWithBotton)

	bTop := widget.NewLabel("")
	bTop.Resize(fyne.NewSize(1, 20))
	bBottom := widget.NewLabel("")
	bBottom.Resize(fyne.NewSize(1, 20))
	boader := layout.NewBorderLayout(bTop, bBottom, nil, nil)
	welcomeStrWithBottonWithBoader := container.New(boader, welcomeStrWithBottonAsCenter)

	bg := NewBg(fmax(welcomeStrWithBottonWithBoader.MinSize().Width, welcomeStrWithBottonWithBoader.Size().Width, 350),
		fmax(welcomeStrWithBottonWithBoader.MinSize().Height, welcomeStrWithBottonWithBoader.Size().Height, 200))

	lowerContainer := container.NewStack(pic, bg, welcomeStrWithBottonWithBoader)

	helloWindow.SetContent(lowerContainer)
	helloWindow.SetMaster()
	helloWindow.SetFixedSize(true)

	helloWindow.SetOnClosed(func() {
		rt.App().Quit()
	})

	helloWindow.SetCloseIntercept(func() {
		helloWindow.Close()
		rt.App().Quit()
	})
	fmt.Println("TAG B")
}

func ShowHelloWindow(rt runtime.RunTime) {
	createHelloWindow(rt)
	helloWindow.Show()
	helloWindow.CenterOnScreen()
}

func ShowHelloWindowTimeout(rt runtime.RunTime) {
	ShowHelloWindow(rt)
	dialog.ShowInformation("超时未操作", "您已超过一定时间未操作本系统。", helloWindow)
}

func HideHelloWindow(rt runtime.RunTime) {
	createHelloWindow(rt)
	helloWindow.Hide()
}
