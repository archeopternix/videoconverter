// video.go
package video

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"

	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/u2takey/ffmpeg-go"
)

type VideoSettings struct {
	VideoContainer string
	VideoFormat    string
	AudioFormat    string
	CRF            string
	Preset         string
	Profile        string
	AQuality       string
}

func ConvertVideo(inputFile, outputFile string, settings VideoSettings) error {
	var args ffmpeg_go.KwArgs
	switch settings.VideoFormat {
	case "prores":
		args = ffmpeg_go.KwArgs{
			"c:v":       settings.VideoFormat,
			"c:a":       settings.AudioFormat, // Use AAC codec for audio
			"b:a":       settings.AQuality,    // Medium quality audio bitrate
			"profile:v": settings.Profile,     // Medium quality audio bitrate
			"vf":        "format=yuv420p",
		}
	case "libsvtav1":
		presmap := map[string]string{
			"ultrafast": "0",
			"superfast": "1",
			"veryfast":  "2",
			"faster":    "3",
			"fast":      "4",
			"medium":    "5",
			"slow":      "6",
			"slower":    "7",
			"veryslow":  "8",
		}

		preset, ok := presmap[settings.Preset]
		if !ok {
			preset = "5"
		}

		args = ffmpeg_go.KwArgs{
			"c:v":     settings.VideoFormat,
			"c:a":     settings.AudioFormat, // Use AAC codec for audio
			"b:a":     settings.AQuality,    // Medium quality audio bitrate
			"preset":  preset,               // Medium quality audio bitrate
			"crf":     settings.CRF,         // Medium quality audio bitrate
			"pix_fmt": "yuv420p",            // Standard pixel format
		}
	default:
		args = ffmpeg_go.KwArgs{
			"c:v":     settings.VideoFormat,
			"c:a":     settings.AudioFormat, // Use AAC codec for audio
			"b:a":     settings.AQuality,    // Medium quality audio bitrate
			"preset":  settings.Preset,      // Medium quality audio bitrate
			"crf":     settings.CRF,         // Medium quality audio bitrate
			"pix_fmt": "yuv420p",            // Standard pixel format
		}
	}
	// ffmpeg -i "in.avi" -c:v libsvtav1 -crf 28 -preset 6 -vf scale=720:576 -r 25 -strict experimental  -c:a aac -b:a 192k "out.mp4"

	buf := bytes.NewBuffer(nil)

	err := ffmpeg_go.Input(inputFile).
		Output(outputFile,
			args,
		).
		WithOutput(buf, os.Stdout).
		Run()

	return err
}

// checkIfVideo checks if the given file is a video based on its MIME type
func CheckIfVideo(filePath string) (bool, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to read file: %w", err)
	}
	mime := mimetype.Detect(data)

	if strings.Contains(mime.String(), "video") {
		return true, nil
	}

	return false, nil
}

// extractFrame uses ffmpeg-go to extract a frame from a video at a specified time
func ExtractFrame(videoPath, frameTime string) (image.Image, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg_go.Input(videoPath, ffmpeg_go.KwArgs{"ss": frameTime}).
		Output("pipe:", ffmpeg_go.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, nil).
		Run()
	if err != nil {
		return nil, fmt.Errorf("failed to extract frame: %w", err)
	}

	img, err := jpeg.Decode(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return img, nil
}
