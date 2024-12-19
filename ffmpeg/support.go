package ffmpeg

import "path/filepath"

var SupportedVideoExts = []string{".mp4", ".mkv", ".avi", ".mov", ".flv", ".wmv", ".webm", ".mpeg", ".3gp", ".m4v"}
var SupportedVideoCodecs = []string{"libx264", "libx265"}
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
