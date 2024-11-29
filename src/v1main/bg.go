package v1main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/SuperH-0630/hdangan/src/systeminit"
	"image/color"
)

type Bg struct {
	canvas.Text
	x, y float32
}

func NewBg(x, y float32) *Bg {
	t := "宋子桓"
	c, err := systeminit.GetInit()
	if err == nil {
		t = c.Yaml.Report.Name
	}

	text := canvas.NewText(t, color.NRGBA{R: 0, B: 0, G: 0, A: 0})
	return &Bg{
		Text: *text,
		x:    x,
		y:    y,
	}
}

func (b *Bg) MinSize() fyne.Size {
	b.Text.MinSize()
	return fyne.Size{Width: b.x, Height: b.y}
}

func (b *Bg) Size() fyne.Size {
	mSize := b.MinSize()
	size := b.Text.Size()
	return fyne.NewSize(fmax(mSize.Width, size.Width), fmax(mSize.Height, size.Height))
}

func fmax(a float32, b ...float32) float32 {
	for _, v := range b {
		if v > a {
			a = v
		}
	}
	return a
}
