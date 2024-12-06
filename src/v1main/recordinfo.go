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

func readRecordData(rt runtime.RunTime, fc model.File, r *model.FileMoveRecord) (*fyne.Container, *fyne.Container) {
	f := fc.GetFile()

	moveComment := widget.NewMultiLineEntry()
	moveComment.Text = strToStr(r.MoveComment, "")
	moveComment.Wrapping = fyne.TextWrapWord
	moveComment.OnChanged = func(s string) {
		rt.Action()
		moveComment.Text = strToStr(r.MoveComment, "")
	}

	leftLayout := layout.NewFormLayout()
	left := container.New(leftLayout,
		widget.NewLabel("卷宗号："),
		widget.NewLabel(fmt.Sprintf("%d", f.FileSetID)),

		widget.NewLabel("卷宗类型："),
		widget.NewLabel(fmt.Sprintf("%s", model.FileSetTypeName[f.FileSetType])),

		widget.NewLabel("联合文件编号：："),
		widget.NewLabel(fmt.Sprintf("%d", f.FileUnionID)),

		widget.NewLabel("出借记录号："),
		widget.NewLabel(fmt.Sprintf("%d", r.ID)),

		widget.NewLabel("档案类型："),
		widget.NewLabel(fmt.Sprintf("%s", model.FileSetTypeName[f.FileSetType])),

		widget.NewLabel("姓名："),
		widget.NewLabel(fmt.Sprintf("%s", f.Name)),

		widget.NewLabel("借出状态："),
		widget.NewLabel(fmt.Sprintf("%s", r.MoveStatus)),
	)

	rightLayout := layout.NewFormLayout()
	right := container.New(rightLayout,
		widget.NewLabel("借出时间："),
		widget.NewLabel(fmt.Sprintf("%s", r.MoveTime.Format("2006-01-02 15:04:05"))),

		widget.NewLabel("借入人："),
		widget.NewLabel(fmt.Sprintf("%s", strToStr(r.MoveOutPeopleName, "暂无"))),

		widget.NewLabel("借入单位："),
		widget.NewLabel(fmt.Sprintf("%s", strToStr(r.MoveOutPeopleUnit, "暂无"))),

		widget.NewLabel("借出人："),
		widget.NewLabel(fmt.Sprintf("%s", strToStr(r.MoveOutPeopleName, "暂无"))),

		widget.NewLabel("借出单位："),
		widget.NewLabel(fmt.Sprintf("%s", strToStr(r.MoveOutPeopleUnit, "暂无"))),

		widget.NewLabel("最后迁出备注："),
		moveComment,
	)

	return left, right
}

func ShowRecordInfo(rt runtime.RunTime, recordTableWindow fyne.Window, fileWindow fyne.Window, fc model.File, r *model.FileMoveRecord, refresh func(rt runtime.RunTime)) {
	f := fc.GetFile()

	infoWindow := rt.App().NewWindow(fmt.Sprintf("档案借出信息-%s-%d", f.Name, f.FileID))

	infoWindow.SetOnClosed(func() {
		rt.Action()
		infoWindow = nil
	})
	infoWindow.SetCloseIntercept(func() {
		rt.Action()
		WinClose(infoWindow)
		infoWindow = nil
	})

	left, right := readRecordData(rt, f, r)

	upBox := container.NewHBox(left, right)

	warpRefresh := func(rt runtime.RunTime) {
		left, right := readData(rt, f)
		upBox.RemoveAll()
		upBox.Add(left)
		upBox.Add(right)
		upBox.Refresh()
		refresh(rt)
	}

	change := widget.NewButton("编辑本条记录", func() {
		rt.Action()
		ShowMoveEdit(rt, fc, warpRefresh)
	})

	new_ := widget.NewButton("新建记录", func() {
		rt.Action()
		ShowNewMove(rt, fc, warpRefresh)
	})

	other := widget.NewButton("同档案的其他迁入迁出记录", func() {
		rt.Action()
		recordTableWindow.Show()
		recordTableWindow.CenterOnScreen()
	})

	move := widget.NewButton("查看档案", func() {
		rt.Action()
		fileWindow.Show()
		fileWindow.CenterOnScreen()
	})

	del := widget.NewButton("删除此条记录", func() {
		rt.Action()
		dialog.ShowInformation("你的操作很危险", "删除档案迁移记录是被视作危险的行为，将不被允许操作。", infoWindow)
	})

	downBox := container.NewHBox(change, new_, other, move, del)
	downCenterBox := container.NewCenter(downBox)

	gg := NewBg(5, 20)

	box := container.NewVBox(upBox, gg, downCenterBox)
	cbox := container.NewCenter(box)

	bg := NewBg(fmax(cbox.MinSize().Width, cbox.Size().Width, 500),
		fmax(cbox.MinSize().Height, cbox.Size().Height, 300))

	lastContainer := container.NewStack(bg, cbox)
	infoWindow.SetContent(lastContainer)

	infoWindow.Show()
	infoWindow.CenterOnScreen()
	infoWindow.SetFixedSize(true)
}
