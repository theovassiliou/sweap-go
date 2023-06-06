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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGuests(t *testing.T) {
	guests, err := api.GetGuests("9a96ba92-46b4-4e41-bcc0-fb273dbf22b7")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	fmt.Println(guests)
	assert.Equal(t, 6, len(*guests))
}

func TestGetGuestById(t *testing.T) {

	guest, err := api.GetGuestById("85291ecc-0b0f-4651-8343-13c1744ea944")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	assert.NotNil(t, guest, "Wanted a guest")
	fmt.Println(guest)
}

func TestGetGuest(t *testing.T) {
	guests, err := api.GetGuests("9a96ba92-46b4-4e41-bcc0-fb273dbf22b7")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	assert.Equal(t, 6, len(*guests))
}

func TestCreateGuest(t *testing.T) {
	guest := Guest{
		ID:              "",
		Version:         0,
		EventID:         "9a96ba92-46b4-4e41-bcc0-fb273dbf22b7",
		FirstName:       "Theo Test",
		LastName:        "Vassiliou",
		EntourageCount:  0,
		Comment:         nil,
		Email:           "vassiliou@web.de",
		CustomFields:    CustomFields{},
		InvitationID:    "",
		InvitationState: NONE,
		TicketID:        "",
	}
	guestCreated, err := api.CreateGuest(guest)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	assert.Equal(t, guest.FirstName, guestCreated.FirstName)
	assert.Equal(t, guest.LastName, guestCreated.LastName)
	assert.NotEqual(t, "", guestCreated.ID)
}

func TestUpdateGuest(t *testing.T) {
	guest := Guest{
		ID:              "",
		Version:         0,
		EventID:         "9a96ba92-46b4-4e41-bcc0-fb273dbf22b7",
		FirstName:       "Theo Test",
		LastName:        "Vassiliou",
		EntourageCount:  0,
		Comment:         nil,
		Email:           "vassiliou@web.de",
		CustomFields:    CustomFields{},
		InvitationID:    "",
		InvitationState: NONE,
		TicketID:        "",
	}
	guestCreated, err := api.CreateGuest(guest)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	assert.Equal(t, guest.LastName, guestCreated.LastName)

	guestNew := *guestCreated
	guestNew.LastName = "Vassiliou-Gioles"

	updatedGuest, err := api.UpdateGuest(guestNew)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	assert.Equal(t, guestNew.LastName, updatedGuest.LastName)
	assert.NotEqual(t, guestCreated.LastName, updatedGuest.LastName)
	err = api.DeleteGuest(updatedGuest.ID)
	assert.Nil(t, err, "err should be nil")
}

func TestDeleteGuest(t *testing.T) {
	guest := Guest{
		ID:              "",
		Version:         0,
		EventID:         "9a96ba92-46b4-4e41-bcc0-fb273dbf22b7",
		FirstName:       "Theo Test",
		LastName:        "Vassiliou",
		EntourageCount:  0,
		Comment:         nil,
		Email:           "vassiliou@web.de",
		CustomFields:    CustomFields{},
		InvitationID:    "",
		InvitationState: NONE,
		TicketID:        "",
	}
	guestCreated, err := api.CreateGuest(guest)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	err = api.DeleteGuest(guestCreated.ID)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	guestNotFound, err := api.GetGuestById(guestCreated.ID)
	assert.Nilf(t, guestNotFound, "guest with guestId %v found but should not", guestCreated.ID)
	assert.NotNil(t, err, "error nil but expected non-nil value")
	assert.Equal(t, 4040, err.(*SweapError).Code)
}

func TestDeleteNonExistingGuest(t *testing.T) {
	err := api.DeleteGuest("ID_NON_EXISTENT")
	assert.NotNil(t, err, "expected not nill error")
	fmt.Println(err)
}

// ------ Aux functions ------

