// Fit a video into a 8mb file (Discord nitro pls?)
// Needs ffmpeg ffprobe
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

// in kilobits per second
const AUDIO_BITRATE = 48

const DOWNSCALE = "-vf scale=iw/%f:ih/%f"

func main() {
	size := flag.Float64("size", 8, "target size in MB")
	preset := flag.String("preset", "slow", "h264 encode preset")
	downscale := flag.Float64("downscale", 0, "downscale multiplier")
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

	// construct output filename
	arr := strings.Split(file, ".")
	output := strings.Join(arr[0:len(arr)-1], ".")
	output = "8mb." + output + ".mp4"

	pass1 := exec.Command("ffmpeg", "-y", "-i", file, "-c:v", "libx264", "-preset", *preset,
		"-b:v", fmt.Sprintf("%dk", bitrate), "-pass", "1", "-passlogfile", file,
		"-an", "-f", "null", "/dev/null")
	pass2 := exec.Command("ffmpeg", "-y", "-i", file, "-c:v", "libx264", "-preset", *preset,
		"-b:v", fmt.Sprintf("%dk", bitrate), "-pass", "2", "-passlogfile", file, "-ac", "1",
		"-c:a", "libopus", "-b:a", fmt.Sprintf("%dk", AUDIO_BITRATE), output)

	if *downscale > 0 {
		pass1.Args = slices.Insert(pass1.Args, 4,
			strings.Split(fmt.Sprintf(DOWNSCALE, *downscale, *downscale), " ")...)
		pass2.Args = slices.Insert(pass2.Args, 4,
			strings.Split(fmt.Sprintf(DOWNSCALE, *downscale, *downscale), " ")...)
	}

	pass1.Stderr = os.Stderr
	pass1.Stdout = os.Stdout
	pass2.Stderr = os.Stderr
	pass2.Stdout = os.Stdout

	cleanup := func() {
		os.Remove(file + "-0.log")
		os.Remove(file + "-0.log.mbtree")
		os.Remove(file + "-0.log.temp")
		os.Remove(file + "-0.log.mbtree.temp")
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		if pass1.Process != nil {
			pass1.Process.Kill()
		}
		if pass2.Process != nil {
			pass2.Process.Kill()
		}
	}()

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
