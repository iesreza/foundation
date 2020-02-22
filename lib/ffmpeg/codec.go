package ffmpeg

import (
	"strings"
)

type Codec struct {
	Flag          string
	Slug          string
	Name          string
	Type          CodecType
	FrameLevelMT  bool
	SliceLevelMT  bool
	Experimental  bool
	DrawHorizBand bool
	DirectRender  bool
}

func (d *Codec) parseFlag() {
	if len(d.Flag) != 6 {
		return
	}
	switch d.Flag[0] {
	case 'V':
		d.Type = CODEC_VIDEO
		break
	case 'A':
		d.Type = CODEC_AUDIO
		break
	case 'S':
		d.Type = CODEC_SUBTITLE
		break
	}
	if d.Flag[1] == 'F' {
		d.FrameLevelMT = true
	}
	if d.Flag[2] == 'S' {
		d.SliceLevelMT = true
	}
	if d.Flag[3] == 'X' {
		d.Experimental = true
	}
	if d.Flag[4] == 'B' {
		d.DrawHorizBand = true
	}
	if d.Flag[5] == 'D' {
		d.DirectRender = true
	}
}

func HasDecoder(name string) bool {
	name = strings.ToLower(name)
	for _, item := range GetDecoders() {
		if item.Slug == name {
			return true
		}
	}
	return false
}

func SearchDecoder(name string) []*Codec {
	name = strings.ToLower(name)
	res := []*Codec{}
	for _, item := range GetDecoders() {
		if strings.Contains(item.Slug, name) {
			res = append(res, item)
		}
	}
	return res
}

func GetDecoders() []*Codec {
	if len(decoders) != 0 {
		return decoders
	}
	b, err := runWithTimeout("-decoders", defaultTimeout)

	if err == nil {
		lines := strings.Split(string(b), "\n")

		for i := 10; i < len(lines); i++ {
			fields := strings.Fields(lines[i])
			if len(fields) == 3 {
				dec := &Codec{
					Flag: fields[0],
					Slug: fields[1],
					Name: fields[2],
				}
				dec.parseFlag()
				decoders = append(decoders, dec)
			}
		}

	}
	return decoders
}

func HasEncoder(name string) bool {
	name = strings.ToLower(name)
	for _, item := range GetEncoders() {
		if item.Slug == name {
			return true
		}
	}
	return false
}

func SearchEncoder(name string) []*Codec {
	name = strings.ToLower(name)
	res := []*Codec{}
	for _, item := range GetEncoders() {
		if strings.Contains(item.Slug, name) {
			res = append(res, item)
		}
	}
	return res
}

func GetEncoders() []*Codec {
	if len(encoders) != 0 {
		return encoders
	}
	b, err := runWithTimeout("-encoders", defaultTimeout)

	if err == nil {
		lines := strings.Split(string(b), "\n")

		for i := 10; i < len(lines); i++ {
			fields := strings.Fields(lines[i])
			if len(fields) == 3 {
				dec := &Codec{
					Flag: fields[0],
					Slug: fields[1],
					Name: fields[2],
				}
				dec.parseFlag()
				encoders = append(encoders, dec)
			}
		}

	}
	return encoders
}
