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
