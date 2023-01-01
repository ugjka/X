// Fit a video into a 8mb file (Discord nitro pls?)
// Needs ffmpeg ffprobe
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// in kilobits per second
const AUDIO_BITRATE = 96

func main() {
	size := flag.Float64("size", 8, "target size in MB")
	preset := flag.String("preset", "slow", "h264 encode preset")
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "error: no filename given")
		os.Exit(1)
	}
	file := flag.Args()[0]
	probe := exec.Command(
		"ffprobe", "-i", file, "-show_entries",
		"format=duration", "-v", "quiet", "-of",
		"csv=p=0")
	secbytes, err := probe.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	secbytes = bytes.TrimSpace(secbytes)
	seconds, err := strconv.ParseFloat(string(secbytes), 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	bitfloat := *size * 1024.0 * 8.0 / seconds
	// video bitrate
	bitrate := int(bitfloat) - AUDIO_BITRATE

	pass1 := exec.Command("ffmpeg", "-y", "-i", file, "-c:v", "libx264", "-preset", *preset,
		"-b:v", fmt.Sprintf("%dk", bitrate), "-pass", "1", "-passlogfile", "fflogfile", "-c:a", "aac", "-b:a", fmt.Sprintf("%dk", AUDIO_BITRATE), "-f", "mp4", "/dev/null")
	pass2 := exec.Command("ffmpeg", "-y", "-i", file, "-c:v", "libx264", "-preset", *preset,
		"-b:v", fmt.Sprintf("%dk", bitrate), "-pass", "2", "-passlogfile", "fflogfile", "-c:a", "aac", "-b:a", fmt.Sprintf("%dk", AUDIO_BITRATE), "8mb."+file)

	pass1.Stderr = os.Stderr
	pass1.Stdout = os.Stdout
	pass2.Stderr = os.Stderr
	pass2.Stdout = os.Stdout

	cleanup := func() {
		os.Remove("fflogfile-0.log")
		os.Remove("fflogfile-0.log.mbtree")
	}

	err = pass1.Run()
	if err != nil {
		cleanup()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	err = pass2.Run()
	if err != nil {
		cleanup()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cleanup()
}
