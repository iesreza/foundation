package ffmpeg

import (
	"fmt"
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
	accelerators   = []string{}
	muxers         = []*DMuxer{}
	demuxers       = []*DMuxer{}
	defaultTimeout = 10 * time.Second
	instance       = FFMPEG{}
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
	input         []Input
	output        Output
	applyVideo    ApplyVideo
	applyAudio    ApplyAudio
	hwaccel       string
	OverWrite     bool
	Loop          string `validate="\\d+" switch="-stream_loop"`
	LimitFileSize string `validate="\\d+[K,M,k,m,G,g]*" switch="-fs"`
	Threads       string `validate="\\d+" switch="-threads"`
	Metadata      string `validate=".+" switch="-metadata" separated="true"`
	VideoSync     bool   `switch="-vsync 1"`
	AudioSync     bool   `switch="-async 1"`
	Longest       bool   `switch="-longest"`
	Shortest      bool   `switch="-shortest"`
}

type Progress struct {
	FramesProcessed string
	CurrentTime     string
	CurrentBitrate  string
	Progress        float64
	Speed           string
}

type Tags struct {
	Encoder string `json:"ENCODER"`
}

func New() *FFMPEG {
	p := FFMPEG{}
	return &p
}

func (f *FFMPEG) HwAcceleration(v string) {
	f.hwaccel = v
}

func (f *FFMPEG) AutoHwAcceleration() {
	f.hwaccel = "auto"
}

func (f *FFMPEG) buildCommand() ([]string, error) {

	pipe := []string{}
	if len(f.input) == 0 {
		return pipe, fmt.Errorf("input is empty")
	}
	if f.hwaccel != "" {
		pipe = append(pipe, "-hwaccel", f.hwaccel)
	}
	for _, item := range f.input {
		p, err := item.buildCommand()
		if err != nil {
			return pipe, err
		}
		pipe = append(pipe, p...)
	}

	video, err := f.applyVideo.buildCommand()
	pipe = append(pipe, video...)
	if err != nil {
		return pipe, err
	}
	audio, err := f.applyAudio.buildCommand()
	pipe = append(pipe, audio...)
	if err != nil {
		return pipe, err
	}
	p, err := f.output.buildCommand()
	if err != nil {
		return pipe, err
	}
	pipe = append(pipe, p...)

	return pipe, nil
}
