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
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type Guests []Guest

type Guest struct {
	// Common fields
	ID         string      `json:"id"`
	Version    int         `json:"version"`
	UpdatedAt  *time.Time  `json:"updatedAt,omitempty"`
	CreatedAt  *time.Time  `json:"createdAt,omitempty"`
	ExternalID interface{} `json:"externalId"`
	// Specific fields
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
	ParentGuestID   string          `json:"parentGuestId"`
	CategoryID      string          `json:"categoryId"`
	AttendanceState AttendanceState `json:"attendanceState"`
}

type GuestPages struct {
	Status   string `json:"status"`
	Content  Guests `json:"content"`
	Error    Error  `json:"error"`
	Pageable struct {
		Size          int `json:"size"`
		TotalElements int `json:"totalElements"`
		TotalPages    int `json:"totalPages"`
		Page          int `json:"page"`
	} `json:"pageable"`

	/*!SECTION
	{
	    "status": "OK",
	    "content": [
	        {
	            "eventId": "8baf13a5-d3c6-419b-848f-b049a9504d2e",
	            "firstName": "Gabriel",
				 ...
	       },
	        {
	            "eventId": "8baf13a5-d3c6-419b-848f-b049a9504d2e",
	            "firstName": "Elena",
				...
			}
	    ],
	    "pageable": {
	        "size": 2,
	        "totalElements": 133,
	        "totalPages": 67,
	        "page": 0
	    }
	}	*/
}
type AttendanceState string

const (
	NONEATTENDANCE AttendanceState = "NONE"
	PRESENT        AttendanceState = "PRESENT"
	GONE           AttendanceState = "GONE"
)

// GuestUpdate carries all update events for guest, let it be new guests, or updated guests
type GuestUpdate struct {
	EventType string
	Value     interface{}
}

const (
	NEWGUEST    = "NEW_GUEST"
	UPDATEGUEST = "UPDATE_GUEST"
)

type CustomFields struct {
	DefaultMetaAttributeTitle string `json:"default_meta_attribute__title"`
}

// CreateGuest creates a Guest in an event and returns a Guest containing it's unigue ID
func (api *Client) CreateGuest(g Guest) (*Guest, error) {
	return api.CreateGuestContext(context.Background(), g)
}

