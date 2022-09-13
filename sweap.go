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
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/oauth2/clientcredentials"
)

const (
	// APIURL of the slack api.
	APIURL   = "https://api.sweap.io/core/v1/"
	TOKENURL = "https://auth.sweap.io/realms/users/protocol/openid-connect/token"
)

// httpClient defines the minimal interface needed for an http.Client to be implemented.
type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// AuthTestResponse ...
type AuthTestResponse struct {
	URL string `json:"url"`
}

// ResponseMetadata holds pagination metadata
type ResponseMetadata struct {
	Cursor   string   `json:"next_cursor"`
	Messages []string `json:"messages"`
	Warnings []string `json:"warnings"`
}

// Client for the Sweap API
type Client struct {
	clientID     string
	clientSecret string
	authStyle    int
	debug        bool
	endpoint     string
	tokenurl     string
	log          ilogger
	httpclient   httpClient
}

// NewSweap creates a new Sweap object with given credentials
// if authentication fails error will be non-nil
func New(clientId, clientSecret string, options ...SweapOptions) (*Client, error) {

	s := &Client{
		clientID:     clientId,
		clientSecret: clientSecret,
		authStyle:    0,
		endpoint:     APIURL,
		tokenurl:     TOKENURL,
		log:          log.New(os.Stderr, "theovassiliou/sweap-go", log.LstdFlags|log.Lshortfile),
	}

	for _, opt := range options {
		opt(s)
	}

	c := clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     s.tokenurl,
		AuthStyle:    0,
	}

	ctx := context.Background()
	client := c.Client(ctx)

	s.httpclient = client

	// TODO: Remove this from here, if a real api call for testing authentification has been included.
	/* _, err := s.GetEvents()

	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	} */
	return s, nil
}

// ----- Options -------

// OptionClientCredentials sets an app-level token for the client.
func OptionClientCredentials(id, secret string) func(*Client) {
	return func(c *Client) {
		c.clientID = id
		c.clientSecret = secret
	}
}

type SweapOptions func(*Client)

// OptionLog set logging for client.
func OptionLog(l logger) func(*Client) {
	return func(c *Client) {
		c.log = internalLog{logger: l}
	}
}

// OptionDebug enable debugging for the client
func OptionDebug(b bool) func(*Client) {
	return func(c *Client) {
		c.debug = b
	}
}

// OptionAPIURL set the url for the client. only useful for testing.
func OptionAPIURL(u string) func(*Client) {
	return func(c *Client) { c.endpoint = u }
}

// OptionTOKENURL set the url for the client. only useful for testing.
func OptionTOKENURL(u string) func(*Client) {
	return func(c *Client) { c.tokenurl = u }
}

// ----- Communication -------

// get a Sweap web method.
func (api *Client) getMethod(ctx context.Context, path string, values url.Values, intf interface{}) error {
	return getResource(ctx, api.httpclient, api.endpoint+path, values, intf, api)
}

// ----- Logging ------
// Debugf print a formatted debug line.
func (api *Client) Debugf(format string, v ...interface{}) {
	if api.debug {
		api.log.Output(2, fmt.Sprintf(format, v...))
	}
}

// Debugln print a debug line.
func (api *Client) Debugln(v ...interface{}) {
	if api.debug {
		api.log.Output(2, fmt.Sprintln(v...))
	}
}

// Debug returns if debug is enabled.
func (api *Client) Debug() bool {
	return api.debug
}
