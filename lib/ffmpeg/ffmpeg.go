package ffmpeg

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

var (
	// ErrBinNotFound is returned when the ffprobe binary was not found
	ErrBinNotFound = fmt.Errorf("ffprobe bin not found")
	// ErrTimeout is returned when the ffprobe process did not succeed within the given time
	ErrTimeout = fmt.Errorf("process timeout exceeded")

	binPath        = "ffmpeg"
	decoders       = []*Codec{}
	encoders       = []*Codec{}
	defaultTimeout = 10 * time.Second
)

// SetFFPMPEGBinPath sets the global path to find and execute the ffmpeg program
func SetFFMPEGBinPath(newBinPath string) {
	binPath = newBinPath
}

type CodecType int

const (
	CODEC_VIDEO = iota
	CODEC_AUDIO
	CODEC_SUBTITLE
)

type FFMPEG struct {
}

type Progress struct {
	FramesProcessed string
	CurrentTime     string
	CurrentBitrate  string
	Progress        float64
	Speed           string
}

func New() *FFMPEG {
	p := FFMPEG{}
	return &p
}

func runWithTimeout(command string, timeout time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return runContext(ctx, command)
}

func runContext(ctx context.Context, command string) ([]byte, error) {
	var cmd *exec.Cmd
	var err error
	command = binPath + " " + command

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}

	var outputBuf bytes.Buffer
	cmd.Stdout = &outputBuf

	err = cmd.Start()
	if err == exec.ErrNotFound {
		return []byte{}, ErrBinNotFound
	} else if err != nil {
		return []byte{}, err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		err = cmd.Process.Kill()
		if err == nil {
			return []byte{}, ErrTimeout
		}
		return []byte{}, err
	case err = <-done:
		if err != nil {
			return []byte{}, err
		}
	}

	if err != nil {
		return outputBuf.Bytes(), err
	}

	return outputBuf.Bytes(), nil
}
