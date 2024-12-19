package widgets

import (
	"path/filepath"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	progressBars = make(map[string]*widget.ProgressBar)
	statusLabels = make(map[string]binding.String)
	mu           sync.Mutex
)

type OnCancelFunc func(path string)

func NewVideoEntry(path string, onCancel OnCancelFunc) fyne.CanvasObject {
	progress := binding.NewFloat()
	status := binding.NewString()

	status.Set("Waiting")

	progressBar := widget.NewProgressBarWithData(progress)
	statusLabel := NewColoredLabel("", theme.DisabledColor())

	status.AddListener(binding.NewDataListener(func() {
		value, _ := status.Get()

		switch value {
		case "Working...":
			statusLabel.Color = theme.PrimaryColor()
		case "Completed":
			statusLabel.Color = theme.SuccessColor()
		case "Error":
			statusLabel.Color = theme.ErrorColor()
		default:
			statusLabel.Color = theme.DisabledColor()
		}

		statusLabel.Text = value
		statusLabel.Refresh()
	}))

	mu.Lock()
	progressBars[path] = progressBar
	statusLabels[path] = status
	mu.Unlock()

	return container.NewVBox(
		container.NewHBox(
			widget.NewLabel(filepath.Base(path)),
			statusLabel,
		),
		container.NewBorder(
			nil,
			nil,
			nil,
			widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				onCancel(path)
			}),
			progressBar,
		),
	)
}

func UpdateVideoEntryProgress(path string, progress float64, status string) {
	mu.Lock()
	defer mu.Unlock()

	if bar, exists := progressBars[path]; exists {
		bar.SetValue(progress)
	}

	if label, exists := statusLabels[path]; exists {
		label.Set(status)
	}
}
