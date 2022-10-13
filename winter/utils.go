package winter

import "strings"

// newPad returns a function that pads every string given to it to the longest
// length it has seen so far.
//
// For example:
//
//	p := newPad()
//	p("hi")    // "hi"
//	p("hello") // "hello"
//	p("hi")    // "hi   "
//	p("hello") // "hello"
//	q := newPad()
//	q("hi")    // "hi"
func newPad() func(string) string {
	var longest int
	return func(s string) string {
		if len(s) >= longest {
			longest = len(s)
			return s
		}
		return s + strings.Repeat(" ", longest-len(s))
	}
}
