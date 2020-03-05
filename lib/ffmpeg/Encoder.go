package ffmpeg

import (
	"github.com/iesreza/foundation/lib/ref"
	"regexp"
	"strings"
)

type ApplyVideo struct {
	Aspect             string   `validate:"\\d+\\:\\d+" switch:"-aspect"`
	Bitrate            string   `validate:"\\d+[K,M,k,m]*" switch:"-b:v"`
	CropBottom         string   `validate:"\\d+" switch:"-cropbottom"`
	CropTop            string   `validate:"\\d+" switch:"-croptop"`
	CropLeft           string   `validate:"\\d+" switch:"-cropleft"`
	CropRight          string   `validate:"\\d+" switch:"-cropright"`
	Codec              string   `validate:".+" switch:"-c:v"`
	MaxBitrate         string   `validate:"\\d+[K,M,k,m]*" switch:"-maxrate"`
	MinBitrate         string   `validate:"\\d+[K,M,k,m]*" switch:"-minrate"`
	BufferSize         string   `validate:"\\d+[K,M,k,m]*" switch:"-bufsize"`
	Size               string   `validate:"\\d+\\x\\d+" switch:"-s"`
	StripAudio         bool     `switch:"-an"`
	Preset             string   `validate:".+" switch:"-preset"`
	ConstantRateFactor string   `validate:"\\d+" switch:"-crf"`
	Tune               string   `validate:".+" switch:"-tune"`
	Pass               string   `validate:"\\d+" switch:"-pass"`
	PixelFormat        string   `validate:".+" switch:"-pix_fmt"`
	StartTime          string   `validate:"\\d{2}\:\\d{2}:\\d{2}" switch:"-ss"`
	Duration           string   `validate:"\\d{2}\:\\d{2}:\\d{2}" switch:"-t"`
	Filter             []string `validate:".+" switch:"-vf"`
	ComplexFilter      []string `validate:".+" switch:"-filter_complex"`
	QuantizationMin    string   `validate:"\\d+" switch:"-qmin"`
	QuantizationMax    string   `validate:"\\d+" switch:"-qmax"`
	Subtitle           string   `validate:".+" switch:"-subtitles"`
}

type ApplyAudio struct {
	Bitrate    string   `validate:"\\d+[K,M,k,m]*" switch:"-b:a"`
	Codec      string   `validate:".+" switch:"-c:a"`
	SampleRate string   `validate:"\\d+" switch:"-ar"`
	StripVideo bool     `switch:"-vn"`
	Filter     []string `validate:".+" switch:"-filter:a"`
	Channels   string   `validate:"\\d+" switch:"-ac"`
}

func parsePipe(v interface{}) ([]string, error) {
	pipe := []string{}
	obj := ref.Parse(v)

	for _, field := range obj.Fields {
		//reflect.Indirect(v).FieldByName(field.Name).
		if field.Tag.Get("switch") == "" {
			continue
		}
		v, err := obj.Get(field.Name)
		if err != nil {
			return pipe, err
		}

		if field.Name == "Codec" && v.String() != "" {
			codec, err := GetDecoder(v.String())
			if err != nil {
				return pipe, err
			}
			pipe = append(pipe, field.Tag.Get("switch"), codec.Slug)
		}
		if v.Type().String() == "string" && v.String() != "" {

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

			continue
		}

		if v.Type().String() == "[]string" {
			slice := v.Interface().([]string)
			if len(slice) > 0 {
				pipe = append(pipe, field.Tag.Get("switch"), `"`+strings.Join(slice, " , ")+`"`)
			}
			continue
		}

		if v.Type().String() == "bool" && v.Bool() {

			pipe = append(pipe, field.Tag.Get("switch"))

		}
		continue
	}

	return pipe, nil
}

func (v *ApplyVideo) buildCommand() ([]string, error) {
	return parsePipe(v)
}

func (v *ApplyAudio) buildCommand() ([]string, error) {
	return parsePipe(v)
}

func (f *FFMPEG) ApplyVideo(v ApplyVideo) {
	f.applyVideo = v
}

func (f *FFMPEG) ApplyAudio(v ApplyAudio) {
	f.applyAudio = v
}
