package v1main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	happ "github.com/SuperH-0630/hdangan/src/app"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"github.com/SuperH-0630/hdangan/src/systeminit"
	"os"
)

func Main() {
	start()

	a := happ.NewApp()

	rt := runtime.NewRunTime(a)

	err := model.AutoCreateModel(rt)
	if err != nil {
		dbFail(rt, fmt.Sprintf("数据库构建失败: %s。", err.Error()), 1)
		return
	}

	ShowHelloWindow(rt)
	StartTheGame(rt)

	a.Run()

	exit()
}

func start() {
	fmt.Println("START")
}

func exit() {
	fmt.Println("EXIT")
}

func dbFail(rt runtime.RunTime, res string, exitCode int) {
	defer func() {
		if exitCode <= 0 {
			exitCode = 1
		}

		os.Exit(exitCode)
	}()

	w1 := rt.App().NewWindow("数据库错误")

	w1.Resize(fyne.NewSize(300, 300))

	var msg string
	d, err := systeminit.GetInit()
	if errors.Is(err, systeminit.LuckyError) {
		msg = fmt.Sprintf("I am sorry, that we has miss some wrong.\n%s\n\n\n寻求帮助：怨天由人", res)
	} else if err == nil {
		msg = fmt.Sprintf("I am sorry, that we has miss some wrong.\n%s\n\n\n寻求帮助：%s <%s>", res, d.Yaml.Report.Name, d.Yaml.Report.Email)
	} else {
		msg = fmt.Sprintf("I am sorry, that we has miss some wrong.\n%s", res)
	}

	dialog.ShowError(fmt.Errorf("%s", msg), w1)

	w1.SetContent(widget.NewLabel(""))
	w1.CenterOnScreen()
	w1.SetFixedSize(true)
	w1.Show()
}
