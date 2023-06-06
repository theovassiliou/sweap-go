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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// SweapResponse handles parsing out errors from the web api
type SweapResponse struct {
	Ok               bool             `json:"ok,omitempty"`
	Error            SweapError       `json:"error,omitempty"`
	ResponseMetadata ResponseMetadata `json:"response_metadata,omitempty"`
}

type SweapError struct {
	Error_  string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (se *SweapError) Error() string {
	return fmt.Sprintf("Code %v: %v - %v", se.Code, se.Error_, se.Message)
}

func jsonReq(ctx context.Context, endpoint string, body interface{}) (req *http.Request, err error) {
	buffer := bytes.NewBuffer([]byte{})
	if err = json.NewEncoder(buffer).Encode(body); err != nil {
		return nil, err
	}

	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buffer); err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	return req, nil
}

func parseResponseBody(body io.ReadCloser, intf interface{}, d Debug) error {
	response, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	if d.Debug() {
		d.Debugln("parseResponseBody", string(response))
	}

	return json.Unmarshal(response, intf)
}

func logRequest(req *http.Request, d Debug) error {
	if d.Debug() {
		text, err := httputil.DumpRequest(req, true)
		if err != nil {
			return err
		}
		d.Debugln(string(text))
	}

	return nil
}

func doPost(ctx context.Context, client httpClient, req *http.Request, parser responseParser, d Debug) error {
	logRequest(req, d)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = checkStatusCode(resp, d)
	if err != nil {
		sweapError := &SweapError{}
		e := newJSONParser(sweapError)(resp)
		if e == nil {
			return sweapError
		}
		return err
	}
	logResponse(resp, d)
	if resp.ContentLength != 0 {
		return parser(resp)
	}
	return nil
}

func doDelete(ctx context.Context, client httpClient, req *http.Request, d Debug) error {
	logRequest(req, d)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = checkStatusCode(resp, d)
	if err != nil {
		sweapError := &SweapError{}
		e := newJSONParser(sweapError)(resp)
		if e == nil {
			return sweapError
		}
		return err
	}

	return nil
}

// post JSON.
func postJSON(ctx context.Context, client httpClient, endpoint string, json []byte, intf interface{}, d Debug) error {
	reqBody := bytes.NewBuffer(json)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	return doPost(ctx, client, req, newJSONParser(intf), d)
}

// putJSON.
func putJSON(ctx context.Context, client httpClient, endpoint string, json []byte, intf interface{}, d Debug) error {
	reqBody := bytes.NewBuffer(json)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	return doPost(ctx, client, req, newJSONParser(intf), d)
}

func getResource(ctx context.Context, client httpClient, endpoint string, values url.Values, intf interface{}, d Debug) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	req.URL.RawQuery = values.Encode()

	return doPost(ctx, client, req, newJSONParser(intf), d)
}

func deleteResource(ctx context.Context, client httpClient, endpoint string, values url.Values, d Debug) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	req.URL.RawQuery = values.Encode()

	return doDelete(ctx, client, req, d)
}

func logResponse(resp *http.Response, d Debug) error {
	if d.Debug() {
		text, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return err
		}
		d.Debugln(string(text))
	}

	return nil
}

// timerReset safely reset a timer, see time.Timer.Reset for details.
func timerReset(t *time.Timer, d time.Duration) {
	if !t.Stop() {
		<-t.C
	}
	t.Reset(d)
}

func checkStatusCode(resp *http.Response, d Debug) error {
	// Sweap seems to send an HTML body along with 5xx error codes. Don't parse it.
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		logResponse(resp, d)
		return StatusCodeError{Code: resp.StatusCode, Status: resp.Status}
	}

	return nil
}

type responseParser func(*http.Response) error

func newJSONParser(dst interface{}) responseParser {
	return func(resp *http.Response) error {
		return json.NewDecoder(resp.Body).Decode(dst)
	}
}
