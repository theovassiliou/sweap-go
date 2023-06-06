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
	"time"
)

const messageBufferSize int = 256

type ListenOptions func(*Client)

// Listen listen to guest changes (new, updates)
func (api *Client) Listen() (chan *GuestUpdate, error) {

	messageCh := make(chan *GuestUpdate, messageBufferSize)

	search := EventSearchParameter{
		Name: "Integration Test Event",
	}

	events, err := api.SearchEvents(search)
	guests, err := api.GetGuests((*events)[0].ID)

	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}
	go func() {
		storedGuests := make([]Guest, 0)

		for {

			for _, g := range *guests {

				if !contains(g, storedGuests) {

					messageCh <- &GuestUpdate{NEWGUEST, g}
				} else {

				}
			}
			storedGuests = *guests
			fmt.Println("[INFO] Going to sleep")
			time.Sleep(15 * time.Second)
			guests, _ = api.GetGuests((*events)[0].ID)

		}
	}()

	return messageCh, nil
}

func contains(g1 Guest, g2 Guests) bool {
	for _, g := range g2 {
		if g1.ID == g.ID {
			return true
		}
	}
	return false
}
