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
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

var (
	serverAddr string
	once       sync.Once
	api        *Client
)

func TestMain(m *testing.M) {
	// log.Println("Do stuff BEFORE the tests!")
	http.HandleFunc("/events", getEvents)
	http.HandleFunc("/events/9a96ba92-46b4-4e41-bcc0-fb273dbf22b7", getEvent)
	http.HandleFunc("/guests", getGuests)

	once.Do(startServer)
	api, _ = New("testing-client-id", "testing-client-secret", OptionDebug(false), OptionEnvFile("./.env"))

	exitVal := m.Run()

	os.Exit(exitVal)
}

func startServer() {
	http.HandleFunc("/realms/users/protocol/openid-connect/token", getToken)
	server := httptest.NewServer(nil)
	serverAddr = server.Listener.Addr().String()
	log.Print("Test WebSocket server listening on ", serverAddr)
}

func getToken(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`
{
	"access_token":"e.e.Y-s-d-o",
	"expires_in":300,
	"refresh_expires_in":0,
	"token_type":"Bearer",
	"not-before-policy":0,
	"scope":"profile email"
}`)
	rw.Write(response)
}
