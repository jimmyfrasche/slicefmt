package slicefmt_test

import (
	"fmt"
	"io"

	"github.com/jimmyfrasche/slicefmt"
)

func Example_default() {
	// This cfg replicates the default slice formatting.
	// Not especially useful!
	cfg := &slicefmt.Config{
		Empty:   "[]",
		Prefix:  "[",
		Postfix: "]",
		Sep:     " ",
	}

	// Empty is used for both cases when the length is 0.
	var nil []float64
	fmt.Println(cfg.Fmt(nil))
	fmt.Println(cfg.Fmt([]byte{}))

	fmt.Println(cfg.Fmt([]int{1, 2, 3}))

	// Output:
	// []
	// []
	// [1 2 3]
}

func Example_simple() {
	cfg := slicefmt.Config{
		Nil:  "<nil>",
		Len0: "<none>",
		Sep:  " ~ ",
	}

	// Separate formatting of the two 0 length cases
	var nil []float64
	fmt.Println(cfg.Fmt(nil))

	fmt.Printf("%q\n", cfg.Fmt([]byte{}))

	fmt.Printf("%02d\n", cfg.Fmt([]int{1, 2, 3}))

	// Output:
	// <nil>
	// <none>
	// 01 ~ 02 ~ 03
}

func Example_commas() {
	notOxford := slicefmt.Config{
		SepFunc: func(w io.Writer, n, end int) {
			// Use ", " for all but the last separator
			sep := ", "
			if n == end {
				sep = " and "
			}
			fmt.Fprint(w, sep)
		},
	}

	oxford := slicefmt.Config{
		SepFunc: func(w io.Writer, n, end int) {
			// Default to ", "
			sep := ", "
			if n == 0 && n == end {
				// If there are only 2 items (1 space) just use " and "
				sep = " and "
			} else if n == end {
				// If there are more than two items and this is the last
				sep = ", and "
			}
			fmt.Fprint(w, sep)
		},
	}

	var s []int
	for i := 0; i < 4; i++ {
		s = append(s, i)
		fmt.Println(notOxford.Fmt(s))
		fmt.Println(oxford.Fmt(s))
	}

	// Output:
	// 0
	// 0
	// 0 and 1
	// 0 and 1
	// 0, 1 and 2
	// 0, 1, and 2
	// 0, 1, 2 and 3
	// 0, 1, 2, and 3
}

func Example_summary() {
	simple := slicefmt.Config{
		Sep:     "/",
		CutOff:  2,
		Summary: " and so on",
	}
	complex := &slicefmt.Config{
		Sep:    "/",
		CutOff: 2,
		SummaryFunc: func(w io.Writer, n int) {
			fmt.Fprintf(w, " (and %d more)", n)
		},
	}

	var s []int
	for i := 0; i < 4; i++ {
		s = append(s, i)
		fmt.Println(simple.Fmt(s))
		fmt.Println(complex.Fmt(s))
	}

	// Output:
	// 0
	// 0
	// 0/1
	// 0/1
	// 0/1 and so on
	// 0/1 (and 1 more)
	// 0/1 and so on
	// 0/1 (and 2 more)
}

func Example_complex() {
	// SepFunc and SummaryFunc work well together.
	cfg := slicefmt.Config{
		Empty:  "none",
		CutOff: 3,
		SepFunc: func(w io.Writer, n, end int) {
			// Use ", " for all but the last separator
			sep := ", "
			if n == end {
				sep = " and "
			}
			fmt.Fprint(w, sep)
		},
		SummaryFunc: func(w io.Writer, n int) {
			fmt.Fprintf(w, " (and %d more)", n)
		},
	}

	s := []string{"a", "b", "c", "d", "e", "f"}
	for i := range s {
		fmt.Printf("%q\n", cfg.Fmt(s[:i]))
	}

	// Output:
	// none
	// "a"
	// "a" and "b"
	// "a", "b" and "c"
	// "a", "b" and "c" (and 1 more)
	// "a", "b" and "c" (and 2 more)
}

func Example_bad() {
	// a nil *slicefmt.Config is invalid
	var badConfig *slicefmt.Config
	fmt.Println(badConfig.Fmt(nil))

	// a zero slicefmt.Config is valid
	// but Fmt must be called with a slice
	var goodConfig slicefmt.Config
	fmt.Println(goodConfig.Fmt(nil))
	fmt.Println(goodConfig.Fmt(7))

	// Output:
	// %!(slice formatter given nil *Config)
	// %!(slice formatter only formats slices)
	// %!(slice formatter only formats slices)
}
