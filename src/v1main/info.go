package v1main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
)

func readData(rt runtime.RunTime, fc model.File) (*fyne.Container, *fyne.Container) {
	f := fc.GetFile()

	comment := widget.NewMultiLineEntry()
	comment.Text = strToStr(f.Comment, "")
	comment.Wrapping = fyne.TextWrapWord
	comment.OnChanged = func(s string) {
		rt.Action()
		comment.Text = strToStr(f.Comment, "")
	}

	material := widget.NewMultiLineEntry()
	material.Text = strToStr(f.Material, "")
	material.Wrapping = fyne.TextWrapWord
	material.OnChanged = func(s string) {
		rt.Action()
		material.Text = strToStr(f.Material, "")
	}

	if !f.OldName.Valid {
		f.OldName.String = ""
	}

	if !f.IDCard.Valid {
		f.IDCard.String = ""
	}

	if !f.Comment.Valid {
		f.Comment.String = ""
	}

	sex := "男性"
	if !f.IsMan {
		sex = "女性"
	}

	leftLayout := layout.NewFormLayout()
	left := container.New(leftLayout,
		widget.NewLabel("卷宗号："),
		widget.NewLabel(fmt.Sprintf("%d", f.FileID)),

		widget.NewLabel("卷宗类型："),
		widget.NewLabel(fmt.Sprintf("%s", model.FileSetTypeName[f.FileSetType])),

		widget.NewLabel("文件联合编号："),
		widget.NewLabel(fmt.Sprintf("%d", f.FileUnionID)),

		widget.NewLabel("文件编号："),
		widget.NewLabel(fmt.Sprintf("%d", f.FileID)),

		widget.NewLabel("文件组内编号："),
		widget.NewLabel(fmt.Sprintf("%d", f.FileGroupID)),

		widget.NewLabel("姓名："),
		widget.NewLabel(fmt.Sprintf("%s", f.Name)),

		widget.NewLabel("曾用名："),
		widget.NewLabel(fmt.Sprintf("%s", f.OldName)),

		widget.NewLabel("身份证号："),
		widget.NewLabel(fmt.Sprintf("%s", f.IDCard)),

		widget.NewLabel("性别："),
		widget.NewLabel(fmt.Sprintf("%s", sex)),

		widget.NewLabel("出生日期："),
		widget.NewLabel(fmt.Sprintf("%s", f.Birthday.Format("2006-01-02"))),

		widget.NewLabel("卷宗备注："),
		comment,
	)

	rightDigit := make([]fyne.CanvasObject, 0, 10)

	switch ff := fc.(type) {
	case *model.FileQianRu:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			widget.NewLabel(fmt.Sprintf("%s", ff.Type)),

			widget.NewLabel("旧地址："),
			widget.NewLabel(fmt.Sprintf("%s", ff.OldLocation)),

			widget.NewLabel("新地址："),
			widget.NewLabel(fmt.Sprintf("%s", ff.NewLocation)),
		)
	case *model.FileChuSheng:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			widget.NewLabel(fmt.Sprintf("%s", ff.Type)),

			widget.NewLabel("地址："),
			widget.NewLabel(fmt.Sprintf("%s", ff.NewLocation)),
		)
	case *model.FileQianChu:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			widget.NewLabel(fmt.Sprintf("%s", ff.Type)),

			widget.NewLabel("新地址："),
			widget.NewLabel(fmt.Sprintf("%s", ff.NewLocation)),
		)
	case *model.FileSiWang:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			widget.NewLabel(fmt.Sprintf("%s", ff.Type)),

			widget.NewLabel("地址："),
			widget.NewLabel(fmt.Sprintf("%s", ff.Location)),
		)
	case *model.FileBianGeng:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			widget.NewLabel(fmt.Sprintf("%s", ff.Type)),

			widget.NewLabel("地址："),
			widget.NewLabel(fmt.Sprintf("%s", ff.Location)),
		)
	case *model.FileSuoNeiYiJu:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			widget.NewLabel(fmt.Sprintf("%s", ff.Type)),

			widget.NewLabel("地址："),
			widget.NewLabel(fmt.Sprintf("%s", ff.Location)),
		)
	case *model.FileSuoJianYiJu:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			widget.NewLabel(fmt.Sprintf("%s", ff.Type)),

			widget.NewLabel("地址："),
			widget.NewLabel(fmt.Sprintf("%s", ff.Location)),
		)
	case *model.FileNongZiZhuanFei:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			widget.NewLabel(fmt.Sprintf("%s", ff.Type)),

			widget.NewLabel("地址："),
			widget.NewLabel(fmt.Sprintf("%s", ff.Location)),
		)
	case *model.FileYiZhanShiQianYiZheng:
		rightDigit = append(rightDigit,
			widget.NewLabel("类型："),
			widget.NewLabel(fmt.Sprintf("%s", ff.Type)),

			widget.NewLabel("旧地址："),
			widget.NewLabel(fmt.Sprintf("%s", ff.Location)),
		)
	}

	rightDigit = append(rightDigit,
		widget.NewLabel("人数："),
		widget.NewLabel(fmt.Sprintf("%d", f.PeopleCount)),

		widget.NewLabel("办理时间："),
		widget.NewLabel(fmt.Sprintf("%s", f.Time.Format("2006-01-02 15:04:05"))),

		widget.NewLabel("页码："),
		widget.NewLabel(fmt.Sprintf("%d-%d（共%d页）", f.PageStart, f.PageEnd, f.PageCount)),

		widget.NewLabel("备考："),
		widget.NewLabel(fmt.Sprintf("%s", f.BeiKao.String)),

		widget.NewLabel("材料："),
		material,
	)

	rightLayout := layout.NewFormLayout()
	right := container.New(rightLayout, rightDigit...)

	return left, right
}