func getGuests(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		response := []byte(`
		{
			"id": "85291ecc-0b0f-4651-8343-13c1744ea944",
			"externalId": null,
			"version": 0,
			"createdAt": "2022-07-06T10:23:34.931Z",
			"updatedAt": "2022-07-06T10:23:34.931Z",
			"eventId": "9a96ba92-46b4-4e41-bcc0-fb273dbf22b7",
			"firstName": "Theo Test",
			"lastName": "Vassiliou",
			"entourageCount": 0,
			"comment": null,
			"email": "vassiliou@web.de",
			"customFields": {},
			"invitationId": "5c59qqypkv6e4s8ue8ay76e9ryqd8742wwgeateacd",
			"invitationState": "ACCEPTED",
			"ticketId": "B44ZZ2XBF69H"
		}	
		`)
		rw.Write(response)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	response := []byte(`[
				{
					"id": "85291ecc-0b0f-4651-8343-13c1744ea944",
					"externalId": null,
					"version": 0,
					"createdAt": "2022-07-06T10:23:34.931Z",
					"updatedAt": "2022-07-06T10:23:34.931Z",
					"eventId": "9a96ba92-46b4-4e41-bcc0-fb273dbf22b7",
					"firstName": "Sven",
					"lastName": "Frauen",
					"entourageCount": 0,
					"comment": null,
					"email": "sven.frauen@sweap.io",
					"customFields": {},
					"invitationId": "5c59qqypkv6e4s8ue8ay76e9ryqd8742wwgeateacd",
					"invitationState": "ACCEPTED",
					"ticketId": "B44ZZ2XBF69H"
				},
				{
					"id": "cc3a5ae8-49b2-4a8f-bf85-2647d0a0fe25",
					"externalId": null,
					"version": 0,
					"createdAt": "2022-07-06T10:23:34.931Z",
					"updatedAt": "2022-07-06T10:23:34.931Z",
					"eventId": "9a96ba92-46b4-4e41-bcc0-fb273dbf22b7",
					"firstName": "Wissam",
					"lastName": "Ghozlan",
					"entourageCount": 0,
					"comment": null,
					"email": "wissam.ghozlan@sweap.io",
					"customFields": {},
					"invitationId": "ts5am7ufz3689fufc7muuwqjahrzcbnhub5qsnwr40",
					"invitationState": "ACCEPTED",
					"ticketId": "H7VXRCC8HVZK"
				},
				{
					"id": "dd169957-e88c-40fb-95b4-7e823d844de9",
					"externalId": null,
					"version": 0,
					"createdAt": "2022-07-06T10:23:34.931Z",
					"updatedAt": "2022-07-06T10:23:34.931Z",
					"eventId": "9a96ba92-46b4-4e41-bcc0-fb273dbf22b7",
					"firstName": "Matthias",
					"lastName": "Heicke",
					"entourageCount": 0,
					"comment": null,
					"email": "matthias.heicke@sweap.io",
					"customFields": {},
					"invitationId": "maw1qfpttrwfedstiluln7zz9kcrdh9mw4ojehp1kx",
					"invitationState": "ACCEPTED",
					"ticketId": "ZFX99HHXTS5D"
				},
				{
					"id": "c8eb73d7-d481-48eb-bfcd-d2c60aebbec5",
					"externalId": null,
					"version": 0,
					"createdAt": "2022-07-06T10:23:34.931Z",
					"updatedAt": "2022-07-06T10:23:34.931Z",
					"eventId": "9a96ba92-46b4-4e41-bcc0-fb273dbf22b7",
					"firstName": "Sebastian",
					"lastName": "Prestel",
					"entourageCount": 0,
					"comment": null,
					"email": "sebastian.prestel@sweap.io",
					"customFields": {},
					"invitationId": "4q5te5fmm5b0oxa0d7tklehfl4rkl5i3ybabwf7v0q",
					"invitationState": "ACCEPTED",
					"ticketId": "4TX8M2YGKQF2"
				},
				{
					"id": "430de42a-cde1-49e3-9eae-5d8163ed3559",
					"externalId": null,
					"version": 0,
					"createdAt": "2022-07-06T10:23:34.931Z",
					"updatedAt": "2022-07-06T10:23:34.931Z",
					"eventId": "9a96ba92-46b4-4e41-bcc0-fb273dbf22b7",
					"firstName": "Theo",
					"lastName": "Vassiliou",
					"entourageCount": 0,
					"comment": null,
					"email": "theo.vassiliou@sweap.io",
					"customFields": {},
					"invitationId": "ppj3uwr944951408cb1ls0zboogm5uokktogaqzlep",
					"invitationState": "ACCEPTED",
					"ticketId": "TZ9Q46762N6P"
				},
				{
					"id": "c71bc5a3-df81-4e5c-9908-8d657bcdaec4",
					"externalId": null,
					"version": 0,
					"createdAt": "2022-07-06T10:23:34.931Z",
					"updatedAt": "2022-07-06T10:23:34.931Z",
					"eventId": "9a96ba92-46b4-4e41-bcc0-fb273dbf22b7",
					"firstName": "Matija",
					"lastName": "Vojvodic",
					"entourageCount": 0,
					"comment": null,
					"email": "matija.vojvodic@sweap.io",
					"customFields": {
						"default_meta_attribute__title": "Product Design"
					},
					"invitationId": "4tl0rewsta7rdeic0vj2emdlaknaksrm0ktrtlmnng",
					"invitationState": "ACCEPTED",
					"ticketId": "24Z5P7PPSCC9"
				}
			]`)
	rw.Write(response)
}

func TestNewGuestIterator(t *testing.T) {
	guestList := []Guest{
		{
			ID:        "1",
			FirstName: "F1",
			LastName:  "L1",
		},
		{
			ID:        "2",
			FirstName: "F2",
			LastName:  "L2",
		},
		{

			ID:        "3",
			FirstName: "F3",
			LastName:  "L3",
		},
	}
	nextGuest := NewGuestIterator(guestList)
	assert.NotNil(t, nextGuest, "expected not nil guestIterator")
	assert.Equal(t, guestList[0], *nextGuest())
	assert.Equal(t, guestList[1], *nextGuest())
	assert.Equal(t, guestList[2], *nextGuest())
	ng := nextGuest()
	assert.Nil(t, ng, "expected nil")
}

func TestEmptytNewGuestIterator(t *testing.T) {
	guestList := []Guest{}
	nextGuest := NewGuestIterator(guestList)
	assert.NotNil(t, nextGuest, "expected not nil guestIterator")
	ng := nextGuest()
	assert.Nil(t, ng, "expected not nil guest")
}