// CreateGuestContext creates the provided Guest (with a given eventID) and returns a Guest containing it's unigue ID with custom context
// It takes a context.Context object.
func (api *Client) CreateGuestContext(ctx context.Context, g Guest) (*Guest, error) {

	if g.EventID == "" {
		return nil, fmt.Errorf("no eventId provided in Guest %v", g)
	}

	request, _ := json.Marshal(g)

	resp := new(Guest)
	err := postJSON(ctx, api.httpclient, api.endpoint+"guests", request, &resp, api)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetGuests will retrieve the complete list of guests for a given event
// You can only get guests for a specific event, so you need to pass the event id via query parameter eventId.
// The guests are not further filtered.
// For more fine granular control use #GetGuestsContext
func (api *Client) GetGuests(eventId string) (*Guests, error) {
	return api.GetGuestsContext(context.Background(), eventId, NewGuestSearchParameters())
}

// GetGuestsContext retrieves the complete list of guests for a given event with a custom context.
// It takes a context.Context object, an event ID, and GuestSearchParameter as input.
// If successful, it returns a pointer to Guests and nil error.
// If either the event ID or GuestSearchParameter is not provided, it returns nil and an error.
// The function can be called on a Client object.
func (api *Client) GetGuestsContext(ctx context.Context, eventId string, params GuestSearchParameter) (*Guests, error) {
	if eventId == "" && (params.Id == "" && params.InvitationId == "") {
		return nil, errors.New("neither eventId nor Guest or InvitationId provided")
	}

	values := url.Values{}
	if eventId != "" {
		values.Set("eventId", eventId)
	}
	if params.Id != "" {
		values.Set("id", params.Id)
	}
	if params.InvitationId != "" {
		values.Set("invitationId", params.InvitationId)
	}

	// Define a helper function to set non-empty parameters
	setIfNotEmpty := func(key, value string) {
		if value != "" {
			values.Set(key, value)
		}
	}

	setIfNotEmpty("firstName", params.FirstName)
	setIfNotEmpty("firstNameContains", params.FirstNameContains)
	setIfNotEmpty("lastName", params.LastName)
	setIfNotEmpty("lastNameContains", params.LastNameContains)
	setIfNotEmpty("email", params.Email)
	setIfNotEmpty("invitationState", string(params.InvitationState))
	setIfNotEmpty("externalId", params.ExternalID)
	setIfNotEmpty("ticketId", params.TicketID)

	if params.CreatedAfter != nil {
		values.Set("createdAfter", params.CreatedAfter.Format(time.RFC3339))
	}
	if params.UpdatedAfter != nil {
		values.Set("updatedAfter", params.UpdatedAfter.Format(time.RFC3339))
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

// GetGuestByIdContext retrieves the guest with the given guest ID using a custom context.
// It takes a context.Context object and a guest ID as input.
// If successful, it returns a pointer to the Guest and nil error.
// If an error occurs during the request, it returns nil and an error.
// The function can be called on a Client object.
func (api *Client) GetGuestByIdContext(ctx context.Context, guestID string) (*Guest, error) {
	// Construct the endpoint URL by formatting the guest ID into the string.
	endpoint := fmt.Sprintf("guests/%s", guestID)

	// Call the guestRequest method on the api object, passing the context,
	// the constructed endpoint, and the empty URL values as parameters.
	response, err := api.guestRequest(ctx, endpoint, url.Values{})
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateGuest updates the guest with the given guest ID.
// It takes a Guest object as input.
// If successful, it returns a pointer to the updated Guest and nil error.
// If an error occurs during the update, it returns nil and the specific error encountered.
// The function can be called on a Client object.
func (api *Client) UpdateGuest(guest Guest) (*Guest, error) {
	// Delegate the update operation to UpdateGuestContext with a background context
	return api.UpdateGuestContext(context.Background(), guest)
}

// UpdateGuestContext updates the guest with the given guest ID using a custom context.
// It takes a context.Context object and a Guest object as input.
// If successful, it returns a pointer to the updated Guest and nil error.
// If the Guest is missing the EventID or ID, it returns a SweapLibraryError with the appropriate error message.
// If an error occurs during the update, it returns nil and the specific error encountered.
// The function can be called on a Client object.
func (api *Client) UpdateGuestContext(ctx context.Context, g Guest) (*Guest, error) {
	// Check if the EventID is missing in the Guest object
	if g.EventID == "" {
		return nil, SweapLibraryError{Message: fmt.Sprintf("no event ID given in Guest %v", g)}
	}

	// Check if the ID is missing in the Guest object
	if g.ID == "" {
		return nil, SweapLibraryError{Message: "no guest ID given"}
	}

	// Marshal the Guest object into JSON
	request, err := json.Marshal(g)
	if err != nil {
		return nil, err
	}

	// Perform the update request and store the response in the updatedGuest variable
	var updatedGuest Guest
	err = putJSON(ctx, api.httpclient, api.endpoint+"guests/"+g.ID, request, &updatedGuest, api)
	if err != nil {
		return nil, err
	}

	return &updatedGuest, nil
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
	response := new(Guests)

	err := api.getMethod(ctx, path, values, response)
	return response, err
}

func (api *Client) guestsPaginatedRequest(ctx context.Context, path string, values url.Values) (GuestPages, error) {
	response := new(GuestPages)

	err := api.getMethod(ctx, path, values, response)
	return *response, err
}

func (api *Client) guestRequest(ctx context.Context, path string, values url.Values) (*Guest, error) {
	response := new(Guest)

	err := api.getMethod(ctx, path, values, response)
	return response, err
}

// NewGuestIterator creates a new iterator function that iterates over a slice of guests.
// It takes a slice of Guest objects as input.
// The returned iterator function can be called repeatedly to retrieve the next Guest from the slice.
// If all guests have been iterated over, the iterator function returns nil.
func NewGuestIterator(guests []Guest) func() *Guest {
	index := 0
	return func() *Guest {
		if index >= len(guests) {
			return nil
		}
		guest := guests[index]
		index++
		return &guest
	}
}

// GetGuestsPaginated will retrieve the complete list of guests for a given event
// You can only get guests for a specific event, so you need to pass the event id via query parameter eventId.
// The guests are not further filtered.
// For more fine granular control use #GetGuestsContext
func (api *Client) GetGuestsPaginated(eventId string, pp PaginationParameter) (GuestPages, error) {
	return api.GetGuestsPaginatedContext(context.Background(), eventId, pp, NewGuestSearchParameters())
}

// GetGuestsPaginatedContext retrieves the paginated list of guests for a given event with a custom context.
// It takes a context.Context object, page and size for pagination, an event ID, and GuestSearchParameter as input.
// If successful, it returns a pointer to Guests and nil error.
// If either the event ID or GuestSearchParameter is not provided, it returns nil and an error.
// If size
// The function can be called on a Client object.
func (api *Client) GetGuestsPaginatedContext(ctx context.Context, eventId string, pages PaginationParameter, params GuestSearchParameter) (GuestPages, error) {
	if eventId == "" && (params.Id == "" && params.InvitationId == "") {
		return GuestPages{}, errors.New("neither eventId nor Guest or InvitationId provided")
	}

	values := url.Values{}
	if eventId != "" {
		values.Set("eventId", eventId)
	}
	if params.Id != "" {
		values.Set("id", params.Id)
	}
	if params.InvitationId != "" {
		values.Set("invitationId", params.InvitationId)
	}

	values.Set("page", strconv.Itoa(pages.Page))
	values.Set("size", strconv.Itoa(pages.Size))

	// Define a helper function to set non-empty parameters
	setIfNotEmpty := func(key, value string) {
		if value != "" {
			values.Set(key, value)
		}
	}

	setIfNotEmpty("firstName", params.FirstName)
	setIfNotEmpty("firstNameContains", params.FirstNameContains)
	setIfNotEmpty("lastName", params.LastName)
	setIfNotEmpty("lastNameContains", params.LastNameContains)
	setIfNotEmpty("email", params.Email)
	setIfNotEmpty("invitationState", string(params.InvitationState))
	setIfNotEmpty("externalId", params.ExternalID)

	if params.CreatedAfter != nil {
		values.Set("createdAfter", params.CreatedAfter.Format(time.RFC3339))
	}
	if params.UpdatedAfter != nil {
		values.Set("updatedAfter", params.UpdatedAfter.Format(time.RFC3339))
	}

	response, err := api.guestsPaginatedRequest(ctx, "guests/paginated", values)
	if err != nil {
		return GuestPages{}, err
	}

	return response, nil
}
