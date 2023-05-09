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
// clipper: sync clipboard between 2 machines
// first arg must be the client's address
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

const port = 43344

var client string

func main() {
	client = os.Args[1]
	var current string
	incoming := make(chan string)
	go listener(incoming)
	for {
		select {
		case txt := <-incoming:
			c := exec.Command("xclip", "-selection", "clipboard", "-in", "-d", ":0")
			c.Stdin = strings.NewReader(txt)
			err := c.Run()
			if err != nil {
				log.Printf("Xclip write error: %v", err)
				break
			}
			current = txt
		default:
			time.Sleep(time.Second)
			c := exec.Command("xclip", "-selection", "clipboard", "-out", "-d", ":0")
			data, err := c.Output()
			if err != nil {
				log.Printf("Xclip read error: %v", err)
				break
			}
			txt := string(data)
			if txt == current {
				break
			}
			current = txt
			err = send(txt)
			if err != nil {
				log.Printf("Send error: %v", err)
			}
		}
	}
}

func send(data string) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client, port))
	if err != nil {
		return err
	}
	conn.SetWriteDeadline(time.Now().Add(time.Second * 5))
	_, err = io.WriteString(conn, data)
	conn.Close()
	return err
}

func listener(t chan<- string) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		data, err := ioutil.ReadAll(conn)
		if err != nil {
			conn.Close()
			log.Printf("Read error: %v", err)
			continue
		}
		conn.Close()
		t <- string(data)
	}

}
