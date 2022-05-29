// Package slicefmt helps format slices.
//
// The Config describes how to format slices but not their elements.
// Given a Config, cfg, and a slice, vs, it can be used as such
//
// 	fmt.Printf("example: %q\n", cfg.Fmt(vs))
//
// In this example the %q verb is applied to each of the elements of vs
// while cfg handles how to print separators between elements, among other things.
package slicefmt

import (
	"fmt"
	"io"
	"reflect"
)

// A SepFunc writes separators based on position.
// The parameter n, between 0 and last inclusive,
// represents the nth separator to format not the nth slice element.
type SepFunc = func(w io.Writer, n, last int)

// A SummaryFunc is called with n>0 for the n elements after the cut off.
type SummaryFunc = func(w io.Writer, n int)

// Config describes how to format a slice.
type Config struct {
	// Prefix and Postfix are added on slices of len > 0,
	// unless otherwise noted.
	Prefix, Postfix string

	// Nil is returned for a slice that == nil.
	// Len0 is returned for a slice of len 0 that != nil.
	// If Empty != "", it will be used for both cases.
	Nil   string
	Len0  string
	Empty string

	// If SepFunc is not nil, it is called between formatting elements.
	// Otherwise Sep is used.
	Sep     string
	SepFunc SepFunc

	// If CutOff > 0, only up to the first CutOff elements will be formatted
	// and one of Summary or SummaryFunc is used instead of Postfix.
	// If SummaryFunc is nil, Summary will be returned after the cutoff.
	// No separator is inserted between the last element and the summary.
	// Namely, if SepFunc is not nil it will be called at most CutOff times.
	CutOff      int
	Summary     string
	SummaryFunc SummaryFunc
}

func print(w io.Writer, s string) {
	if s != "" {
		fmt.Fprint(w, s)
	}
}

func (c *Config) lengths(len int) (int, int) {
	if c.CutOff != 0 && len > c.CutOff {
		return c.CutOff, len - c.CutOff
	}
	return len, 0
}

func (c *Config) fmtPrefix(w io.Writer) {
	print(w, c.Prefix)
}

func (c *Config) fmtEmpty(w io.Writer, isNil bool) {
	s := c.Empty
	if s == "" {
		if isNil {
			s = c.Nil
		} else {
			s = c.Len0
		}
	}
	print(w, s)
}

func (c *Config) fmtSep(w io.Writer, n, last int) {
	if c.SepFunc != nil {
		c.SepFunc(w, n, last)
	} else {
		print(w, c.Sep)
	}
}

func (c *Config) fmtEnd(w io.Writer, n int) {
	if n == 0 {
		print(w, c.Postfix)
	} else if c.SummaryFunc != nil {
		c.SummaryFunc(w, n)
	} else {
		print(w, c.Summary)
	}
}

type pair struct {
	c *Config
	v reflect.Value
}

// Fmt takes a slice and returns an object
// that formats the list as specified by c
// while formatting the individual elements
// of the slice by the format string.
//
// The formatter is bound to the config and the slice. It should not be reused.
//
// It should always be used immediately, like:
//
//	fmt.Printf("%q\n", myConfig.Fmt(mySlice))
//	fmt.Println(myConfig.Fmt(mySlice))
func (c *Config) Fmt(slice any) fmt.Formatter {
	return pair{
		c: c,
		v: reflect.ValueOf(slice),
	}
}

func (p pair) Format(f fmt.State, verb rune) {
	cfg, v := p.c, p.v

	// Handle any obscure error cases and bail
	if cfg == nil {
		fmt.Fprint(f, "%!(slice formatter given nil *Config)")
		return
	}
	if !v.IsValid() || v.Type().Kind() != reflect.Slice {
		fmt.Fprint(f, "%!(slice formatter only formats slices)")
		return
	}

	L := v.Len()
	if L == 0 {
		cfg.fmtEmpty(f, v.IsNil())
		return
	}

	N, leftovers := cfg.lengths(L)

	cfg.fmtPrefix(f)
	for i := 0; i < N; i++ {
		if i != 0 {
			cfg.fmtSep(f, i-1, N-2)
		}

		fmt.Fprintf(f, fmtFormatString(f, verb), v.Index(i))
	}
	cfg.fmtEnd(f, leftovers)
}
