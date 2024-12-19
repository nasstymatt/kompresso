package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
)

func NewVideoList(files binding.StringList) fyne.CanvasObject {
	videoList := container.NewVBox()

	updateList := func() {
		videoList.RemoveAll()
		paths, _ := files.Get()
		for _, path := range paths {
			videoList.Add(
				NewVideoEntry(path, func(path string) {
					files.Remove(path)
				}),
			)
		}
		videoList.Refresh()
	}

	files.AddListener(binding.NewDataListener(updateList))

	return videoList
}
