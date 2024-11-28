package fail

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/SuperH-0630/hdangan/src/systeminit"
	"os"
)

func ToFail(res string) {
	a := app.New()

	w1 := a.NewWindow("Error")

	w1.Resize(fyne.NewSize(300, 300))
	w1.SetMaster()

	l1 := widget.NewLabel("")
	w1.SetContent(l1)

	var msg string
	d, err := systeminit.GetInit()
	if errors.Is(err, systeminit.LuckyError) {
		msg = err.Error()
	} else if err == nil {
		msg = fmt.Sprintf("I am sorry, that we has miss some wrong.\n%s\n\n\n寻求帮助：%s <%s>", res, d.Yaml.Report.Name, d.Yaml.Report.Email)
	} else {
		msg = fmt.Sprintf("I am sorry, that we has miss some wrong.\n%s", res)
	}

	l1.SetText(msg)

	w1.Show()
	a.Run()

	os.Exit(1)
}
