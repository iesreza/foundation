package ffmpeg

import (
	"fmt"
	"testing"
)

func TestGetCodecs(t *testing.T) {

	/*	var err error

		// Create new instance of transcoder
		trans := new(transcoder.Transcoder)

		// Initialize transcoder passing the input file path and output file path
		err = trans.Initialize( "d:/300.mkv", "d:/330.mkv" )
		// Handle error...
		fmt.Println( trans.GetCommand() )
		// Start transcoder process with progress checking
		done := trans.Run(true)

		// Returns a channel to get the transcoding progress
		progress := trans.Output()

		// Example of printing transcoding progress
		for msg := range progress {
			fmt.Println(msg)
		}

		// This channel is used to wait for the transcoding process to end
		err = <-done
		fmt.Println(err)*/

	ffmpeg := New()
	ffmpeg.Input(Input{
		Src: "D:/300.mkv",
	})

	ffmpeg.Output(Output{
		Path:  "D:/300.2.mkv",
		Muxer: "mp4",
	})

	s := ApplyVideo{
		Codec:   "h264",
		CropTop: "300",
		Aspect:  "16:9",
		Filter:  []string{"a", "b", "c"},
	}

	ffmpeg.ApplyVideo(s)

	ffmpeg.HwAcceleration("cuvid")

	fmt.Println(ffmpeg.buildCommand())
}
