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

func readRecordData(rt runtime.RunTime, f *model.File, r *model.FileMoveRecord) (*fyne.Container, *fyne.Container) {
	fileComment := widget.NewMultiLineEntry()
	fileComment.Text = strToStr(f.FileComment, "")
	fileComment.Wrapping = fyne.TextWrapWord
	fileComment.OnChanged = func(s string) {
		rt.Action()
		fileComment.Text = strToStr(f.FileComment, "")
	}

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
		widget.NewLabel(fmt.Sprintf("%d", f.FileID)),

		widget.NewLabel("姓名："),
		widget.NewLabel(fmt.Sprintf("%s", f.Name)),

		widget.NewLabel("身份证号："),
		widget.NewLabel(fmt.Sprintf("%s", f.IDCard)),

		widget.NewLabel("户籍地："),
		widget.NewLabel(fmt.Sprintf("%s", f.Location)),

		widget.NewLabel("卷宗标题："),
		widget.NewLabel(fmt.Sprintf("%s", f.FileTitle)),

		widget.NewLabel("卷宗类型："),
		widget.NewLabel(fmt.Sprintf("%s", f.FileType)),

		widget.NewLabel("卷宗备注："),
		fileComment,
	)

	rightLayout := layout.NewFormLayout()
	right := container.New(rightLayout,
		widget.NewLabel("迁入迁出状态："),
		widget.NewLabel(fmt.Sprintf("%s", r.MoveStatus)),

		widget.NewLabel("动作发生时间："),
		widget.NewLabel(fmt.Sprintf("%s", r.MoveTime.Format("2006-01-02 15:04:05"))),

		widget.NewLabel("最后迁出人："),
		widget.NewLabel(fmt.Sprintf("%s", strToStr(r.MoveOutPeopleName, "暂无"))),

		widget.NewLabel("最后迁出单位："),
		widget.NewLabel(fmt.Sprintf("%s", strToStr(r.MoveOutPeopleUnit, "暂无"))),

		widget.NewLabel("最后迁出备注："),
		moveComment,
	)

	return left, right
}

func ShowRecordInfo(rt runtime.RunTime, recordTableWindow fyne.Window, fileWindow fyne.Window, f *model.File, r *model.FileMoveRecord, refresh func(rt runtime.RunTime)) {
	infoWindow := rt.App().NewWindow(fmt.Sprintf("迁入迁出详细信息-%s-%d", f.Name, f.FileID))

	infoWindow.SetOnClosed(func() {
		rt.Action()
		WinClose(infoWindow)
		infoWindow = nil
	})
	infoWindow.SetCloseIntercept(func() {
		rt.Action()
		WinClose(infoWindow)
		infoWindow = nil
	})

	left, right := readRecordData(rt, f, r)

	upBox := container.NewHBox(left, right)

	change := widget.NewButton("同档案的其他迁入迁出记录", func() {
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

	downBox := container.NewHBox(change, move, del)
	downCenterBox := container.NewCenter(downBox)

	gg := NewBg(5, 20)

	box := container.NewVBox(upBox, gg, downCenterBox)
	cbox := container.NewCenter(box)

	bg := NewBg(max(cbox.MinSize().Width, cbox.Size().Width, 500),
		max(cbox.MinSize().Height, cbox.Size().Height, 300))

	lastContainer := container.NewStack(bg, cbox)
	infoWindow.SetContent(lastContainer)

	infoWindow.Show()
	infoWindow.CenterOnScreen()
}
