package scenes

import (
	"fmt"
	"kompresso/ffmpeg"
	"kompresso/widgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func ShowCompressorScene(w fyne.Window, filePaths []string) {
	files := binding.NewStringList()
	var outputSettings *widget.Form

	files.Set(filePaths)
	outputSettings = widgets.NewOutputSettings(func(codec string, quality float64) {
		outputSettings.SubmitText = "Compressing..."
		outputSettings.Disable()
		outputSettings.Refresh()

		go func() {
			for _, file := range filePaths {
				widgets.UpdateVideoEntryProgress(file, 0, "Waiting")
				err := ffmpeg.ProcessVideo(file, codec, quality, func(progress float64) {
					if progress == 1 {
						widgets.UpdateVideoEntryProgress(file, progress, "Completed")
					} else {
						widgets.UpdateVideoEntryProgress(file, progress, "Working...")
					}
				})
				if err != nil {
					widgets.UpdateVideoEntryProgress(file, 0, "Error")
					fmt.Printf("Error processing %s: %v\n", file, err)
				}
			}
			outputSettings.SubmitText = "Start Compression"
			outputSettings.Refresh()
			outputSettings.Enable()
		}()
	})
	videoList := widgets.NewVideoList(files)

	panels := container.NewHSplit(
		container.NewVScroll(videoList),
		container.NewVScroll(
			outputSettings,
		),
	)

	panels.Offset = 0.65

	w.SetContent(panels)
}
