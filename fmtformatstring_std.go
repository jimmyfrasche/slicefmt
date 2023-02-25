//go:build go1.20

package slicefmt

import "fmt"

func fmtFormatString(state fmt.State, verb rune) string {
	return fmt.FormatString(state, verb)
}
