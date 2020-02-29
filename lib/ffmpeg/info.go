package ffmpeg

import (
	"fmt"
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

type DMuxer struct {
	Slug    string
	Name    string
	Muxer   bool
	Demuxer bool
}

func GetHwAccelerators() []string {
	if len(accelerators) != 0 {
		return accelerators
	}
	b, err := runWithTimeout("-hwaccels", defaultTimeout)
	if err == nil {
		lines := strings.Split(string(b), "\n")

		for i := 1; i < len(lines); i++ {
			if len(strings.TrimSpace(lines[i])) > 2 {
				accelerators = append(accelerators, strings.TrimSpace(lines[i]))
			}
		}
	}
	return accelerators
}

func HasHwAccelerator(accl string) bool {
	for _, item := range accelerators {
		if item == accl {
			return true
		}
	}
	return false
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

func GetDecoder(name string) (*Codec, error) {
	name = strings.ToLower(name)
	for _, item := range GetDecoders() {
		if item.Slug == name {
			return item, nil
		}
	}
	return nil, fmt.Errorf("%s decoder not found", name)
}

func SearchDecoder(name string) []*Codec {
	name = strings.ToLower(name)
	res := []*Codec{}
	for _, item := range GetDecoders() {
		if strings.Contains(strings.ToLower(item.Slug+item.Name), name) {
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

			if len(fields) >= 3 {
				dec := &Codec{
					Flag: fields[0],
					Slug: fields[1],
					Name: strings.Join(fields[2:], " "),
				}
				dec.parseFlag()
				decoders = append(decoders, dec)
			}
		}

	}
	return decoders
}

func GetEncoder(name string) (*Codec, error) {
	name = strings.ToLower(name)
	for _, item := range GetEncoders() {
		if item.Slug == name {
			return item, nil
		}
	}
	return nil, fmt.Errorf("%s encoder not found", name)
}

func SearchEncoder(name string) []*Codec {
	name = strings.ToLower(name)
	res := []*Codec{}
	for _, item := range GetEncoders() {
		if strings.Contains(strings.ToLower(item.Slug+item.Name), name) {
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
			if len(fields) >= 3 {
				dec := &Codec{
					Flag: fields[0],
					Slug: fields[1],
					Name: strings.Join(fields[2:], ""),
				}
				dec.parseFlag()
				encoders = append(encoders, dec)
			}
		}

	}
	return encoders
}

func GetMuxers() []*DMuxer {
	if len(muxers) != 0 {
		return muxers
	}
	b, err := runWithTimeout("-muxers", defaultTimeout)

	if err == nil {
		lines := strings.Split(string(b), "\n")

		for i := 10; i < len(lines); i++ {
			fields := strings.Fields(lines[i])
			if len(fields) >= 3 {
				muxer := &DMuxer{
					Slug:    fields[1],
					Name:    strings.Join(fields[2:], ""),
					Muxer:   true,
					Demuxer: false,
				}

				muxers = append(muxers, muxer)
			}
		}

	}
	return muxers
}

func GetDeMuxers() []*DMuxer {
	if len(demuxers) != 0 {
		return demuxers
	}
	b, err := runWithTimeout("-demuxers", defaultTimeout)

	if err == nil {
		lines := strings.Split(string(b), "\n")

		for i := 10; i < len(lines); i++ {
			fields := strings.Fields(lines[i])
			if len(fields) >= 3 {
				demuxer := &DMuxer{
					Slug:    fields[1],
					Name:    strings.Join(fields[2:], ""),
					Muxer:   false,
					Demuxer: true,
				}

				demuxers = append(demuxers, demuxer)
			}
		}

	}
	return demuxers
}

func GetMuxer(name string) (*DMuxer, error) {
	name = strings.ToLower(name)
	for _, item := range GetMuxers() {
		if item.Slug == name {
			return item, nil
		}
	}
	return nil, fmt.Errorf("%s muxer not found", name)
}

func GetDeMuxer(name string) (*DMuxer, error) {
	name = strings.ToLower(name)
	for _, item := range GetDeMuxers() {
		if item.Slug == name {
			return item, nil
		}
	}
	return nil, fmt.Errorf("%s demuxer not found", name)
}

func SearchMuxer(name string) []*DMuxer {
	name = strings.ToLower(name)
	res := []*DMuxer{}
	for _, item := range GetMuxers() {
		if strings.Contains(strings.ToLower(item.Slug+item.Name), name) {
			res = append(res, item)
		}
	}
	return res
}

func SearchDeMuxer(name string) []*DMuxer {
	name = strings.ToLower(name)
	res := []*DMuxer{}
	for _, item := range GetDeMuxers() {
		if strings.Contains(strings.ToLower(item.Slug+item.Name), name) {
			res = append(res, item)
		}
	}
	return res
}
