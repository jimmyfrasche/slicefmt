//go:build !go1.20

package slicefmt

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

// This should land as fmt.FormatString in Go 1.19
//
// Copied out of https://go-review.googlesource.com/c/go/+/400875
func fmtFormatString(state fmt.State, verb rune) string {
	var tmp [16]byte // Use a local buffer.
	b := append(tmp[:0], '%')
	for _, c := range " +-#0" { // All known flags
		if state.Flag(int(c)) { // The argument is an int for historical reasons.
			b = append(b, byte(c))
		}
	}
	if w, ok := state.Width(); ok {
		b = strconv.AppendInt(b, int64(w), 10)
	}
	if p, ok := state.Precision(); ok {
		b = append(b, '.')
		b = strconv.AppendInt(b, int64(p), 10)
	}
	b = utf8.AppendRune(b, verb)
	return string(b)
}
