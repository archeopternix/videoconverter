// mp3converter project main.go
package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const frameTime = "00:00:01.000" // Time position (HH:MM:SS.mmm) to extract the frame

func main() {
	myApp := app.NewWithID("Video Converter")
	myWindow := myApp.NewWindow("Drag and Drop Example")

	// Create a list to display images and file paths
	list := container.NewVBox()

	// Wrap the list in a scroll container
	scrollList := container.NewScroll(list)

	// Create a progress bar
	progressBar := widget.NewProgressBar()

	myWindow.SetOnDropped(func(pos fyne.Position, files []fyne.URI) {
		progressBar.Show()
		progressBar.Max = float64(len(files))
		progressBar.SetValue(0)

		for i, uri := range files {
			if err := addUriToList(list, uri); err != nil {
				fmt.Println(err)
			}
			// update progressbar
			progressBar.SetValue(float64(i + 1))
		}
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

			if err := addUriToList(list, reader.URI()); err != nil {
				fmt.Println(err)
			}

		}, myWindow).Show()
	})

	startCoding := widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
		settings := VideoSettings{}

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
				convertVideo("D:", "D:", settings)
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

func addUriToList(list *fyne.Container, uri fyne.URI) error {
	fmt.Println(uri.Path())

	// Check if the file is a video
	if isVideo, err := checkIfVideo(uri.Path()); err != nil {
		return fmt.Errorf("Error checking file type: %v", err)
	} else if !isVideo {
		return fmt.Errorf("The file is not a video.")
	}

	// Extract frame from video
	frame, err := extractFrame(uri.Path(), frameTime)
	if err != nil {
		return fmt.Errorf("Error extracting frame: %v", err)
	}

	var cnt *fyne.Container

	// Create an image
	image := canvas.NewImageFromImage(frame)
	image.SetMinSize(fyne.Size{Width: 64, Height: 48})
	image.FillMode = canvas.ImageFillContain

	// Create a label with the file path
	label := widget.NewLabel(uri.Path())

	// Create an button to delete the entry
	button := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		list.Remove(cnt)
	})

	cnt = container.NewHBox(button, image, label)

	list.Add(cnt)
	return nil
}
