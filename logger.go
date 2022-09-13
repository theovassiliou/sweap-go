/*
 Copyright (c) 2022 Theo Vassiliou

 Permission is hereby granted, free of charge, to any person obtaining a copy of
 this software and associated documentation files (the "Software"), to deal in
 the Software without restriction, including without limitation the rights to
 use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 the Software, and to permit persons to whom the Software is furnished to do so,
 subject to the following conditions:

 The above copyright notice and this permission notice shall be included in all
 copies or substantial portions of the Software.

 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package sweap

// Based on logger.go from https://github.com/slack-go/slack

import (
	"fmt"
)

// logger is a logger interface compatible with both stdlib and some
// 3rd party loggers.
type logger interface {
	Output(int, string) error
}

// ilogger represents the internal logging api we use.
type ilogger interface {
	logger
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
}

type Debug interface {
	Debug() bool

	// Debugf print a formatted debug line.
	Debugf(format string, v ...interface{})
	// Debugln print a debug line.
	Debugln(v ...interface{})
}

// internalLog implements the additional methods used by our internal logging.
type internalLog struct {
	logger
}

// Println replicates the behaviour of the standard logger.
func (t internalLog) Println(v ...interface{}) {
	t.Output(2, fmt.Sprintln(v...))
}

// Printf replicates the behaviour of the standard logger.
func (t internalLog) Printf(format string, v ...interface{}) {
	t.Output(2, fmt.Sprintf(format, v...))
}

// Print replicates the behaviour of the standard logger.
func (t internalLog) Print(v ...interface{}) {
	t.Output(2, fmt.Sprint(v...))
}

type discard struct{}

func (t discard) Debug() bool {
	return false
}

// Debugf print a formatted debug line.
func (t discard) Debugf(format string, v ...interface{}) {}

// Debugln print a debug line.
func (t discard) Debugln(v ...interface{}) {}
