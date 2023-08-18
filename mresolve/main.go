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
// Resolve mdns domains e.g "mresolve raspberrypi.local"
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/pion/mdns"
	"golang.org/x/net/ipv4"
)

func main() {
	if len(os.Args) < 2 {
		stderr(fmt.Errorf("error: no domain supplied"))
	}

	addr, err := net.ResolveUDPAddr("udp", mdns.DefaultAddress)
	stderr(err)

	l, err := net.ListenUDP("udp4", addr)
	stderr(err)

	server, err := mdns.Server(ipv4.NewPacketConn(l), &mdns.Config{})
	stderr(err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, src, err := server.Query(ctx, os.Args[1])
	stderr(err)

	fmt.Println(src)
}

func stderr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
