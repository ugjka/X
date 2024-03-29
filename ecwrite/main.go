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
