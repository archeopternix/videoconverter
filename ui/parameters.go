package ui

import (
	video "videoconverter/video"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func SetupParameterDialog() (*fyne.Container, *video.VideoSettings) {

	settings := &video.VideoSettings{}

	// Container select box
	selectContainer := widget.NewSelect([]string{"mp4", "avi", "mov"}, func(value string) {
		settings.VideoContainer = value
	})
	selectContainerLabel := widget.NewLabel("Filetype")

	// Create select box
	selectVFormat := widget.NewSelect([]string{"H.264/MPEG-4 AVC", "H.265/HEVC", "AV1", "ProRes", "AVI"}, func(value string) {
		switch value {
		case "H.264/MPEG-4 AVC":
			settings.Profile = ""
			settings.VideoFormat = "libx264"
			settings.CRF = "21"
		case "H.265/HEVC":
			settings.Profile = ""
			settings.VideoFormat = "libx265"
			settings.CRF = "25"
		case "AV1":
			settings.Profile = ""
			settings.VideoFormat = "libsvtav1"
			settings.CRF = "26"
		case "ProRes":
			settings.Profile = "3"
			settings.VideoFormat = "prores"
		default:
			settings.Profile = ""
			settings.VideoFormat = "avi"
		}
	})
	selectVFormatLabel := widget.NewLabel("Codec   ")

	// Create entry fields
	crf := widget.NewEntry()
	crf.SetPlaceHolder("CRF")
	crf.SetText("23")
	crfLabel := widget.NewLabel("CRF 0..50")

	// Create select box
	selectPreset := widget.NewSelect([]string{"ultrafast", "faster", "medium", "slow", "ultraslow", ""}, func(value string) {
		settings.Preset = value
	})
	selectPresetLabel := widget.NewLabel("Speed   ")

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
	selectAFormatLabel := widget.NewLabel("Audio Format")

	// Profiles select box
	profiles := widget.NewSelect([]string{"H.264/AAC/mp4", "HEVC/AAC/mp4", "AV1/AAC/mp4", "ProRes/AAC/mov"}, func(value string) {
		switch value {
		case "H.264/AAC/mp4":
			selectContainer.SetSelected("mp4")
			selectVFormat.SetSelected("H.264/MPEG-4 AVC")
			selectPreset.SetSelected("medium")
			selectAFormat.SetSelected("AAC")
			crf.SetText("21")
			settings.CRF = "21"
			settings.Profile = ""
		case "HEVC/AAC/mp4":
			selectContainer.SetSelected("mp4")
			selectVFormat.SetSelected("H.265/HEVC")
			selectPreset.SetSelected("medium")
			selectAFormat.SetSelected("AAC")
			crf.SetText("25")
			settings.CRF = "25"
			settings.Profile = ""
		case "AV1/AAC/mp4":
			selectContainer.SetSelected("mp4")
			selectVFormat.SetSelected("AV1")
			selectPreset.SetSelected("medium")
			selectAFormat.SetSelected("AAC")
			crf.SetText("28")
			settings.CRF = "28"
			settings.Profile = ""
		case "ProRes/AAC/mov":
			selectContainer.SetSelected("mov")
			selectVFormat.SetSelected("ProRes")
			selectPreset.SetSelected("                              ")
			selectAFormat.SetSelected("AAC")
			crf.SetText("")
			settings.CRF = ""
			settings.Profile = "3"
		}
	})
	profilesLabel := widget.NewLabel("Profiles")
	profiles.SetSelected("H.264/AAC/mp4")

	// Create a container for the dialog content

	dialogContent := container.NewGridWithColumns(2,
		profilesLabel,
		profiles,
		widget.NewSeparator(),
		widget.NewSeparator(),
		selectContainerLabel,
		selectContainer,
		selectVFormatLabel,
		selectVFormat,
		crfLabel,
		crf,
		selectPresetLabel,
		selectPreset,
		selectAFormatLabel,
		selectAFormat,
	)
	return dialogContent, settings
}
