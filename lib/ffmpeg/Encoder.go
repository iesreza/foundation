package ffmpeg

import (
	"github.com/iesreza/foundation/lib/ref"
	"regexp"
)

type ApplyVideo struct {
	Aspect     string `validate:"\\d+\\:\\d+" switch:"-aspect"`
	Bitrate    string `validate:"\\d+[K,M,k,m]*" switch:"-b:v"`
	CropBottom string `validate:"\\d+" switch:"-cropbottom"`
	CropTop    string `validate:"\\d+" switch:"-croptop"`
	CropLeft   string `validate:"\\d+" switch:"-cropleft"`
	CropRight  string `validate:"\\d+" switch:"-cropright"`
	Codec      string `validate:".+" switch:"-b:a"`
	MaxBitrate string `validate:"\\d+[K,M,k,m]*" switch:"-maxrate"`
	MinBitrate string `validate:"\\d+[K,M,k,m]*" switch:"-minrate"`
	BufferSize string `validate:"\\d+[K,M,k,m]*" switch:"-bufsize"`
	Size       string `validate:"\\d+\\x\\d+" switch:"-s"`
	StripAudio bool   `switch:"-an"`
}

type ApplyAudio struct {
	Bitrate    string `validate:"\\d+[K,M,k,m]*" switch:"-b:a"`
	Codec      string `validate:".+" switch:"-b:a"`
	SampleRate string `validate:"\\d+" switch:"-ar"`
	StripVideo bool   `switch:"-vn"`
}

func (v *ApplyVideo) buildCommand() []string {
	pipe := []string{}
	obj := ref.Parse(v)

	for _, field := range obj.Fields {
		//reflect.Indirect(v).FieldByName(field.Name).
		v, _ := obj.Get(field.Name)
		if v.Type().String() == "string" {
			if v.String() != "" && field.Tag.Get("switch") != "" {
				if field.Tag.Get("validate") != "" {
					matched, err := regexp.MatchString(field.Tag.Get("validate"), v.String())
					if err != nil {
						continue
					}
					if !matched {
						continue
					}
				}
				pipe = append(pipe, field.Tag.Get("switch"), v.String())
			}

		}

		if v.Type().String() == "bool" {
			if v.Bool() && field.Tag.Get("switch") != "" {
				pipe = append(pipe, field.Tag.Get("switch"))
			}

		}
	}

	return pipe
}

func (v *ApplyAudio) buildCommand() []string {
	pipe := []string{}
	obj := ref.Parse(v)

	for _, field := range obj.Fields {
		//reflect.Indirect(v).FieldByName(field.Name).
		v, _ := obj.Get(field.Name)
		if v.Type().String() == "string" {
			if v.String() != "" && field.Tag.Get("switch") != "" {
				if field.Tag.Get("validate") != "" {
					matched, err := regexp.MatchString(field.Tag.Get("validate"), v.String())
					if err != nil {
						continue
					}
					if !matched {
						continue
					}
				}
				pipe = append(pipe, field.Tag.Get("switch"), v.String())
			}

		}

		if v.Type().String() == "bool" {
			if v.Bool() && field.Tag.Get("switch") != "" {
				pipe = append(pipe, field.Tag.Get("switch"))
			}

		}
	}

	return pipe
}

func (f *FFMPEG) ApplyVideo(v ApplyVideo) {
	f.applyVideo = v
}

func (f *FFMPEG) ApplyAudio(v ApplyAudio) {
	f.applyAudio = v
}
