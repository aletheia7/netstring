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
	"bytes"
	"errors"
	"io"
	"strconv"
)

const (
	comma = byte(44)
	zero  = byte(48)
	nine  = byte(57)
	colon = byte(58)
)

var Err_invalid_netstring = errors.New("invalid netstrng")

type scanner struct {
	ns_len    int
	colon_pos int
	comma_pos int
}

func (o *scanner) Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if o.ns_len == 0 {
		for i, b := range data {
			switch {
			case zero <= b && b <= nine:
			case b == colon:
				o.colon_pos = i
				var num_err error
				if o.ns_len, num_err = strconv.Atoi(string(data[:o.colon_pos])); num_err != nil {
					advance = o.colon_pos + 1
					o.colon_pos = 0
					return
				}
				goto good_length
			default:
				advance = i + 1
				return
			}
		}
		return
	}
good_length:
	o.comma_pos = o.colon_pos + o.ns_len + 1
	if o.comma_pos < len(data) {
		if data[o.comma_pos] == comma {
			token = data[o.colon_pos+1 : o.comma_pos]
		}
		advance = o.comma_pos + 1
		o.ns_len = 0
		o.colon_pos = 0
	}
	return
}

// NewScanner returns a *bufioScanner that parses netstrings.
// The scanner will skip over invalid netstrings in a stream.
//
func NewScanner(r io.Reader) *bufio.Scanner {
	s := bufio.NewScanner(r)
	s.Split((&scanner{}).Split)
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

type reader struct {
	scanner *bufio.Scanner
}

func (o *reader) Read(p []byte) (n int, err error) {
	if o.scanner.Scan() {
		b := o.scanner.Bytes()
		n, err = io.ReadAtLeast(bytes.NewBuffer(b), p, len(b))
	} else {
		err = io.EOF
	}
	return
}

// Reader returns an io.Reader that decodes netstrings
//
func Reader(r io.Reader) io.Reader {
	return &reader{NewScanner(r)}
}

type writer struct {
	w io.Writer
}

func (o *writer) Write(p []byte) (n int, err error) {
	n, err = o.w.Write(B2nsb(p))
	return
}

// Writer returns a io.Writer that makes netstrings
//
func Writer(w io.Writer) io.Writer {
	return &writer{w}
}
