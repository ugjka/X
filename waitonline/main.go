// Wait for machine come online
// Usage: waitonline sshmachine.local 22
package main

import (
	"net"
	"os"
	"time"
)

func main() {
	remote := os.Args[1]
	port := os.Args[2]
	for {
		if _, err := net.Dial("tcp", remote+":"+port); err == nil {
			break
		}
		time.Sleep(time.Second)
	}
}