func ShowInfo(rt runtime.RunTime, f model.File, refresh func(rt runtime.RunTime)) {
	fc := f.GetFile()
	infoWindow := rt.App().NewWindow(fmt.Sprintf("详细信息-%s-%d", fc.Name, fc.FileID))

	infoWindow.SetOnClosed(func() {
		rt.Action()
		infoWindow = nil
	})
	infoWindow.SetCloseIntercept(func() {
		rt.Action()
		WinClose(infoWindow)
		infoWindow = nil
	})

	left, right := readData(rt, f)

	upBox := container.NewHBox(left, right)

	warpRefresh := func(rt runtime.RunTime) {
		left, right := readData(rt, f)
		upBox.RemoveAll()
		upBox.Add(left)
		upBox.Add(right)
		upBox.Refresh()
		refresh(rt)
	}

	change := widget.NewButton("更改信息", func() {
		rt.Action()
		ShowEdit(rt, f, warpRefresh)
	})

	move := widget.NewButton("出借档案", func() {
		rt.Action()
		ShowNewMove(rt, f, warpRefresh)
	})

	record := widget.NewButton("查看迁入迁出记录", func() {
		rt.Action()
		w := CreateRecordWindow(rt, f, infoWindow)
		w.Show()
		w.CenterOnScreen()
	})

	del := widget.NewButton("删除此条记录", func() {
		rt.Action()
		dialog.ShowConfirm("删除提示", "请问是否删除此条档案记录，删除后恢复可能较为苦难。", func(b bool) {
			if b {
				err := model.DeleteFile(rt, fc)
				if err != nil {
					dialog.ShowError(fmt.Errorf("数据库错误: %s", err.Error()), infoWindow)
				}
				refresh(rt)
				infoWindow.Close()
			}
		}, infoWindow)
	})

	downBox := container.NewHBox(change, move, record, del)
	downCenterBox := container.NewCenter(downBox)

	gg := NewBg(5, 20)

	box := container.NewVBox(upBox, gg, downCenterBox)
	cbox := container.NewCenter(box)

	bg := NewBg(fmax(cbox.MinSize().Width, cbox.Size().Width, 400),
		fmax(cbox.MinSize().Height, cbox.Size().Height, 350))

	lastContainer := container.NewStack(bg, cbox)
	infoWindow.SetContent(lastContainer)

	infoWindow.Show()
	infoWindow.CenterOnScreen()
	infoWindow.SetFixedSize(true)
}
