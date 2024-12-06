package v1main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/SuperH-0630/hdangan/src/aboutme"
	"github.com/SuperH-0630/hdangan/src/excelio"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"os"
	"strings"
)

type MainMenu struct {
	Window  *CtrlWindow
	Main    *fyne.MainMenu
	Page    *fyne.Menu
	FileSet *fyne.Menu
	NowPage int64
}

func getMainMenu(rt runtime.RunTime, w *CtrlWindow, refresh func(rt runtime.RunTime)) *MainMenu {
	lucky := fyne.NewMenuItem("启动/关闭彩蛋", func() {
		rt.Action()
		res := ChangeGame(rt)
		if res == TurnOn {
			dialog.ShowInformation("提示", "彩蛋已经触发。", w.window)
		} else {
			dialog.ShowInformation("提示", "彩蛋已经关闭。", w.window)
		}
	})

	wm := &MainMenu{
		Window:  w,
		NowPage: 1,
	}

	search := fyne.NewMenuItem("搜索条件", func() {
		rt.Action()
		ShowWhereWindow(rt, &wm.Window.table.whereInfo, refresh)
	})

	initFile := fyne.NewMenuItem("导入配置文件", func() {
		rt.Action()
		NewInitFile(rt, w.window)
	})

	openFile := fyne.NewMenuItem("打开配置文件", func() {
		rt.Action()
		OpenInit(rt, w.window)
	})

	saveFile := fyne.NewMenuItem("另存配置文件", func() {
		rt.Action()
		SaveInit(rt, w.window)
	})

	copyFile := fyne.NewMenuItem("复制配置文件", func() {
		rt.Action()
		CopyInit(rt, w.window)
	})

	aboutMe := fyne.NewMenuItem("关于", func() {
		dialog.ShowInformation("关于开发者", aboutme.AboutMe, w.window)
		rt.Action()
	})

	quit := fyne.NewMenuItem("退出系统", func() {
		rt.App().Quit()
		os.Exit(0)
	})
	quit.IsQuit = true // declear quit menu

	newFile := fyne.NewMenuItem("新建档案", func() {
		rt.Action()
		ShowNew(rt, w, refresh)
	})

	exclFile := fyne.NewMenuItem("模板导入", func() {
		rt.Action()
		err := AddFromFile(rt, w, refresh)
		if err != nil {
			dialog.ShowError(fmt.Errorf("导入失败：%s", err.Error()), w.window)
		}
	})

	exclRecordFile := fyne.NewMenuItem("模板导入（借出记录）", func() {
		rt.Action()
		err := AddRecordFromFile(rt, w, refresh)
		if err != nil {
			dialog.ShowError(fmt.Errorf("导入失败：%s", err.Error()), w.window)
		}
	})

	template := fyne.NewMenuItem("保存模板", func() {
		rt.Action()
		dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer == nil {
				return
			} else if err != nil {
				dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w.window)
				return
			}

			defer func() {
				_ = writer.Close()
			}()

			if !strings.HasSuffix(writer.URI().Path(), ".xlsx") {
				dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w.window)
			}

			err = excelio.CreateTemplate(rt, writer)
			if err != nil {
				dialog.ShowError(fmt.Errorf("模板保存失败：%s", err), w.window)
			}
		}, w.window)

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
				dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w.window)
				return
			}

			defer func() {
				_ = writer.Close()
			}()

			savepath := writer.URI().Path()

			if !strings.HasSuffix(savepath, ".xlsx") {
				dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w.window)
			}

			err = excelio.OutputFile(rt, w.table.fileSetType, savepath, nil, nil)
			if err != nil {
				dialog.ShowError(fmt.Errorf("生成数据库遇到错误：%s", err), w.window)
			} else {
				dialog.ShowInformation("完成", fmt.Sprintf("你的数据已被保存在 %s ..", savepath), w.window)
			}
		}, w.window)

		dlg.SetFileName("all_data.xlsx")
		dlg.Show()
	})

	outputAllWthSearch := fyne.NewMenuItem("导出全部档案（过滤后）", func() {
		rt.Action()
		dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer == nil {
				return
			} else if err != nil {
				dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w.window)
				return
			}

			defer func() {
				_ = writer.Close()
			}()

			savepath := writer.URI().Path()

			if !strings.HasSuffix(savepath, ".xlsx") {
				dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w.window)
			}

			err = excelio.OutputFile(rt, w.table.fileSetType, savepath, nil, &wm.Window.table.whereInfo)
			if err != nil {
				dialog.ShowError(fmt.Errorf("生成数据库遇到错误：%s", err), w.window)
			} else {
				dialog.ShowInformation("完成", fmt.Sprintf("你的数据已被保存在 %s ..", savepath), w.window)
			}
		}, w.window)

		dlg.SetFileName("all_data_with_condition.xlsx")
		dlg.Show()
	})

	outputNow := fyne.NewMenuItem("保存当前表格", func() {
		rt.Action()
		dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer == nil {
				return
			} else if err != nil {
				dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w.window)
				return
			}

			defer func() {
				_ = writer.Close()
			}()

			savepath := writer.URI().Path()

			if !strings.HasSuffix(savepath, ".xlsx") {
				dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w.window)
			}

			err = excelio.OutputFile(rt, w.table.fileSetType, savepath, wm.Window.table.InfoFile, nil)
			if err != nil {
				dialog.ShowError(fmt.Errorf("生成数据库遇到错误：%s", err), w.window)
			} else {
				dialog.ShowInformation("完成", fmt.Sprintf("你的数据已被保存在 %s ..", savepath), w.window)
			}
		}, w.window)

		dlg.SetFileName("one_page_of_data.xlsx")
		dlg.Show()
	})

	outputAllEveryOne := fyne.NewMenuItem("导出所有档案全部迁移数据", func() {
		rt.Action()
		dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer == nil {
				return
			} else if err != nil {
				dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w.window)
				return
			}

			defer func() {
				_ = writer.Close()
			}()

			savepath := writer.URI().Path()

			if !strings.HasSuffix(savepath, ".xlsx") {
				dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w.window)
			}

			err = excelio.OutputFileRecord(rt, savepath, nil, nil, nil)
			if err != nil {
				dialog.ShowError(fmt.Errorf("生成数据库遇到错误：%s", err), w.window)
			} else {
				dialog.ShowInformation("完成", fmt.Sprintf("你的数据已被保存在 %s ..", savepath), w.window)
			}
		}, w.window)

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
					dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w.window)
					return
				}

				defer func() {
					_ = writer.Close()
				}()

				savepath := writer.URI().Path()

				if !strings.HasSuffix(savepath, ".xlsx") {
					dialog.ShowError(fmt.Errorf("文件名必须以.xlsx结尾"), w.window)
				}

				err = excelio.OutputFileRecord(rt, savepath, nil, nil, &s.SearchRecord)
				if err != nil {
					dialog.ShowError(fmt.Errorf("生成数据库遇到错误：%s", err), w.window)
				} else {
					dialog.ShowInformation("完成", fmt.Sprintf("你的数据已被保存在 %s ..", savepath), w.window)
				}
			}, w.window)

			dlg.SetFileName("all_record_data_with_condition.xlsx")
			dlg.Show()
		})
		ww.Show(rt)
	})

	xitong := fyne.NewMenu("系统", newFile, exclFile, exclRecordFile, template, aboutMe, quit)
	wm.FileSet = fyne.NewMenu("档案类型")
	peizhi := fyne.NewMenu("配置", initFile, openFile, saveFile, copyFile)
	sousuo := fyne.NewMenu("搜索", search)
	wm.Page = fyne.NewMenu("分页")
	tongji := fyne.NewMenu("统计", tj, outputAll, outputAllWthSearch, outputNow, outputAllEveryOne, outputAllWthSearchEveryOne)
	caidan := fyne.NewMenu("彩蛋", lucky)

	wm.Main = fyne.NewMainMenu(xitong, peizhi, sousuo, wm.Page, tongji, caidan)
	return wm
}

