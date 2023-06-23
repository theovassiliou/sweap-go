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

// Listen listens to guest changes (new, updates).
// It returns a channel of GuestUpdate objects and an error.
// The channel can be used to receive updates on new or updated guests.
// If an error occurs during the setup, it returns nil channel and the specific error encountered.
// The function can be called on a Client object.
func (api *Client) Listen() (chan *GuestUpdate, error) {
	messageCh := make(chan *GuestUpdate, messageBufferSize)

	search := EventSearchParameter{
		Name: "Integration Test Event",
	}

	events, err := api.SearchEvents(search)
	if err != nil {
		return nil, err
	}

	guests, err := api.GetGuests((*events)[0].ID)
	if err != nil {
		return nil, err
	}

	go func() {
		storedGuests := make([]Guest, 0)

		for {
			for _, g := range *guests {
				if !contains(g, storedGuests) {
					messageCh <- &GuestUpdate{EventType: NEWGUEST, Value: g}
				} else {
					// Handle guest updates here if needed
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

// contains checks if a guest (g1) is present in a slice of guests (g2).
// It returns true if the guest is found, and false otherwise.
func contains(g1 Guest, g2 Guests) bool {
	for _, g := range g2 {
		if g1.ID == g.ID {
			return true
		}
	}
	return false
}

// isUpdatedAtAfter compares two Guest objects and returns true if the UpdatedAt field of the first guest is after the UpdatedAt field of the second guest.
// It takes two Guest objects as input.
// If either of the UpdatedAt fields is nil, it returns false.
// If the UpdatedAt field of the first guest is after the UpdatedAt field of the second guest, it returns true. Otherwise, it returns false.

func isUpdatedAtAfter(guest1, guest2 Guest) bool {
	if guest1.UpdatedAt == nil || guest2.UpdatedAt == nil {
		return false
	}
	return guest1.UpdatedAt.After(*guest2.UpdatedAt)
}
