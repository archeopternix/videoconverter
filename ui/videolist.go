package ui

import (
	"fmt"
	"slices"
	video "videoconverter/video"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const frameTime = "00:00:01.000" // Time position (HH:MM:SS.mmm) to extract the frame

func DeleteVideoFile(file *VideoFile) {
	List.Delete(file)
}

var List *VideoList

// VideoFile represents one file  in the filelist with an image and the name
type VideoFile struct {
	Uri fyne.URI
	*fyne.Container
}

// NewVideoFile creates a new VideoFile and checks if is a video file,
// renders a preview image
func NewVideoFile(uri fyne.URI) (*VideoFile, error) {
	cnt := &VideoFile{
		Uri: uri,
	}

	// Check if the file is a video
	if isVideo, err := video.CheckIfVideo(uri.Path()); err != nil {
		return nil, fmt.Errorf("Error checking file type: %v", err)
	} else if !isVideo {
		return nil, fmt.Errorf("The file '%s' is not a video.", uri.Path())
	}

	// Extract frame from video
	frame, err := video.ExtractFrame(uri.Path(), frameTime)
	if err != nil {
		return nil, fmt.Errorf("Error extracting frame: %v", err)
	}

	// Create a fyne.Image
	image := canvas.NewImageFromImage(frame)
	image.SetMinSize(fyne.Size{Width: 90, Height: 75})
	image.FillMode = canvas.ImageFillContain

	// Create a label with the file nam3
	label := widget.NewLabel(uri.Name())

	// Create an button to delete the entry
	button := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		DeleteVideoFile(cnt)
	})

	cnt.Container = container.NewHBox(button, image, label)
	return cnt, nil
}

func (v VideoFile) IsEqual(vf *VideoFile) bool {
	if v.Uri.Path() == vf.Uri.Path() {
		return true
	}
	return false
}

type VideoList struct {
	*fyne.Container
	VideoFiles []VideoFile
}

func NewVideoList() *VideoList {
	return &VideoList{
		Container: container.NewVBox(),
	}
}

func (v *VideoList) Delete(file *VideoFile) {
	for i, f := range v.VideoFiles {
		if file.IsEqual(&f) {
			v.VideoFiles = slices.Delete(v.VideoFiles, i, i+1)
			break
		}
	}
	v.Remove(file.Container)
	v.Refresh()
}

func (v *VideoList) Add(file *VideoFile) error {
	for _, f := range v.VideoFiles {
		if file.IsEqual(&f) {
			return fmt.Errorf("File with path %s already exists", f.Uri.Name())
		}
	}
	v.VideoFiles = append(v.VideoFiles, *file)

	v.Container.Add(file.Container)
	return nil
}
