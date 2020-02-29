package ffmpeg

import "fmt"

type Input struct {
	Src    string
	Driver string
}

type Output struct {
	Path  string
	Muxer string
}

func (i Input) buildCommand() ([]string, error) {
	pipe := []string{}
	if i.Src == "" {
		return pipe, fmt.Errorf("input path is empty")
	}
	if i.Driver != "" {
		pipe = append(pipe, "-f", i.Driver)
	}
	pipe = append(pipe, "-i", i.Src)

	return pipe, nil
}

func (i Output) buildCommand() ([]string, error) {
	pipe := []string{}
	if i.Path == "" {
		return pipe, fmt.Errorf("ouput path is empty")
	}
	if i.Muxer == "" {
		return pipe, fmt.Errorf("ouput muxer is empty")
	}
	muxer, err := GetMuxer(i.Muxer)
	if err != nil {
		return pipe, err
	}
	pipe = append(pipe, "-f", muxer.Slug, i.Path)
	return pipe, nil
}

func (f *FFMPEG) Input(input ...Input) {
	f.input = append(f.input, input...)
}

func (f *FFMPEG) Output(output Output) {
	f.output = output
}
