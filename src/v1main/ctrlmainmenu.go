package v1main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/SuperH-0630/hdangan/src/aboutme"
	"github.com/SuperH-0630/hdangan/src/excelreader"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"os"
	"strings"
)

var Menu *fyne.MainMenu
var Page *fyne.Menu
var NowPage int64 = 0

func getMainMenu(rt runtime.RunTime, w fyne.Window, refresh func(rt runtime.RunTime)) *fyne.MainMenu {
	if Menu != nil {
		return Menu
	}

	lucky := fyne.NewMenuItem("启动/关闭彩蛋", func() {
		rt.Action()
		res := ChangeGame(rt)
		if res == TurnOn {
			dialog.ShowInformation("提示", "彩蛋已经触发。", w)
		} else {
			dialog.ShowInformation("提示", "彩蛋已经关闭。", w)
		}
	})

	search := fyne.NewMenuItem("搜索条件", func() {
		rt.Action()
		ShowWhereWindow(rt, &whereInfo, refresh)
	})

	initFile := fyne.NewMenuItem("导入配置文件", func() {
		rt.Action()
		NewInitFile(rt, w)
	})

	openFile := fyne.NewMenuItem("打开配置文件", func() {
		rt.Action()
		OpenInit(rt, w)
	})

	saveFile := fyne.NewMenuItem("另存配置文件", func() {
		rt.Action()
		SaveInit(rt, w)
	})

	copyFile := fyne.NewMenuItem("复制配置文件", func() {
		rt.Action()
		CopyInit(rt, w)
	})

	aboutMe := fyne.NewMenuItem("关于", func() {
		dialog.ShowInformation("关于开发者", aboutme.AboutMe, w)
		rt.Action()
	})

	quit := fyne.NewMenuItem("退出系统", func() {
		rt.App().Quit()
		os.Exit(0)
	})
	quit.IsQuit = true // declear quit menu

	newFile := fyne.NewMenuItem("新建档案", func() {
		rt.Action()
		ShowNew(rt, refresh)
	})

	exclFile := fyne.NewMenuItem("模板导入", func() {
		rt.Action()
		err := AddFromFile(rt, w, refresh)
		if err != nil {
			dialog.ShowError(fmt.Errorf("导入失败：%s", err.Error()), w)
		}
	})

	template := fyne.NewMenuItem("保存模板", func() {
		rt.Action()
		dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer == nil {
				return
			} else if err != nil {
				dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w)
				return
			}

			defer func() {
				_ = writer.Close()
			}()

			savepath := writer.URI().Path()

			if !strings.HasSuffix(savepath, ".xlsx") {
				dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w)
			}

			err = excelreader.CreateTemplate(rt, savepath)
			if err != nil {
				dialog.ShowError(fmt.Errorf("模板保存失败：%s", err), w)
			}
		}, w)

		dlg.SetFileName("template.xlsx")
		dlg.Show()
	})

	tj := fyne.NewMenuItem("数据统计", func() {
		rt.Action()
		TongJi(rt, w)
	})

	outputAll := fyne.NewMenuItem("导出全部档案", func() {
		rt.Action()
		dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer == nil {
				return
			} else if err != nil {
				dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w)
				return
			}

			defer func() {
				_ = writer.Close()
			}()

			savepath := writer.URI().Path()

			if !strings.HasSuffix(savepath, ".xlsx") {
				dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w)
			}

			err = excelreader.OutputFile(rt, savepath, []model.File{}, nil)
			if err != nil {
				dialog.ShowError(fmt.Errorf("生成数据库遇到错误：%s", err), w)
			} else {
				dialog.ShowInformation("完成", fmt.Sprintf("你的数据已被保存在 %s ..", savepath), w)
			}
		}, w)

		dlg.SetFileName("all_data.xlsx")
		dlg.Show()
	})

	outputAllWthSearch := fyne.NewMenuItem("导出全部档案（过滤后）", func() {
		rt.Action()
		dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer == nil {
				return
			} else if err != nil {
				dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w)
				return
			}

			defer func() {
				_ = writer.Close()
			}()

			savepath := writer.URI().Path()

			if !strings.HasSuffix(savepath, ".xlsx") {
				dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w)
			}

			err = excelreader.OutputFile(rt, savepath, []model.File{}, &whereInfo)
			if err != nil {
				dialog.ShowError(fmt.Errorf("生成数据库遇到错误：%s", err), w)
			} else {
				dialog.ShowInformation("完成", fmt.Sprintf("你的数据已被保存在 %s ..", savepath), w)
			}
		}, w)

		dlg.SetFileName("all_data_with_condition.xlsx")
		dlg.Show()
	})

	outputNow := fyne.NewMenuItem("保存当前表格", func() {
		rt.Action()
		dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer == nil {
				return
			} else if err != nil {
				dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w)
				return
			}

			defer func() {
				_ = writer.Close()
			}()

			savepath := writer.URI().Path()

			if !strings.HasSuffix(savepath, ".xlsx") {
				dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w)
			}

			err = excelreader.OutputFile(rt, savepath, InfoFile, nil)
			if err != nil {
				dialog.ShowError(fmt.Errorf("生成数据库遇到错误：%s", err), w)
			} else {
				dialog.ShowInformation("完成", fmt.Sprintf("你的数据已被保存在 %s ..", savepath), w)
			}
		}, w)

		dlg.SetFileName("one_page_of_data.xlsx")
		dlg.Show()
	})

	outputAllEveryOne := fyne.NewMenuItem("导出所有档案全部迁移数据", func() {
		rt.Action()
		dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer == nil {
				return
			} else if err != nil {
				dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w)
				return
			}

			defer func() {
				_ = writer.Close()
			}()

			savepath := writer.URI().Path()

			if !strings.HasSuffix(savepath, ".xlsx") {
				dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w)
			}

			err = excelreader.OutputFileRecord(rt, savepath, nil, nil, nil)
			if err != nil {
				dialog.ShowError(fmt.Errorf("生成数据库遇到错误：%s", err), w)
			} else {
				dialog.ShowInformation("完成", fmt.Sprintf("你的数据已被保存在 %s ..", savepath), w)
			}
		}, w)

		dlg.SetFileName("all_record_data.xlsx")
		dlg.Show()
	})

	outputAllWthSearchEveryOne := fyne.NewMenuItem("导出所有档案全部迁移数据（过滤后）", func() {
		rt.Action()
		ww := NewSaveWhereWindow(rt, func(rt runtime.RunTime, s *SaveWhereWindow) {
			dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
				if writer == nil {
					return
				} else if err != nil {
					dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w)
					return
				}

				defer func() {
					_ = writer.Close()
				}()

				savepath := writer.URI().Path()

				if !strings.HasSuffix(savepath, ".xlsx") {
					dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w)
				}

				err = excelreader.OutputFileRecord(rt, savepath, nil, nil, &s.SearchRecord)
				if err != nil {
					dialog.ShowError(fmt.Errorf("生成数据库遇到错误：%s", err), w)
				} else {
					dialog.ShowInformation("完成", fmt.Sprintf("你的数据已被保存在 %s ..", savepath), w)
				}
			}, w)

			dlg.SetFileName("all_record_data_with_condition.xlsx")
			dlg.Show()
		})
		ww.Show(rt)
	})

	xitong := fyne.NewMenu("系统", newFile, exclFile, template, aboutMe, quit)
	peizhi := fyne.NewMenu("配置", initFile, openFile, saveFile, copyFile)
	sousuo := fyne.NewMenu("搜索", search)
	Page = fyne.NewMenu("分页")
	tongji := fyne.NewMenu("统计", tj, outputAll, outputAllWthSearch, outputNow, outputAllEveryOne, outputAllWthSearchEveryOne)
	caidan := fyne.NewMenu("彩蛋", lucky)

	menu := fyne.NewMainMenu(xitong, peizhi, sousuo, Page, tongji, caidan)

	return menu
}

func ChangePageMenuItem(rt runtime.RunTime, window fyne.Window, table *widget.Table, pageItemCount int, p int64, pageMax int64, message string) {
	pageList := make([]*fyne.MenuItem, 0, pageMax)
	if pageMax <= 0 {
		pageList = append(pageList, fyne.NewMenuItem("暂无数据", func() {
			rt.Action()
		}))
		NowPage = 0
	} else {
		if p > pageMax {
			p = pageMax
		}
		NowPage = p
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
						UpdateTable(rt, window, table, pageItemCount, i)
					}, window)
				})
				pageList = append(pageList, m)
			} else {
				pageList = append(pageList, fyne.NewMenuItem(fmt.Sprintf("第%d页", i), func() {
					rt.Action()
					UpdateTable(rt, window, table, pageItemCount, i)
				}))
			}
		}
	}

	Page.Items = pageList
	Page.Refresh()
	Menu.Refresh()
}
