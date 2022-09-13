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

// Based on logger_test.go from https://github.com/slack-go/slack

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	buf := bytes.NewBufferString("")
	logger := internalLog{logger: log.New(buf, "", 0|log.Lshortfile)}
	logger.Println("test line 123")
	assert.Equal(t, buf.String(), "logger_test.go:37: test line 123\n")
	buf.Truncate(0)
	logger.Print("test line 123")
	assert.Equal(t, buf.String(), "logger_test.go:40: test line 123\n")
	buf.Truncate(0)
	logger.Printf("test line 123\n")
	assert.Equal(t, buf.String(), "logger_test.go:43: test line 123\n")
	buf.Truncate(0)
	logger.Output(1, "test line 123\n")
	assert.Equal(t, buf.String(), "logger_test.go:46: test line 123\n")
	buf.Truncate(0)
}
