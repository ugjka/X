// For Toshiba SATELLITE L850-1LK
// Toggle on and off active cooling
// Might need 'modprobe ec_sys write_support=1'
package main

import (
	"fmt"
	"os"
)

const ec = "/sys/kernel/debug/ec/ec0/io"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Options are 'on', 'off' or 'blast'")
		os.Exit(1)
	}
	fd, err := os.OpenFile(ec, os.O_RDWR, 0600)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	switch os.Args[1] {
	case "on":
		fd.WriteAt([]byte{byte(255)}, 0xE3)
		fd.WriteAt([]byte{0x8}, 0xed)
	case "off":
		fd.WriteAt([]byte{byte(0)}, 0xE3)
	case "blast":
		fd.WriteAt([]byte{byte(255)}, 0xE3)
		fd.WriteAt([]byte{0x12}, 0xed)
	default:
		fmt.Fprintln(os.Stderr, "Options are 'on', 'off' or 'blast'")
		os.Exit(1)
	}
}
