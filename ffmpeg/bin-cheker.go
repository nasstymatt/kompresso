package ffmpeg

import (
	"fmt"
	"os/exec"
	"regexp"
)

const ExecutableName = "ffmpeg"
const ProbeExecutableName = "ffprobe"

func FindProbeExecutableInPath() (string, error) {
	path, err := exec.LookPath(ProbeExecutableName)

	return path, err
}

func FindExecutableInPath() (string, error) {
	path, err := exec.LookPath(ExecutableName)

	return path, err
}

func GetExecutableVersion(path string) (string, error) {
	cmd := exec.Command(path, "-version")

	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	versionLine := string(output)

	re := regexp.MustCompile(`ffmpeg version (\S+)`)
	match := re.FindStringSubmatch(versionLine)

	if len(match) < 2 {
		return "", fmt.Errorf("FFmpeg version not found in output")
	}

	return match[1], nil
}
