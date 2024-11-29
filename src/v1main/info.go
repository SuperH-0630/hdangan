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

func readData(rt runtime.RunTime, f *model.File) (*fyne.Container, *fyne.Container) {
	fileComment := widget.NewMultiLineEntry()
	fileComment.Text = strToStr(f.FileComment, "")
	fileComment.Wrapping = fyne.TextWrapWord
	fileComment.OnChanged = func(s string) {
		rt.Action()
		fileComment.Text = strToStr(f.FileComment, "")
	}

	moveComment := widget.NewMultiLineEntry()
	moveComment.Text = strToStr(f.MoveComment, "")
	moveComment.Wrapping = fyne.TextWrapWord
	moveComment.OnChanged = func(s string) {
		rt.Action()
		moveComment.Text = strToStr(f.MoveComment, "")
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
		widget.NewLabel("最早迁入时间："),
		widget.NewLabel(fmt.Sprintf("%s", f.FirstMoveIn.Format("2006-01-02 15:04:05"))),

		widget.NewLabel("最后迁入时间："),
		widget.NewLabel(fmt.Sprintf("%s", f.LastMoveIn.Format("2006-01-02 15:04:05"))),

		widget.NewLabel("最后迁入迁出时间："),
		widget.NewLabel(fmt.Sprintf("%s", f.MoveStatus)),

		widget.NewLabel("最后迁出人："),
		widget.NewLabel(fmt.Sprintf("%s", strToStr(f.MoveOutPeopleName, "暂无"))),

		widget.NewLabel("最后迁出单位："),
		widget.NewLabel(fmt.Sprintf("%s", strToStr(f.MoveOutPeopleUnit, "暂无"))),

		widget.NewLabel("最后迁出备注："),
		moveComment,
	)

	return left, right
}

func ShowInfo(rt runtime.RunTime, f *model.File, refresh func(rt runtime.RunTime)) {
	infoWindow := rt.App().NewWindow(fmt.Sprintf("详细信息-%s-%d", f.Name, f.FileID))

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

	move := widget.NewButton("迁入迁出", func() {
		rt.Action()
		ShowMove(rt, f, warpRefresh)
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
				err := model.DeleteFile(f)
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
