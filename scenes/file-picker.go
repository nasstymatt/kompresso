package scenes

import (
	"kompresso/ffmpeg"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ShowFilePickerScene(w fyne.Window) {
	label := canvas.NewText("Drop videos here", theme.Color(theme.ColorNameForeground))
	filesLabel := widget.NewLabel("")

	label.Alignment = fyne.TextAlignCenter
	label.TextSize = 24

	w.SetContent(
		container.NewCenter(
			container.NewVBox(
				label,
				filesLabel,
			),
		),
	)

	w.SetOnDropped(func(p fyne.Position, u []fyne.URI) {
		var acceptedFiles []string

		for _, uri := range u {
			if ffmpeg.IsVideoExtSupported(uri.Path()) {
				acceptedFiles = append(acceptedFiles, uri.Path())
			}
		}

		if len(acceptedFiles) == 0 {
			filesLabel.SetText("No video files found")
			return
		} else {
			ShowCompressorScene(w, acceptedFiles)
		}
	})
}
