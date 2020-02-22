package ffmpeg

import (
	"fmt"
	"testing"
)

func TestGetCodecs(t *testing.T) {
	fmt.Println(GetDecoders())
	fmt.Println(SearchDecoder("vorbis"))
	fmt.Println(HasDecoder("vorbis"))

	fmt.Println(GetEncoders())
	fmt.Println(SearchEncoder("VP9"))
	fmt.Println(HasEncoder("libvpx"))
}
