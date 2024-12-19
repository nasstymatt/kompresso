package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type ColoredLabel struct {
	widget.BaseWidget
	Text  string
	Color color.Color
}

func NewColoredLabel(text string, color color.Color) *ColoredLabel {
	label := &ColoredLabel{Text: text, Color: color}
	label.ExtendBaseWidget(label)
	return label
}

func (l *ColoredLabel) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(l.Text, l.Color)
	return &coloredLabelRenderer{label: l, text: text}
}

type coloredLabelRenderer struct {
	label *ColoredLabel
	text  *canvas.Text
}

func (r *coloredLabelRenderer) Layout(size fyne.Size) {
	r.text.Resize(size)
}

func (r *coloredLabelRenderer) MinSize() fyne.Size {
	return r.text.MinSize()
}

func (r *coloredLabelRenderer) Refresh() {
	r.text.Text = r.label.Text
	r.text.Color = r.label.Color
	canvas.Refresh(r.text)
}

func (r *coloredLabelRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *coloredLabelRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.text}
}

func (r *coloredLabelRenderer) Destroy() {}
