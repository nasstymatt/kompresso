package main

import (
	"kompresso/scenes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.New()
	w := app.NewWindow("Kompresso")

	w.Resize(fyne.NewSize(800, 600))
	w.SetFixedSize(true)

	scenes.ShowFfmpegFindScreen(w)

	w.ShowAndRun()
}
