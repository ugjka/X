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
// *******************
// About:
// Get a transcript of a youtube video by ID
//
// Needs yt-dlp
// Tested only on Linux
//
// To build this you need the Go compiler:
// go build -o ytext main.go
//

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

const USAGE = `Usage: %s [OPTIONS] [Youtube ID]

Get a transcript of a youtube video by ID
(default language: 'en')

Options:
-lang string
	  language iso code, 
`

const YTCMD = "yt-dlp -4 --skip-download --write-auto-subs --sub-format json3 --sub-langs %s --"

func main() {
	// check for yt-dlp
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	lang := flag.String("lang", "en", "language (iso code)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, USAGE, path.Base(os.Args[0]))
	}
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "error: no youtube id given")
		os.Exit(1)
	}
	id := flag.Args()[0]

	curdir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cmdarr := strings.Split(fmt.Sprintf(YTCMD, *lang), " ")
	tmp, err := os.MkdirTemp("", "ytext")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	err = os.Chdir(tmp)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cmdarr = append(cmdarr, id)
	cmd := exec.Command(cmdarr[0], cmdarr[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.RemoveAll(tmp)
		os.Exit(1)
	}
	dir, err := os.ReadDir("./")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(dir) == 0 {
		fmt.Fprintln(os.Stderr, "error: transcript not found")
		os.RemoveAll(tmp)
		os.Exit(1)
	}
	name := dir[0].Name()
	jsonfile, err := os.Open(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.RemoveAll(tmp)
		os.Exit(1)
	}
	err = os.Chdir(curdir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	err = os.WriteFile(name[:len(name)-6]+".txt", json3toText(jsonfile), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	err = os.RemoveAll(tmp)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func json3toText(r io.Reader) []byte {
	var data JSON3
	var b = bytes.NewBuffer(nil)
	json.NewDecoder(r).Decode(&data)
	for _, v := range data.Events {
		for _, v := range v.Segs {
			b.WriteString(v.UTF8)
		}
	}
	b.WriteString("\n")
	return b.Bytes()
}

type JSON3 struct {
	Events []struct {
		Segs []struct {
			UTF8 string `json:"utf8"`
		}
	}
}
