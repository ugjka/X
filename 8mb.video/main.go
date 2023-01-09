// MIT License | ugjka@proton.me
// https://8mb.video was down, so...
// Fit a video into a 8mb file (Discord nitro pls?)
// Needs ffmpeg ffprobe fdkaac
// Tested only on Linux
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
)

const USAGE = `Usage: %s [OPTIONS] [FILE]

Compress a video to target size
(default audio: mono 32kbps aac-he v1)

Options:
-down float
	  resolution downscale multiplier (default 1)
-music
	  64kbps stereo audio
-voice
	  16kbps mono audio
-preset string
	  h264 encode preset (default "slow")
-size float
	  target size in MB (default 8)
`

func main() {
	exes := []string{"ffmpeg", "ffprobe", "fdkaac"}
	for _, exe := range exes {
		if _, err := exec.LookPath(exe); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	size := flag.Float64("size", 8, "target size in MB")
	preset := flag.String("preset", "slow", "h264 encode preset")
	down := flag.Float64("down", 1, "resolution downscale multiplier")
	music := flag.Bool("music", false, "64kbps stereo audio")
	voice := flag.Bool("voice", false, "16kbps mono audio")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, USAGE, path.Base(os.Args[0]))
	}
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "error: no filename given")
		os.Exit(1)
	}
	if *down < 1 {
		fmt.Fprintln(os.Stderr, "downscale multiplier cannot be less than 1")
		os.Exit(1)
	}

	file := flag.Args()[0]

	// get video lenght in seconds
	probe := exec.Command(
		"ffprobe",
		"-i", file,
		"-show_entries", "format=duration",
		"-v", "quiet",
		"-of", "csv=p=0",
	)

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
	bitfloat := *size * 8388.608 / seconds

	// muxing overhead (not exact science)
	// based on observed values
	overhead := 96 / bitfloat * 0.07098343
	bitfloat -= bitfloat * overhead

	// deal with chunk overshoot on high bitrates
	switch {
	case bitfloat > 800:
		bitfloat -= 0.25 * 8388.608 / seconds
	case bitfloat > 400:
		bitfloat -= 0.0625 * 8388.608 / seconds
	default:
		bitfloat -= 0.03125 * 8388.608 / seconds
	}
	abitrate := 32
	audioch := 1
	if *music {
		abitrate *= 2
		audioch *= 2
	}
	if *voice {
		abitrate = 16
		audioch = 1
	}

	// video bitrate
	vbitrate := int(bitfloat) - abitrate

	// construct output filename
	arr := strings.Split(file, ".")
	output := strings.Join(arr[0:len(arr)-1], ".")
	output = fmt.Sprintf("%gmb.%s.mp4", *size, output)

	// resolution scale filter
	vfopt := fmt.Sprintf("scale=iw/%f:ih/%f", *down, *down)

	pass1 := exec.Command(
		"ffmpeg", "-y",
		"-i", file,
		"-vf", vfopt,
		"-c:v", "libx264",
		"-preset", *preset,
		"-b:v", fmt.Sprintf("%dk", vbitrate),
		"-pass", "1",
		"-passlogfile", file,
		"-movflags", "+faststart",
		"-an",
		"-f", "null",
		"/dev/null",
	)
	pass1.Stderr = os.Stderr
	pass1.Stdout = os.Stdout

	// we need to do this mumbo jumbo because fdk_aac encoder is disabled
	// on 99.99% of ffmpeg installations (even an Arch)
	// fdkaac standalone encoder is fine though
	wavfile := exec.Command(
		"ffmpeg", "-y",
		"-i", file,
		"-ar", "44100",
		"-ac", fmt.Sprintf("%d", audioch),
		file+".wav",
	)
	wavfile.Stderr = os.Stderr
	wavfile.Stdout = os.Stdout

	aacfile := exec.Command(
		"fdkaac",
		"-p", "5",
		"-b", fmt.Sprintf("%d000", abitrate),
		file+".wav",
	)
	aacfile.Stderr = os.Stderr
	aacfile.Stdout = os.Stdout

	pass2 := exec.Command(
		"ffmpeg", "-y",
		"-i", file,
		"-i", file+".m4a",
		"-vf", vfopt,
		"-c:v", "libx264",
		"-preset", *preset,
		"-b:v", fmt.Sprintf("%dk", vbitrate),
		"-pass", "2",
		"-passlogfile", file,
		"-movflags", "+faststart",
		"-c:a", "copy",
		"-map", "0:v:0",
		"-map", "1:a:0",
		output,
	)
	pass2.Stderr = os.Stderr
	pass2.Stdout = os.Stdout

	// remove tmp files
	cleanup := func() {
		os.Remove(file + "-0.log")
		os.Remove(file + "-0.log.mbtree")
		os.Remove(file + "-0.log.temp")
		os.Remove(file + "-0.log.mbtree.temp")
		os.Remove(file + ".wav")
		os.Remove(file + ".m4a")
	}

	// trap ctrl+c and kill
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		if pass1.Process != nil {
			pass1.Process.Kill()
		}
		if wavfile.Process != nil {
			wavfile.Process.Kill()
		}
		if aacfile.Process != nil {
			aacfile.Process.Kill()
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
	err = wavfile.Run()
	if err != nil {
		cleanup()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	err = aacfile.Run()
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
