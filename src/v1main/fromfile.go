package v1main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/SuperH-0630/hdangan/src/excelio"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"os"
)

func AddFromFile(rt runtime.RunTime, w *CtrlWindow, refresh func(rt runtime.RunTime)) error {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if reader == nil {
			return
		} else if err != nil {
			dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w.window)
			return
		}

		defer func() {
			_ = reader.Close()
		}()

		filename := reader.URI().Path()
		if !IsFileExists(filename) {
			dialog.ShowError(fmt.Errorf("文件不存在：%s", filename), w.window)
			return
		}

		sa, su, fu, err := excelio.ReadFile(rt, w.fileSetType, reader)
		if errors.Is(err, excelio.BadTitle) {
			dialog.ShowError(fmt.Errorf("表格首行（表头）对应错误"), w.window)
			return
		} else if err != nil {
			dialog.ShowError(fmt.Errorf("导入出错：%s", err), w.window)
			return
		}

		dialog.ShowInformation("完成", fmt.Sprintf("恭喜你！共新增%d条数据，共升级%d条数据，但是%d条数据未成功识别。。", sa, su, fu), w.window)
		refresh(rt)
	}, w.window)

	return nil
}

func AddRecordFromFile(rt runtime.RunTime, w *CtrlWindow, refresh func(rt runtime.RunTime)) error {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if reader == nil {
			return
		} else if err != nil {
			dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w.window)
			return
		}

		defer func() {
			_ = reader.Close()
		}()

		filename := reader.URI().Path()
		if !IsFileExists(filename) {
			dialog.ShowError(fmt.Errorf("文件不存在：%s", filename), w.window)
			return
		}

		sa, su, fu, err := excelio.ReadRecord(rt, w.fileSetType, reader)
		if errors.Is(err, excelio.BadTitle) {
			dialog.ShowError(fmt.Errorf("表格首行（表头）对应错误"), w.window)
			return
		} else if err != nil {
			dialog.ShowError(fmt.Errorf("导入出错：%s", err), w.window)
			return
		}

		dialog.ShowInformation("完成", fmt.Sprintf("恭喜你！共新增%d条数据，共升级%d条数据，但是%d条数据未成功识别。。", sa, su, fu), w.window)
		refresh(rt)
	}, w.window)

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
