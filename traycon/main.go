// MIT License | ugjka@proton.me
// Show tray icon while your script/program executes
package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/getlantern/systray"
)

func main() {
	// First Arg: Path to an image
	icon := os.Args[1]
	// Secord Arg: Path to a script, accepts arguments
	script := os.Args[2]
	// extra args
	extra := os.Args[3:]
	data, err := ioutil.ReadFile(icon)
	if err != nil {
		log.Fatal(err)
	}
	systray.SetIcon(data)
	runner := func() {
		cmd := exec.Command(script, extra...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		systray.Quit()
	}
	systray.Run(runner, nil)
}
