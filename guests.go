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
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type Guests []Guest

type Guest struct {
	ID              string          `json:"id"`
	ExternalID      interface{}     `json:"externalId"`
	Version         int             `json:"version"`
	CreatedAt       *time.Time      `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time      `json:"updatedAt,omitempty"`
	EventID         string          `json:"eventId"`
	FirstName       string          `json:"firstName"`
	LastName        string          `json:"lastName"`
	EntourageCount  int             `json:"entourageCount,omitempty"`
	Comment         interface{}     `json:"comment,omitempty"`
	Email           string          `json:"email,omitempty"`
	CustomFields    CustomFields    `json:"customFields,omitempty"`
	InvitationID    string          `json:"invitationId,omitempty"`
	InvitationState InvitationState `json:"invitationState"`
	TicketID        string          `json:"ticketId,omitempty"`
}

type CustomFields struct {
	DefaultMetaAttributeTitle string `json:"default_meta_attribute__title"`
}

// CreateGuest creates a Guest in an event and returns a Guest containing it's unigue ID
func (api *Client) CreateGuest(g Guest) (*Guest, error) {
	return api.CreateGuestConext(context.Background(), g)
}

// CreateGuest creates a Guest in an event and returns a Guest containing it's unigue ID with custom context
func (api *Client) CreateGuestConext(ctx context.Context, g Guest) (*Guest, error) {

	if g.EventID == "" {
		panic(fmt.Sprintf("no event id given in Guest %v", g))
	}

	request, _ := json.Marshal(g)

	resp := &Guest{}
	err := postJSON(ctx, api.httpclient, api.endpoint+"guests", request, &resp, api)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetGuests will retrieve the complete list of guests for a given event
func (api *Client) GetGuests(eventId string) (*Guests, error) {
	return api.GetGuestsContext(context.Background(), eventId, NewGuestSearchParameters())
}

// GetGuestsContext will retrieve the complete list of guests for a given event with a custom context
func (api *Client) GetGuestsContext(ctx context.Context, eventId string, params GuestSearchParameter) (*Guests, error) {
	values := url.Values{}

	if params.FirstName != "" {
		values.Add("first_name", params.FirstName)
	}

	if params.LastName != "" {
		values.Add("last_name", params.LastName)
	}

	if params.Email != "" {
		values.Add("email", params.Email)
	}

	if params.InvitationState != "" {
		values.Add("invitation_state", string(params.InvitationState))
	}

	if params.ExternalID != "" {
		values.Add("external_id", params.ExternalID)
	}
	if params.CreatedAfter != nil {
		values.Add("start_date_after", params.CreatedAfter.Format(time.RFC3339))
	}

	if params.UpdatedAfter != nil {
		values.Add("end_date_after", params.UpdatedAfter.Format(time.RFC3339))
	}

	if eventId != "" {
		values.Add("event_id", eventId)
	} else {
		panic("no event id given")
	}
	response, err := api.guestsRequest(ctx, "guests", values)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// GetGuestById will retrieve the guest with the given guestId
func (api *Client) GetGuestById(guestId string) (*Guest, error) {
	return api.GetGuestByIdContext(context.Background(), guestId)
}

// GetGuestByIdContext will retrieve the guest with the given guestId with a custom context
func (api *Client) GetGuestByIdContext(ctx context.Context, guestId string) (*Guest, error) {
	values := url.Values{}

	response, err := api.guestRequest(ctx, "guests/"+guestId, values)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateGuest will update the given guest with the given guestId
func (api *Client) UpdateGuest(guest Guest) (*Guest, error) {
	return api.UpdateGuestContext(context.Background(), guest)
}

// UpdateGuestContext will update the guest with the given guestId with a custom context
func (api *Client) UpdateGuestContext(ctx context.Context, g Guest) (*Guest, error) {

	if g.EventID == "" {
		return nil, SweapLibraryError{fmt.Sprintf("no event id given in Guest %v", g)}
	}
	if g.ID == "" {
		return nil, SweapLibraryError{"no guest ID given"}
	}

	request, _ := json.Marshal(g)

	resp := &Guest{}
	err := putJSON(ctx, api.httpclient, api.endpoint+"guests/"+g.ID, request, &resp, api)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DeleteGuest will delete the guest with the given guestId. If guest is not found error will be returned
func (api *Client) DeleteGuest(guestId string) error {
	return api.DeleteGuestContext(context.Background(), guestId)
}

// DeleteGuest will delete the guest with the given guestId wiht a custom context.
// If guest is not found error will be returned
func (api *Client) DeleteGuestContext(ctx context.Context, guestId string) error {
	if guestId == "" {
		return SweapLibraryError{"no guest ID given"}
	}

	return deleteResource(ctx, api.httpclient, api.endpoint+"guests/"+guestId, nil, api)

}

func (api *Client) guestsRequest(ctx context.Context, path string, values url.Values) (*Guests, error) {
	response := &Guests{}

	err := api.getMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *Client) guestRequest(ctx context.Context, path string, values url.Values) (*Guest, error) {
	response := &Guest{}

	err := api.getMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
