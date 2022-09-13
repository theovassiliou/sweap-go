/*
 Copyright (c) 2022 Theofanis Vassiliou-Gioles

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

import (
	"fmt"
	"net/http"
)

// StatusCodeError represents an http response error.
// type httpStatusCode interface { HTTPStatusCode() int } to handle it.
type StatusCodeError struct {
	Code   int
	Status string
}

func (t StatusCodeError) Error() string {
	return fmt.Sprintf("sweap server error: %s", t.Status)
}

func (t StatusCodeError) HTTPStatusCode() int {
	return t.Code
}

func (t StatusCodeError) Retryable() bool {
	if t.Code >= 500 || t.Code == http.StatusTooManyRequests {
		return true
	}
	return false
}

// SweapLibraryError represents an internal library error
type SweapLibraryError struct {
	Message string
}

func (se SweapLibraryError) Error() string {
	return fmt.Sprintf("sweap internal internal error: %s", se.Message)
}
