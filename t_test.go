// Copyright 2016 aletheia7. All rights reserved. Use of this source code is
// governed by a BSD-2-Clause license that can be found in the LICENSE file.

package netstring_test

import (
	"bytes"
	"netstring"
	"strings"
	"testing"
	"testing/iotest"
)

func TestString2_and_Bytes2_netstring(t *testing.T) {
	s := "this is a long string"
	ns := `21:` + s + `,`
	if r := netstring.S2nsb(s); !bytes.Equal(r, []byte(ns)) {
		t.Error(ns, string(r))
	}
	if r := netstring.S2ns(s); r != ns {
		t.Error(ns, r)
	}
	if r := netstring.B2ns([]byte(s)); r != ns {
		t.Error(ns, r)
	}
	if r := netstring.B2nsb([]byte(s)); !bytes.Equal(r, []byte(ns)) {
		t.Error(ns, string(r))
	}
}

func TestScanner(t *testing.T) {
	s1 := "abc"
	s2 := strings.Repeat("z", 20000)
	scanner := netstring.NewScanner(iotest.OneByteReader(strings.NewReader(netstring.S2ns(s1) + netstring.S2ns(s2))))
	ret := ""
	for scanner.Scan() {
		ret += scanner.Text()
	}
	if ret != s1+s2 {
		t.Error("failed: ret != s1 + s2")
	}
}

func TestScanner_with_bad(t *testing.T) {
	answer := "abcdef" + "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"
	scanner := netstring.NewScanner(iotest.OneByteReader(strings.NewReader(netstring.S2ns(`abc`) +
		`abc:abc,` + netstring.S2ns("def") + `2:,` + netstring.S2ns("def") + netstring.S2ns(strings.Repeat("z", 300)))))
	result := ``
	for scanner.Scan() {
		result += scanner.Text()
	}
	if answer != result {
		t.Error("failed answer != result")
	}
}
