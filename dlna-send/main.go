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
// ---------------
// send a url to your dlna receiver
// Usage:
//   dlna-send [options] URL
//     -dev string
//       dlna device to use (friendly name))
//     -stop
//       stop playback

package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/av1"
)

const USAGE = `Usage:
  %s [options] URL
    -dev string
      dlna device to use (friendly name))
    -stop
      stop playback
`

func main() {
	dev := flag.String("dev", "", "dlna device to use (friendly name))")
	stop := flag.Bool("stop", false, "stop playback")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, USAGE, path.Base(os.Args[0]))
	}
	flag.Parse()
	var url string
	var err error
	if !*stop {
		if len(flag.Args()) > 0 {
			url = flag.Args()[0]
		} else {
			stderr(fmt.Errorf("no url given"))
		}
		err = checkURL(url)
		if err != nil {
			stderr(err)
		}
	}

	var dlna *goupnp.MaybeRootDevice
	if *dev == "" {
		dlna = chooseDLNADevice()
	} else {
		dlna, err = findDLNAreceiver(*dev)
		if err != nil {
			stderr(err)
		}
	}

	if dlna.Location != nil && !*stop {
		av1SetAndPlay(dlna.Location, url)
		return
	}

	if dlna.Location != nil && *stop {
		av1Stop(dlna.Location)
		return
	}

	stderr(fmt.Errorf("device not ready"))
}

func chooseDLNADevice() *goupnp.MaybeRootDevice {
	fmt.Println("Loading...")

	roots, err := goupnp.DiscoverDevices(av1.URN_AVTransport_1)

	fmt.Print("\033[1A\033[K")
	fmt.Println("----------")

	stderr(err)

	if len(roots) == 0 {
		stderr(fmt.Errorf("no dlna devices on the network found"))
	}
	fmt.Println("DLNA receivers")

	for i, v := range roots {
		fmt.Printf("%d: %s\n", i, v.Root.Device.FriendlyName)
	}

	fmt.Println("----------")
	fmt.Println("Select the DLNA device:")

	selected := selector(roots)
	return &roots[selected]
}

func findDLNAreceiver(friedlyName string) (*goupnp.MaybeRootDevice, error) {
	roots, err := goupnp.DiscoverDevices(av1.URN_AVTransport_1)
	if err != nil {
		return nil, err
	}
	for _, root := range roots {
		if root.Root.Device.FriendlyName == friedlyName {
			return &root, nil
		}
	}
	return nil, fmt.Errorf("'%s' not found", friedlyName)
}

func selector[slice any](s []slice) int {
	var choice int
	for {
		var choiceStr string
		_, err := fmt.Fscanln(os.Stdin, &choiceStr)
		if err != nil {
			fmt.Print("\033[1A\033[K")
			continue
		}
		choice, err = strconv.Atoi(choiceStr)
		if err != nil || choice >= len(s) {
			fmt.Print("\033[1A\033[K")
		} else {
			break
		}
	}
	fmt.Print("\033[1A\033[K")
	fmt.Printf("[%d]\n", choice)
	return choice
}

func stderr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func av1SetAndPlay(loc *url.URL, stream string) {
	client, err := av1.NewAVTransport1ClientsByURL(loc)
	stderr(err)
	err = client[0].SetAVTransportURI(0, stream, "")
	stderr(err)
	err = client[0].Play(0, "1")
	stderr(err)
}

func av1Stop(loc *url.URL) {
	client, err := av1.NewAVTransport1ClientsByURL(loc)
	if err != nil {
		return
	}
	client[0].Stop(0)
}

func checkURL(link string) error {
	parsed, err := url.Parse(link)
	if err != nil {
		return err
	}

	if parsed.Scheme == "" {
		return fmt.Errorf("no url scheme")
	}

	if parsed.Host == "" {
		return fmt.Errorf("no host")
	}
	return nil
}
