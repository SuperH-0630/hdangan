package v1main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
)

func TongJi(rt runtime.RunTime, w fyne.Window) {
	ca, err := model.CountFile(rt)
	if err != nil {
		dialog.ShowError(fmt.Errorf("数据库错误：%s", err.Error()), w)
		return
	}

	fca, err := model.CountDifferentFile(rt)
	if err != nil {
		dialog.ShowError(fmt.Errorf("数据库错误：%s", err.Error()), w)
		return
	}

	m1 := fmt.Sprintf("数据库共记录档案：%d件。", ca)
	m2 := "其中"
	for _, i := range fca {
		m2 += fmt.Sprintf("，%s共%d件", i.File, i.Res)
	}
	m2 += "。"

	m := fmt.Sprintf("%s\n%s", m1, m2)

	dialog.ShowInformation("数据统计", m, w)
}
