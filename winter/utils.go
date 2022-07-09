package winter

import "strings"

// pad returns a function that pads every string given to it to the longest
// length it has seen so far.
//
// For example:
//
//     p := pad()
//     pad("hi")    // "hi"
//     pad("hello") // "hello"
//     pad("hi")    // "hi   "
//     pad("hello") // "hello"
func pad() func(string) string {
	var longest int
	return func(s string) string {
		if len(s) >= longest {
			longest = len(s)
			return s
		}
		return s + strings.Repeat(" ", longest-len(s))
	}
}
