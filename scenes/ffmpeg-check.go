package scenes

import (
	"fmt"
	"kompresso/ffmpeg"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	ffmpegExecutablePath  string
	ffprobeExecutablePath string
	ffmpegVersion         string
)

func ShowFfmpegFindScreen(w fyne.Window) {
	label := widget.NewLabel("Trying to find FFMPEG executable in PATH...")
	spinner := widget.NewActivity()

	spinner.Start()

	content := container.NewCenter(
		container.NewVBox(label, spinner),
	)

	w.SetContent(content)

	go func() {
		executablePath, err := ffmpeg.FindExecutableInPath()

		ffmpegExecutablePath = executablePath

		if err != nil {
			showErrorScreen(w, fmt.Sprintf("Error: %v", err))
			return
		}
		showFindProbeScreen(w)
	}()
}

func showFindProbeScreen(w fyne.Window) {
	label := widget.NewLabel("Trying to find FFPROBE executable in PATH...")
	spinner := widget.NewActivity()

	spinner.Start()

	content := container.NewCenter(
		container.NewVBox(label, spinner),
	)

	w.SetContent(content)

	go func() {
		path, err := ffmpeg.FindProbeExecutableInPath()

		ffprobeExecutablePath = path

		if err != nil {
			showErrorScreen(w, fmt.Sprintf("Error: %v", err))
			return
		}
		showFfmpegCheckingVersionScreen(w)
	}()
}

func showFfmpegCheckingVersionScreen(w fyne.Window) {
	label := widget.NewLabel(fmt.Sprintf("FFmpeg found!\nPath: %s\nChecking version...", ffmpegExecutablePath))
	spinner := widget.NewActivity()

	spinner.Start()

	content := container.NewCenter(
		container.NewVBox(label, spinner),
	)

	w.SetContent(content)

	go func() {
		version, err := ffmpeg.GetExecutableVersion(ffmpegExecutablePath)

		ffmpegVersion = version

		if err != nil {
			showErrorScreen(w, fmt.Sprintf("Error: %v", err))
			return
		}
		showSuccessScreen(w)
	}()
}

func showSuccessScreen(w fyne.Window) {
	label := widget.NewLabel(fmt.Sprintf("FFmpeg found!\nPath: %s\nVersion: %s", ffmpegExecutablePath, ffmpegVersion))
	spinner := widget.NewActivity()

	spinner.Start()

	content := container.NewCenter(
		container.NewVBox(
			label,
			widget.NewLabel("Launching Kompresso..."),
			spinner,
		),
	)

	w.SetContent(content)

	time.Sleep(time.Second)

	ShowFilePickerScene(w)
}

func showErrorScreen(w fyne.Window, message string) {
	label := widget.NewLabel(message)

	content := container.NewCenter(
		container.NewVBox(
			label,
			widget.NewButton("Close", func() {
				w.Close()
			}),
		),
	)

	w.SetContent(content)
}
