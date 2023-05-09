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
// For Toshiba SATELLITE L850-1LK
// Toggle on and off active cooling
// Might need 'modprobe ec_sys write_support=1'
// Yeet
package main

import (
	"fmt"
	"os"
	"os/exec"
)

const ec = "/sys/kernel/debug/ec/ec0/io"

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	exec.Command("modprobe", "ec_sys", "write_support").Run()
	fd, err := os.OpenFile(ec, os.O_RDWR, 0600)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	switch os.Args[1] {
	case "on":
		fd.WriteAt([]byte{byte(255)}, 0xE3)
		fd.WriteAt([]byte{0x8}, 0xED)
	case "off":
		fd.WriteAt([]byte{byte(0)}, 0xE3)
	case "blast":
		fd.WriteAt([]byte{byte(255)}, 0xE3)
		fd.WriteAt([]byte{0x12}, 0xED)
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Options are 'on', 'off' or 'blast'")
	os.Exit(1)
}
