// Fit a video into a 8mb file (Discord nitro pls?)
// Needs ffmpeg ffprobe
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

// in kilobits per second
const AUDIO_BITRATE = 96

func main() {
	file := os.Args[1]
	probe := exec.Command(
		"ffprobe", "-i", file, "-show_entries",
		"format=duration", "-v", "quiet", "-of",
		"csv=p=0")
	secbytes, err := probe.Output()
	if err != nil {
		log.Fatal(err)
	}
	secbytes = bytes.TrimSpace(secbytes)
	sec, err := strconv.ParseFloat(string(secbytes), 64)
	if err != nil {
		log.Fatal(err)
	}
	bitfloat := 8.0 * 1024.0 * 8.0 / sec
	// video bitrate
	bitrate := int(bitfloat) - AUDIO_BITRATE

	pass1 := exec.Command("ffmpeg", "-y", "-i", file, "-c:v", "libx264", "-preset", "medium",
		"-b:v", fmt.Sprintf("%dk", bitrate), "-pass", "1", "-passlogfile", "fflogfile", "-c:a", "aac", "-b:a", fmt.Sprintf("%dk", AUDIO_BITRATE), "-f", "mp4", "/dev/null")
	pass2 := exec.Command("ffmpeg", "-y", "-i", file, "-c:v", "libx264", "-preset", "medium",
		"-b:v", fmt.Sprintf("%dk", bitrate), "-pass", "2", "-passlogfile", "fflogfile", "-c:a", "aac", "-b:a", fmt.Sprintf("%dk", AUDIO_BITRATE), "8mb."+file)

	pass1.Stderr = os.Stderr
	pass1.Stdout = os.Stdout
	pass2.Stderr = os.Stderr
	pass2.Stdout = os.Stdout

	defer os.Remove("fflogfile-0.log")
	defer os.Remove("fflogfile-0.log.mbtree")

	err = pass1.Run()
	if err != nil {
		log.Fatal(err)
	}
	err = pass2.Run()
	if err != nil {
		log.Fatal(err)
	}
}
