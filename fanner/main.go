// For Toshiba SATELLITE L850-1LK
// Toggle on and off active cooling
// Might need 'modprobe ec_sys write_support=1'
package main

import (
	"fmt"
	"log"
	"os"
)

const ec = "/sys/kernel/debug/ec/ec0/io"
const register = 0xE3

func main() {
	var arg string
	if arg = os.Args[1]; len(os.Args) < 2 && arg != "on" && arg != "off" {
		fmt.Fprintln(os.Stderr, "Options are 'on' or 'off'")
		return
	}
	fd, err := os.OpenFile(ec, os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}
	if arg == "on" {
		fd.WriteAt([]byte{byte(255)}, register)
	} else {
		fd.WriteAt([]byte{byte(0)}, register)
	}
}
