// Suspend your machine when bluetooth headset disconnects
package main

import (
	"bytes"
	"os/exec"
	"sync"
	"time"
)

func main() {
	var out []byte
	var lines int
	var once sync.Once
	var total int
	for {
		out, _ = exec.Command("pactl", "list", "cards", "short").Output()
		lines = bytes.Count(out, []byte{'\n'})
		once.Do(func() {
			total = lines
		})
		if lines < total {
			exec.Command("systemctl", "suspend").Run()
			return
		}
		time.Sleep(time.Second * 5)
	}
}
