// MIT License | ugjka@proton.me
// Write to EC register
// Might need 'modprobe ec_sys write_support=1'
// Values must be Hex
// Example: ecwrite e3 ff
package main

import (
	"fmt"
	"os"
)

const ec = "/sys/kernel/debug/ec/ec0/io"

func main() {

	register := os.Args[1]
	var r uint8
	if _, err := fmt.Sscanf(register, "%x", &r); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid register %s\n", register)
		os.Exit(1)
	}
	value := os.Args[2]
	var v uint8
	if _, err := fmt.Sscanf(value, "%x", &v); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid value %s\n", value)
		os.Exit(1)
	}
	fd, err := os.OpenFile(ec, os.O_RDWR, 0600)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_, err = fd.WriteAt([]byte{byte(v)}, int64(r))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
