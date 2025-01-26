// mp3converter project main.go
package main

import (
	"fmt"
	ui "videoconverter/ui"
	video "videoconverter/video"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.NewWithID("Video Converter")
	myWindow := myApp.NewWindow("Drag and Drop Example")

	// Create a list to display images and file paths
	ui.List = ui.NewVideoList() //container.NewVBox()

	// Wrap the list in a scroll container
	scrollList := container.NewScroll(ui.List.Container)

	// Create a progress bar
	progressBar := widget.NewProgressBar()

	myWindow.SetOnDropped(func(pos fyne.Position, files []fyne.URI) {
		progressBar.Show()
		progressBar.Max = float64(len(files))
		progressBar.SetValue(0)

		for i, uri := range files {
			v, err := ui.NewVideoFile(uri)
			if err != nil {
				dialog.ShowError(err, myWindow)
				continue
			}

			err = ui.List.Add(v)
			if err != nil {
				dialog.ShowError(err, myWindow)
				continue
			}

			// update progressbar
			progressBar.SetValue(float64(i + 1))
		}
		ui.List.Refresh()
		progressBar.Hide()
	})

	openFile := widget.NewToolbarAction(theme.FileIcon(), func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()

			progressBar.Show()
			progressBar.Max = 1
			progressBar.SetValue(0)

			if v, err := ui.NewVideoFile(reader.URI()); err == nil {
				ui.List.Add(v)
			} else {
				dialog.ShowError(err, myWindow)
			}
			progressBar.SetValue(1)
			progressBar.Hide()
		}, myWindow).Show()
	})

	startCoding := widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
		settings := video.VideoSettings{}

		// Container select box
		selectContainer := widget.NewSelect([]string{"mp4", "avi", "mov"}, func(value string) {
			settings.VideoContainer = value
		})
		//	selectContainer.SetSelected("mp4")

		// Create select box
		selectVFormat := widget.NewSelect([]string{"H.264/MPEG-4 AVC", "H.265/HEVC", "AV1", "ProRes", "AVI"}, func(value string) {
			switch value {
			case "H.264/MPEG-4 AVC":
				settings.Profile = ""
				settings.VideoFormat = "libx264"
				settings.CRF = "23"
			case "H.265/HEVC":
				settings.Profile = ""
				settings.VideoFormat = "libx265"
				settings.CRF = "28"
			case "AV1":
				settings.Profile = ""
				settings.VideoFormat = "libaom-av1"
				settings.CRF = "30"
			case "ProRes":
				settings.Profile = "3"
				settings.VideoFormat = "prores"
			default:
				settings.Profile = ""
				settings.VideoFormat = "avi"
			}
		})

		// Create entry fields
		crf := widget.NewEntry()
		crf.SetPlaceHolder("CRF")
		crf.SetText("23")

		// Create select box
		selectPreset := widget.NewSelect([]string{"ultrafast", "fast", "medium", "slow", "ultraslow", ""}, func(value string) {
			settings.Preset = value
		})

		// Create select box
		selectAFormat := widget.NewSelect([]string{"AAC", "WAV", "MP3"}, func(value string) {
			switch value {
			case "AAC":
				settings.AudioFormat = "aac"
			case "WAV":
				settings.AudioFormat = "wav"
			case "MP3":
				settings.AudioFormat = "mp3"
			}
			settings.AQuality = "128k"
		})
		selectAFormat.SetSelected("AAC")

		// Container select box
		preSelect := widget.NewSelect([]string{"H.264/AAC/mp4", "HEVC/AAC/mp4", "AV1/AAC/mp4", "ProRes/AAC/mov"}, func(value string) {
			switch value {
			case "H.264/AAC/mp4":
				selectContainer.SetSelected("mp4")
				selectVFormat.SetSelected("H.264/MPEG-4 AVC")
				settings.CRF = "23"
				selectPreset.SetSelected("medium")
				selectAFormat.SetSelected("AAC")
				settings.Profile = ""

			case "HEVC/AAC/mp4":
				selectContainer.SetSelected("mp4")
				selectVFormat.SetSelected("H.265/HEVC")
				settings.CRF = "28"
				selectPreset.SetSelected("medium")
				selectAFormat.SetSelected("AAC")
				settings.Profile = ""
			case "AV1/AAC/mp4":
				selectContainer.SetSelected("mp4")
				selectVFormat.SetSelected("AV1")
				settings.CRF = "30"
				selectPreset.SetSelected("medium")
				selectAFormat.SetSelected("AAC")
				settings.Profile = ""
			case "ProRes/AAC/mov":
				selectContainer.SetSelected("mov")
				selectVFormat.SetSelected("ProRes")
				settings.CRF = ""
				selectPreset.SetSelected("")
				selectAFormat.SetSelected("AAC")
				settings.Profile = "3"
			}
		})
		preSelect.SetSelected("H.264/AAC/mp4")

		path := widget.NewEntry()
		path.SetPlaceHolder("target path for converted files")

		// Create a container for the dialog content
		dialogContent := container.NewVBox(
			preSelect,
			widget.NewSeparator(),
			selectContainer,
			selectVFormat,
			crf,
			selectPreset,
			selectAFormat,
			widget.NewSeparator(),
			path,
		)

		dialogContent.Resize(fyne.NewSize(500, 400))
		dialog.ShowCustomConfirm("Start Coding", "Run", "Close", dialogContent, func(run bool) {
			if run {
				video.ConvertVideo("D:", "D:", settings)
				fmt.Println("Run button clicked")
			} else {
				fmt.Println("Cancel button clicked")
			}
			// Add your run logic here
		}, myWindow)
	})

	exit := widget.NewToolbarAction(theme.CancelIcon(), func() {
		myApp.Quit()
	})

	toolbar := widget.NewToolbar(
		openFile,
		startCoding,
		exit,
	)

	content := container.NewBorder(toolbar, progressBar, nil, nil, scrollList)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}
