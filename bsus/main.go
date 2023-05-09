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
			exec.Command("systemctl", "suspend", "-i").Run()
			return
		}
		time.Sleep(time.Second * 5)
	}
}
