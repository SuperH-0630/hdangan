package v1main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/SuperH-0630/hdangan/src/excelio"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"strings"
)

type Controller struct {
	Record *RecordWindow
	Menu   *fyne.MainMenu
	Page   *fyne.Menu
	Search *WhereWindow
}

func GetMainMenuRecord(rt runtime.RunTime, w *RecordWindow, refresh func(rt runtime.RunTime)) *Controller {
	show := fyne.NewMenuItem("查看档案", func() {
		rt.Action()
		w.Window.Show()
		w.Window.CenterOnScreen()
	})

	ss := NewWhereWindow(rt, w, refresh)

	search := fyne.NewMenuItem("搜索条件", func() {
		rt.Action()
		ss.Show(rt)
	})

	outputAll := fyne.NewMenuItem("导出全部迁移数据", func() {
		rt.Action()
		dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer == nil {
				return
			} else if err != nil {
				dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w.Window)
				return
			}

			defer func() {
				_ = writer.Close()
			}()

			savepath := writer.URI().Path()

			if !strings.HasSuffix(savepath, ".xlsx") {
				dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w.Window)
			}

			err = excelio.OutputFileRecord(rt, savepath, w.File, []model.FileMoveRecord{}, nil)
			if err != nil {
				dialog.ShowError(fmt.Errorf("生成数据库遇到错误：%s", err), w.Window)
			} else {
				dialog.ShowInformation("完成", fmt.Sprintf("你的数据已被保存在 %s ..", savepath), w.Window)
			}
		}, w.Window)

		dlg.SetFileName("all_record_data.xlsx")
		dlg.Show()
	})

	outputAllWthSearch := fyne.NewMenuItem("导出全部迁移数据（过滤后）", func() {
		rt.Action()
		dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer == nil {
				return
			} else if err != nil {
				dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w.Window)
				return
			}

			defer func() {
				_ = writer.Close()
			}()

			savepath := writer.URI().Path()

			if !strings.HasSuffix(savepath, ".xlsx") {
				dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w.Window)
			}

			err = excelio.OutputFileRecord(rt, savepath, w.File, []model.FileMoveRecord{}, &w.SearchRecord)
			if err != nil {
				dialog.ShowError(fmt.Errorf("生成数据库遇到错误：%s", err), w.Window)
			} else {
				dialog.ShowInformation("完成", fmt.Sprintf("你的数据已被保存在 %s ..", savepath), w.Window)
			}
		}, w.Window)

		dlg.SetFileName("all_record_data_with_condition.xlsx")
		dlg.Show()
	})

	outputAllEveryOne := fyne.NewMenuItem("导出所有档案全部迁移数据", func() {
		rt.Action()
		dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer == nil {
				return
			} else if err != nil {
				dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w.Window)
				return
			}

			defer func() {
				_ = writer.Close()
			}()

			savepath := writer.URI().Path()

			if !strings.HasSuffix(savepath, ".xlsx") {
				dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w.Window)
			}

			err = excelio.OutputFileRecord(rt, savepath, nil, []model.FileMoveRecord{}, nil)
			if err != nil {
				dialog.ShowError(fmt.Errorf("生成数据库遇到错误：%s", err), w.Window)
			} else {
				dialog.ShowInformation("完成", fmt.Sprintf("你的数据已被保存在 %s ..", savepath), w.Window)
			}
		}, w.Window)

		dlg.SetFileName("all_record_data.xlsx")
		dlg.Show()
	})

	outputAllWthSearchEveryOne := fyne.NewMenuItem("导出所有档案全部迁移数据（过滤后）", func() {
		rt.Action()
		dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer == nil {
				return
			} else if err != nil {
				dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w.Window)
				return
			}

			defer func() {
				_ = writer.Close()
			}()

			savepath := writer.URI().Path()

			if !strings.HasSuffix(savepath, ".xlsx") {
				dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w.Window)
			}

			err = excelio.OutputFileRecord(rt, savepath, nil, []model.FileMoveRecord{}, &w.SearchRecord)
			if err != nil {
				dialog.ShowError(fmt.Errorf("生成数据库遇到错误：%s", err), w.Window)
			} else {
				dialog.ShowInformation("完成", fmt.Sprintf("你的数据已被保存在 %s ..", savepath), w.Window)
			}
		}, w.Window)

		dlg.SetFileName("all_record_data_with_condition.xlsx")
		dlg.Show()
	})

	outputNow := fyne.NewMenuItem("导出本页数据", func() {
		rt.Action()
		dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer == nil {
				return
			} else if err != nil {
				dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w.Window)
				return
			}

			defer func() {
				_ = writer.Close()
			}()

			savepath := writer.URI().Path()

			if !strings.HasSuffix(savepath, ".xlsx") {
				dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w.Window)
			}

			err = excelio.OutputFileRecord(rt, savepath, w.File, w.InfoRecord, nil)
			if err != nil {
				dialog.ShowError(fmt.Errorf("生成数据库遇到错误：%s", err), w.Window)
			} else {
				dialog.ShowInformation("完成", fmt.Sprintf("你的数据已被保存在 %s ..", savepath), w.Window)
			}
		}, w.Window)

		dlg.SetFileName("one_page_of_record_data.xlsx")
		dlg.Show()
	})

	quit := fyne.NewMenuItem("关闭", func() {
		WinClose(w.Window)
		w.Window = nil
	})
	quit.IsQuit = true

	dangan := fyne.NewMenu("档案", show, quit)
	page := fyne.NewMenu("分页")
	sousuo := fyne.NewMenu("搜索", search)
	fenxi := fyne.NewMenu("统计", outputAll, outputAllWthSearch, outputNow, outputAllEveryOne, outputAllWthSearchEveryOne)

	menu := fyne.NewMainMenu(dangan, sousuo, page, fenxi)

	w.Menu = &Controller{
		Record: w,
		Menu:   menu,
		Page:   page,
		Search: ss,
	}

	w.Window.SetMainMenu(w.Menu.Menu)

	return w.Menu
}

func (c *Controller) ChangePageMenuItemRecord(rt runtime.RunTime, pageItemCount int, p int64, pageMax int64, message string) {
	pageList := make([]*fyne.MenuItem, 0, pageMax)
	if pageMax <= 0 {
		pageList = append(pageList, fyne.NewMenuItem("暂无数据", func() {
			rt.Action()
		}))
		c.Record.NowPage = 0
	} else {
		if p > pageMax {
			p = pageMax
		}
		c.Record.NowPage = p
		for k := int64(1); k <= pageMax; k++ {
			i := k
			if i == p {
				m := fyne.NewMenuItem(fmt.Sprintf("第%d页（当前页）", i), func() {
					rt.Action()
					dialog.ShowConfirm("是否需要重载？", message, func(b bool) {
						rt.Action()
						if !b {
							return
						}
						c.Record.Table.UpdateTableRecord(rt, pageItemCount, i)
					}, c.Record.Window)
				})
				pageList = append(pageList, m)
			} else {
				pageList = append(pageList, fyne.NewMenuItem(fmt.Sprintf("第%d页", i), func() {
					rt.Action()
					c.Record.Table.UpdateTableRecord(rt, pageItemCount, i)
				}))
			}
		}
	}

	c.Page.Items = pageList
	c.Page.Refresh()
	c.Menu.Refresh()
}
