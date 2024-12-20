package ffmpeg

import (
	"fmt"
	"path/filepath"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var SupportedVideoExts = []string{".mp4", ".mkv", ".avi", ".mov", ".flv", ".wmv", ".webm", ".mpeg", ".3gp", ".m4v"}
var SupportedVideoCodecs = []string{"libx264", "libx265", "h264_nvenc", "hevc_nvenc"}
var DefautlCodec = SupportedVideoCodecs[0]
var DefaultQuality = 50.0

// IsVideoExtSupported checks if the file extension of the given path is supported for video files.
// It iterates through the SupportedVideoExts map and compares the file extension with the supported extensions.
// If a match is found, it returns true, indicating that the video extension is supported.
// Otherwise, it returns false.
//
// Parameters:
//   - path: The file path to check the extension of.
//
// Returns:
//   - bool: True if the file extension is supported, false otherwise.
func IsVideoExtSupported(path string) bool {
	for _, ext := range SupportedVideoExts {
		if filepath.Ext(path) == ext {
			return true
		}
	}
	return false
}

func GetGlobalArgsForCodec(codec string) []string {
	baseArgs := []string{"-hide_banner", "-nostats", "-loglevel", "error"}

	switch codec {
	case "h264_nvenc":
	case "hevc_nvenc":
		baseArgs = append(baseArgs, "-hwaccel", "cuda", "-hwaccel_output_format", "cuda")
	}

	return baseArgs
}

func GetKWArgsForCodec(codec string, quality float64) ffmpeg.KwArgs {
	baseArgs := ffmpeg.KwArgs{
		"c:v":      codec,
		"pix_fmt":  "yuv420p",
		"vf":       "pad=ceil(iw/2)*2:ceil(ih/2)*2",
		"movflags": "+faststart",
	}

	switch codec {
	case "libx264", "libx265":
		baseArgs["preset"] = "slow"
		baseArgs["crf"] = calculateCRF(int(quality))
	case "h264_nvenc", "hevc_nvenc":
		baseArgs["preset"] = "p4"
		baseArgs["cq"] = calculateCQ(int(quality)) // Or use bitrate: baseArgs["b:v"] = calculateBitrate(quality)
	}

	return baseArgs
}

func calculateCRF(quality int) string {
	const (
		maxCRF     = 36
		minCRF     = 16
		defaultCRF = 28
	)

	if quality < 0 || quality > 100 {
		return fmt.Sprintf("%d", defaultCRF)
	}

	diff := (maxCRF - minCRF) * (100 - quality) / 100
	crf := minCRF + diff

	return fmt.Sprintf("%d", crf)
}

func calculateCQ(quality int) string {
	const (
		maxCQ     = 31
		minCQ     = 14
		defaultCQ = 23
	)

	if quality < 0 || quality > 100 {
		return fmt.Sprintf("%d", defaultCQ)
	}

	diff := (maxCQ - minCQ) * (100 - quality) / 100
	cq := minCQ + diff

	return fmt.Sprintf("%d", cq)
}

func calculateBitrate(quality float64) string {
	const (
		maxBitrate = 10000 // kbps
		minBitrate = 2000  // kbps
	)

	bitrate := minBitrate + (maxBitrate-minBitrate)*(quality/100)
	return fmt.Sprintf("%dk", int(bitrate))
}
