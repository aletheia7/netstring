// Copyright 2016 aletheia7. All rights reserved. Use of this source code is
// governed by a BSD-2-Clause license that can be found in the LICENSE file.
//
// Package netstring provides methods to make netstrings and parse netstrings
// from a stream using a bufio.Scanner.
//
// Netstring "11:Hello World," = "Hello World"

package netstring

import (
	"bufio"
	"io"
	"strconv"
)

const (
	comma = byte(44)
	zero  = byte(48)
	nine  = byte(57)
	colon = byte(58)
)

// NewScanner returns a *bufioScanner that parses netstrings.
// The scanner will skip over invalid netstrings in a stream.
//
func NewScanner(r io.Reader) *bufio.Scanner {
	var (
		s         = bufio.NewScanner(r)
		ns_len    int
		colon_pos int
		comma_pos int
	)
	s.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if ns_len == 0 {
			for i, b := range data {
				switch {
				case zero <= b && b <= nine:
				case b == colon:
					colon_pos = i
					var num_err error
					if ns_len, num_err = strconv.Atoi(string(data[:colon_pos])); num_err != nil {
						advance = colon_pos + 1
						colon_pos = 0
						return
					}
					break
				default:
					advance = i + 1
					return
				}
			}
			return
		}
		comma_pos = colon_pos + ns_len + 1
		if comma_pos < len(data) {
			if data[comma_pos] == comma {
				token = data[colon_pos+1 : comma_pos]
			}
			advance = comma_pos + 1
			ns_len = 0
			colon_pos = 0
		}
		return
	})
	return s
}

// String to netstring
//
func S2nsb(s string) []byte {
	return B2nsb([]byte(s))
}

// String to netstring
//
func S2ns(s string) string {
	return string(B2nsb([]byte(s)))
}

// Bytes to netstring
//
func B2ns(b []byte) string {
	return string(B2nsb(b))
}

// Bytes to netstring
//
func B2nsb(b []byte) []byte {
	nslen := strconv.AppendInt(nil, int64(len(b)), 10)
	r := make([]byte, len(nslen)+1+len(b)+1)
	copy(r, nslen)
	r[len(nslen)] = colon
	copy(r[len(nslen)+1:], b)
	r[len(r)-1] = comma
	return r
}