func (m *MainMenu) ChangeFileSetModelItem(rt runtime.RunTime, pageItemCount int, message string) {
	setFileList := make([]*fyne.MenuItem, 0, len(model.FileSetTypeList))

	for _, t := range model.FileSetTypeList {
		name, ok := model.FileSetTypeName[t]
		if !ok {
			continue
		}

		if t == m.Window.table.fileSetType {
			name += "（当前）"
			setFileList = append(setFileList, fyne.NewMenuItem(name, func() {
				dialog.ShowConfirm("是否需要重载？", message, func(b bool) {
					rt.Action()
					if !b {
						return
					}
					m.Window.table.UpdateTable(rt, t, pageItemCount, m.NowPage)
				}, m.Window.window)
			}))
		} else {
			setFileList = append(setFileList, fyne.NewMenuItem(name, func() {
				rt.Action()
				m.Window.table.UpdateTable(rt, t, pageItemCount, 1)
			}))
		}
	}

	m.FileSet.Items = setFileList
	m.FileSet.Refresh()
	m.Main.Refresh()
}

func (m *MainMenu) ChangePageMenuItem(rt runtime.RunTime, pageItemCount int, p int64, pageMax int64, message string) {
	pageList := make([]*fyne.MenuItem, 0, pageMax)
	if pageMax <= 0 {
		pageList = append(pageList, fyne.NewMenuItem("暂无数据", func() {
			rt.Action()
		}))
		m.NowPage = 0
	} else {
		if p > pageMax {
			p = pageMax
		}
		m.NowPage = p
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
						m.Window.table.UpdateTable(rt, m.Window.table.fileSetType, pageItemCount, i)
					}, m.Window.window)
				})
				pageList = append(pageList, m)
			} else {
				pageList = append(pageList, fyne.NewMenuItem(fmt.Sprintf("第%d页", i), func() {
					rt.Action()
					m.Window.table.UpdateTable(rt, m.Window.table.fileSetType, pageItemCount, i)
				}))
			}
		}
	}

	m.Page.Items = pageList
	m.Page.Refresh()
	m.Main.Refresh()
}
