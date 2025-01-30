// mp3converter project main.go
package main

import (
	"fmt"
	"strings"
	ui "videoconverter/ui"
	video "videoconverter/video"

	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ffmpeg -i "/home/archeopternix/Videos/2002 Weihnachten Vasoldsberg_opt2a.avi" -b:a 128k -c:a aac -c:v libx264 -crf 21 -pix_fmt yuv420p -preset medium "/home/archeopternix/Videos/2002 Weihnachten Vasoldsberg_opt2a.mp4"
// ffmpeg -i "/home/archeopternix/Videos/2002 Weihnachten Vasoldsberg_opt1.avi" -b:a 128k -c:a aac -c:v prores -profile:v 3 -vf format=yuv420p "/home/archeopternix/Videos/NEU/2002 Weihnachten Vasoldsberg_opt1.mov"

func main() {
	myApp := app.NewWithID("Video Converter")
	myWindow := myApp.NewWindow("Drag and Drop Example")
	var videoSettings *video.VideoSettings

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

	setParameter := widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
		dialogContent, settings := ui.SetupParameterDialog()
		dialogContent.Resize(fyne.NewSize(500, 400))
		dialog.ShowCustomConfirm("Start Coding", "Save", "Cancel", dialogContent, func(ok bool) {
			if ok {
				videoSettings = settings
			}

		}, myWindow)
	})

	exit := widget.NewToolbarAction(theme.CancelIcon(), func() {
		myApp.Quit()
	})

	run := widget.NewToolbarAction(theme.DownloadIcon(), func() {
		if videoSettings == nil {
			dialog.ShowError(fmt.Errorf("Video Settings are not set"), myWindow)
			return
		}

		if len(ui.List.VideoFiles) < 1 {
			dialog.ShowError(fmt.Errorf("No Videos to convert"), myWindow)
			return
		}

		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri == nil {
				return // User cancelled
			}
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}

			fmt.Println(uri.Path())
			progressBar.Show()
			progressBar.Max = float64(len(ui.List.VideoFiles))
			progressBar.SetValue(0)

			succ := 0
			for i, file := range ui.List.VideoFiles {
				targetpath := filepath.Join(uri.Path(), fileNameWithoutExtension(file.Uri.Name())+"."+videoSettings.VideoContainer)
				if err := video.ConvertVideo(file.Uri.Path(), targetpath, *videoSettings); err != nil {
					fmt.Println(err)
					continue
				}
				progressBar.SetValue(float64(i))
				succ += 1
			}
			progressBar.Hide()
			dialog.ShowInformation("Conversion to "+videoSettings.VideoContainer, fmt.Sprintf("%d/%d files converted successfully", succ, len(ui.List.VideoFiles)), myWindow)
		}, myWindow).Show()
	})

	toolbar := widget.NewToolbar(
		openFile,
		run,
		setParameter,
		exit,
	)

	content := container.NewBorder(toolbar, progressBar, nil, nil, scrollList)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}
