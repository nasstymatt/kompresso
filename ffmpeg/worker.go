package ffmpeg

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/natefinch/npipe"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ProcessVideo(path string, codec string, quality float64, updateProgress func(progress float64)) error {
	outputFile := filepath.Join(filepath.Dir(path), fmt.Sprintf("%s_compressed%s", filepath.Base(path), filepath.Ext(path)))

	probeData, err := ffmpeg.Probe(path)
	if err != nil {
		return fmt.Errorf("error probing file: %v", err)
	}
	totalDuration, err := probeDuration(probeData)
	if err != nil {
		return fmt.Errorf("error extracting duration: %v", err)
	}

	progressSocket, socketType := progressTempSock(totalDuration, updateProgress)

	var progressArg string
	if socketType == "pipe" {
		progressArg = progressSocket
	} else {
		progressArg = "unix://" + progressSocket
	}

	fmt.Println("quality", quality)

	kwArgs := GetKWArgsForCodec(codec, quality)
	globalArgs := GetGlobalArgsForCodec(codec)

	globalArgs = append(globalArgs, "-progress", progressArg)

	fmt.Printf("Compressing video: %f\n", quality)

	err = ffmpeg.Input(path).
		Output(outputFile, kwArgs).
		GlobalArgs(globalArgs...).
		OverWriteOutput().
		Run()

	if err != nil {
		return fmt.Errorf("error running ffmpeg: %v", err)
	}

	return nil
}

func progressTempSock(totalDuration float64, updateProgress func(progress float64)) (string, string) {
	rand.Seed(time.Now().Unix())
	var sockFileName string
	var listener net.Listener
	var err error
	var socketType string

	if runtime.GOOS == "windows" {
		sockFileName = fmt.Sprintf(`\\.\pipe\%d_sock`, rand.Int())
		listener, err = npipe.Listen(sockFileName)
		socketType = "pipe"
	} else {
		sockFileName = path.Join(os.TempDir(), fmt.Sprintf("%d_sock", rand.Int()))
		listener, err = net.Listen("unix", sockFileName)
		socketType = "unix"
	}

	if err != nil {
		log.Fatalf("error creating socket: %v", err)
	}

	go func() {
		defer listener.Close()
		re := regexp.MustCompile(`out_time_ms=(\d+)`)
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("error accepting connection: %v", err)
				return
			}
			defer conn.Close()

			buf := make([]byte, 16)
			data := ""
			for {
				_, err := conn.Read(buf)
				if err != nil {
					break
				}
				data += string(buf)
				matches := re.FindAllStringSubmatch(data, -1)
				if len(matches) > 0 {
					lastMatch := matches[len(matches)-1]
					outTimeMS, _ := strconv.Atoi(lastMatch[1])
					progress := float64(outTimeMS) / (totalDuration * 1_000_000)
					updateProgress(progress)
				}
				if strings.Contains(data, "progress=end") {
					updateProgress(1.0)
					return
				}
			}
		}
	}()

	return sockFileName, socketType
}

func probeDuration(a string) (float64, error) {
	type probeFormat struct {
		Duration string `json:"duration"`
	}
	type probeData struct {
		Format probeFormat `json:"format"`
	}
	pd := probeData{}
	err := json.Unmarshal([]byte(a), &pd)
	if err != nil {
		return 0, err
	}
	duration, err := strconv.ParseFloat(pd.Format.Duration, 64)
	if err != nil {
		return 0, err
	}
	return duration, nil
}
