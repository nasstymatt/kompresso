package widgets

import (
	"kompresso/ffmpeg"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type OnSubmitFunc func(codec string, quality float64)

func NewOutputSettings(onSubmit OnSubmitFunc) *widget.Form {
	selectedCodec := binding.NewString()
	quality := binding.NewFloat()

	codecSelect := widget.NewSelect(ffmpeg.SupportedVideoCodecs, func(codec string) {
		selectedCodec.Set(codec)
	})
	qualitySlider := widget.NewSliderWithData(0, 100, quality)

	codecSelect.SetSelected(ffmpeg.DefautlCodec)
	quality.Set(ffmpeg.DefaultQuality)

	return &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Codec", Widget: codecSelect},
			{Text: "Quality", Widget: widget.NewLabelWithData(binding.FloatToString(quality))},
			{Widget: qualitySlider},
		},
		SubmitText:  "Start Compression",
		Orientation: widget.Vertical,
		OnSubmit: func() {
			codec, _ := selectedCodec.Get()
			q, _ := quality.Get()

			onSubmit(codec, q)
		},
	}
}
