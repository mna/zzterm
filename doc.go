// Package zzterm efficiently reads and decodes terminal input keys and mouse
// events without any memory allocation. It is intended to be used with a
// terminal set in raw mode as its io.Reader. Setting the terminal in raw
// mode is not handled by this package, there are a number of Go packages that
// can do this (see the example).
//
//
package zzterm // import "git.sr.ht/~mna/zzterm"
