// MIT+NoAI License
//
// Copyright (c) 2023 ugjka <ugjka@proton.me>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights///
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// This code may not be used to train artificial intelligence computer models 
// or retrieved by artificial intelligence software or hardware.
//
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
