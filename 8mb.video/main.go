// https://8mb.video was down, so...
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
)

func main() {
	size := flag.Float64("size", 8, "target size in MB")
	preset := flag.String("preset", "slow", "h264 encode preset")
	down := flag.Float64("down", 1, "downscale multiplier")
	music := flag.Bool("music", false, "stereo audio")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "error: no filename given")
		os.Exit(1)
	}
	if *down < 0 {
		fmt.Fprintln(os.Stderr, "downscale multiplier cannot be negative")
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
	// -0.1mb because it sometimes overshoots
	bitfloat := (*size - 0.1) * 1024.0 * 8.0 / seconds
	//audio bitrate
	bitratea := 64
	audioch := 1
	if *music {
		bitratea *= 2
		audioch *= 2
	}
	// video bitrate
	bitratev := int(bitfloat) - bitratea

	// construct output filename
	arr := strings.Split(file, ".")
	output := strings.Join(arr[0:len(arr)-1], ".")
	output = "8mb." + output + ".mp4"
	vfopt := fmt.Sprintf("scale=iw/%f:ih/%f", *down, *down)

	pass1 := exec.Command("ffmpeg", "-y", "-i", file, "-vf", vfopt, "-c:v", "libx264", "-preset", *preset,
		"-b:v", fmt.Sprintf("%dk", bitratev), "-pass", "1", "-passlogfile", file,
		"-an", "-f", "null", "/dev/null")
	pass2 := exec.Command("ffmpeg", "-y", "-i", file, "-vf", vfopt, "-c:v", "libx264", "-preset", *preset,
		"-b:v", fmt.Sprintf("%dk", bitratev), "-pass", "2", "-passlogfile", file,
		"-ac", fmt.Sprintf("%d", audioch), "-c:a", "aac", "-b:a", fmt.Sprintf("%dk", bitratea), output)

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
