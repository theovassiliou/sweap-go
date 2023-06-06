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

	"github.com/joho/godotenv"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	// Official endpoints for normal users
	// APIURL   = "https://api.sweap.io/core/v1/"
	APIURL   = "https://api.sweap.io/core/v1/"
	TOKENURL = "https://auth.sweap.io/realms/users/protocol/openid-connect/token"

	// Inofficial endpoints to the staging or dev area.
	// Normal API users won't need this.
	// For special development purposes, only
	APIURL_ST   = "https://api.sweap.st/core/v1/"
	TOKENURL_ST = "https://auth.sweap.st/realms/users/protocol/openid-connect/token"

	APIURL_DEV   = "https://api.sweap.dev/core/v1/"
	TOKENURL_DEV = "https://auth.sweap.dev/realms/users/protocol/openid-connect/token"
)

type Status string

const (
	OK                   Status = "OK"
	BAD_REQUEST          Status = "BAD_REQUEST"
	UNAUTHORIZED         Status = "UNAUTHORIZED"
	VALIDATION_EXCEPTION Status = "VALIDATION_EXCEPTION"
	EXCEPTION            Status = "EXCEPTION"
	WRONG_CREDENTIALS    Status = "WRONG_CREDENTIALS"
	ACCESS_DENIED        Status = "ACCESS_DENIED"
	NOT_FOUND            Status = "NOT_FOUND"
	DUPLICATE_ENTITY     Status = "DUPLICATE_ENTITY"
)

type SortBy string

const (
	CREATED_AT SortBy = "createdAt"
	UPDATED_AT SortBy = "updatedAt"
	START_DATE SortBy = "startDate"
	END_DATE   SortBy = "endDate"
	NAME       SortBy = "name"
)

type SortOrder string

const (
	ASC  SortOrder = "ASC"
	DESC SortOrder = "DESC"
)

type Error struct {
	Error   string
	Code    int
	Message string
}

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
	envFile      string
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
		ClientID:     s.clientID,
		ClientSecret: s.clientSecret,
		TokenURL:     s.tokenurl,
		AuthStyle:    0,
	}

	ctx := context.Background()
	client := c.Client(ctx)

	s.httpclient = client

	if c.TokenURL == TOKENURL_DEV && !s.checkCredentials() {
		return nil, fmt.Errorf("authorization at endopoint %v failed. Check credentials", s.tokenurl)
	}

	return s, nil
}

func (s *Client) checkCredentials() bool {
	_, err := s.CheckCredentials()
	return err == nil
}

// ----- Options -------
type SweapOptions func(*Client)

// OptionClientCredentials sets an app-level token for the client.
func OptionClientCredentials(id, secret string) func(*Client) {
	return func(c *Client) {
		c.clientID = id
		c.clientSecret = secret
	}
}

// OptionUseStagingEnv selects the staging environment for the client.
func OptionUseStagingEnv() func(*Client) {
	return func(c *Client) {
		c.tokenurl = TOKENURL_ST
		c.endpoint = APIURL_ST
	}
}

// OptionUseDevEnv selects the dev environment for the client.
func OptionUseDevEnv() func(*Client) {
	return func(c *Client) {
		c.tokenurl = TOKENURL_DEV
		c.endpoint = APIURL_DEV
	}
}

// OptionLog set logging for client.
func OptionLog(l logger) func(*Client) {
	return func(c *Client) {
		c.log = internalLog{logger: l}
	}
}

// OptionEnvFile uses an env file for reading the credentials
func OptionEnvFile(fileName string) func(*Client) {
	return func(c *Client) {
		c.envFile = fileName
		err := godotenv.Load(fileName)

		if err != nil {
			log.Fatalf("Error loading .env file: %v\n", fileName)
		}

		if env, present := os.LookupEnv("CLIENTID"); present {
			c.clientID = env
		}

		if env, present := os.LookupEnv("CLIENT_SECRET"); present {
			c.clientSecret = env
		}
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
