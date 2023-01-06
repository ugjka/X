// MIT License | ugjka@proton.me
// Dump laptop's EC content
// Might need modprobe ec_sys
package main

import (
	"fmt"
	"os"
)

const ec = "/sys/kernel/debug/ec/ec0/io"

const reset = "\033[0m"
const red = "\033[31m"
const green = "\033[32m"

func main() {
	fd, err := os.Open(ec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	buf := make([]byte, 1)
	for i := 0; i < 16; i++ {
		fmt.Fprintf(os.Stdout, " %s%02X%s", green, i*16, reset)
		for j := 0; j < 16; j++ {
			fd.Read(buf)
			if int(buf[0]) == 0 {
				fmt.Print(" 00")
				continue
			}
			fmt.Fprintf(os.Stdout, " %s%02X%s", red, buf[0], reset)
		}
		fmt.Println()
	}
	fmt.Print("   ")
	for i := 0; i < 16; i++ {
		fmt.Fprintf(os.Stdout, " %s%02X%s", green, i, reset)
	}
	fmt.Println()
}
