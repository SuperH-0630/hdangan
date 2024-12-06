package v1main

import (
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
)

func TongJi(rt runtime.RunTime, w *CtrlWindow) {
	fca := make([]int64, len(model.FileSetTypeList))
	sum := int64(0)

	for _, k := range model.FileSetTypeList {
		var err error
		maker, ok := model.FileSetTypeMaker[k]
		if !ok {
			fca[k] = 0
			continue
		}

		fca[k], err = model.CountFile(rt, maker())
		if err != nil {
			fca[k] = 0
			continue
		}

		sum += fca[k]
	}

	m1 := fmt.Sprintf("数据库共记录档案：%d件。", sum)
	m2 := ""
	if sum != 0 {
		m2 = "其中"
		for _, i := range model.FileSetTypeList {
			ca := fca[i]
			if ca <= 0 {
				continue
			}

			m2 += fmt.Sprintf("，%s共%d件", model.FileSetTypeName[i], ca)
		}
		m2 += "。"
	}

	m := fmt.Sprintf("%s\n%s", m1, m2)

	dialog.ShowInformation("数据统计", m, w.window)
}
