package util

import (
	"background/common/logger"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
)

/*
	This method allows call ffmpeg ffprobe command and return the json output.
*/
func ffmpegProbeInfo(file string) ([]byte, error) {
	cmdName := fmt.Sprintf("ffprobe -v quiet -print_format json -show_format -show_streams %s ", file)
	logger.Debug(cmdName)
	cmdArgs := strings.Fields(cmdName)
	out, err := exec.Command(cmdArgs[0], cmdArgs[1:]...).Output()

	logger.Error(string(out))
	return out, err
}

/*
	This method allows to retrieve the corresponding file information from audio/video file.

	input:
		file 		:	input audio/video file path

*/
func FfmpegVideoInfo(file string) (name string, frameRate float32, videoCodec, audioCodec string, bitrate, videoBitrate, audioBitrate, duration, width, height uint32, size uint64, err error) {
	name, videoCodec, audioCodec = "", "", ""
	frameRate = 0
	bitrate, videoBitrate, audioBitrate, duration, width, height = 0, 0, 0, 0, 0, 0
	size = 0

	info, err := ffmpegProbeInfo(file)
	if err != nil {
		return
	}

	js, err := simplejson.NewJson(info)
	if err != nil {
		return
	}

	for i := 0; i < 3; i++ {
		stream := js.Get("streams").GetIndex(i)
		if stream == nil {
			continue
		}

		codecType, _ := stream.Get("codec_type").String()
		if codecType == "video" {
			w, _ := stream.Get("width").Int()
			h, _ := stream.Get("height").Int()
			width = uint32(w)
			height = uint32(h)

			videoCodec, _ = stream.Get("codec_name").String()
			avgFrameRate, _ := stream.Get("avg_frame_rate").String()
			frameRate, _ = calcRational(avgFrameRate)

			vBitrate, _ := stream.Get("bit_rate").String()
			vBitrateInt, _ := strconv.Atoi(vBitrate)
			videoBitrate = uint32(vBitrateInt)
		} else if codecType == "audio" {
			audioCodec, _ = stream.Get("codec_name").String()
			aBitrate, _ := stream.Get("bit_rate").String()
			aBitrateInt, _ := strconv.Atoi(aBitrate)
			audioBitrate = uint32(aBitrateInt)
		}
	}

	nStr, _ := js.Get("format").Get("filename").String()
	names := strings.Split(nStr, "/")
	name = names[len(names)-1]
	sizeStr, _ := js.Get("format").Get("size").String()
	size = uint64(Atoi64(sizeStr))
	durationStr, _ := js.Get("format").Get("duration").String()
	duration = uint32(Atof64(durationStr))
	bitrateStr, _ := js.Get("format").Get("bit_rate").String()
	bitrate = uint32(Atof64(bitrateStr))

	return
}

func calcRational(rational string) (float32, error) {
	sp := strings.Split(rational, "/")
	if len(sp) != 2 {
		err := errors.New("invalid rationl")
		logger.Error(err)
		return 0, err
	}

	numerator, err := strconv.ParseFloat(sp[0], 32)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	denominator, err := strconv.ParseFloat(sp[1], 32)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return float32(numerator / denominator), nil
}

func FfmpegMusicInfo(file string) (title, artist, album string, bitrate, duration, size uint64, err error) {
	title, artist, album = "", "", ""
	bitrate, duration, size = 0, 0, 0

	info, err := ffmpegProbeInfo(file)
	if err != nil {
		return
	}

	js, err := simplejson.NewJson(info)
	if err != nil {
		return
	}

	sizeStr, _ := js.Get("format").Get("size").String()
	durationStr, _ := js.Get("format").Get("duration").String()
	bitrateStr, _ := js.Get("format").Get("bit_rate").String()

	title, _ = js.Get("format").Get("tags").Get("title").String()
	artist, _ = js.Get("format").Get("tags").Get("artist").String()
	album, _ = js.Get("format").Get("tags").Get("album").String()
	bitrate = uint64(Atof64(bitrateStr))
	duration = uint64(Atof64(durationStr))
	size = uint64(Atoi64(sizeStr))

	return
}
