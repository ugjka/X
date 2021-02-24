// Suspend your machine when bluetooth headset disconnects
// Assumes 1 pulseaudio card without bluetooth
package main

import (
	"bytes"
	"os/exec"
	"time"
)

func main() {
	var out []byte
	var lines int
	for {
		out, _ = exec.Command("pactl", "list", "cards", "short").Output()
		lines = bytes.Count(out, []byte{'\n'})
		if lines < 2 {
			exec.Command("systemctl", "suspend").Run()
			return
		}
		time.Sleep(time.Second * 5)
	}
}
