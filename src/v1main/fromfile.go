package v1main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/SuperH-0630/hdangan/src/excelreader"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"os"
)

func AddFromFile(rt runtime.RunTime, w fyne.Window, refresh func(rt runtime.RunTime)) error {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if reader == nil {
			return
		} else if err != nil {
			dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w)
			return
		}

		defer func() {
			_ = reader.Close()
		}()

		filename := reader.URI().Path()
		if !IsFileExists(filename) {
			dialog.ShowError(fmt.Errorf("文件不存在：%s", filename), w)
			return
		}

		sa, su, fu, err := excelreader.ReadFile(rt, reader)
		if errors.Is(err, excelreader.BadTitle) {
			dialog.ShowError(fmt.Errorf("表格首行（表头）对应错误"), w)
			return
		} else if err != nil {
			dialog.ShowError(fmt.Errorf("导入出错：%s", err), w)
			return
		}

		dialog.ShowInformation("完成", fmt.Sprintf("恭喜你！共新增%d条数据，共升级%d条数据，但是%d条数据未成功识别。。", sa, su, fu), w)
		refresh(rt)
	}, w)

	return nil
}

func IsFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}

	return false
}
