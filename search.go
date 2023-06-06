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
	"time"
)

type EventState string

const (
	DRAFT  EventState = "DRAFT"
	ACTIVE EventState = "ACTIVE"
	CLOSED EventState = "CLOSED"
)

type EventSearchParameter struct {
	Id             string     // optional, filters equal event id
	Name           string     // optional, filters equal case insensitive event name
	NameContains   string     // optional, filters equal case insensitive event name
	State          EventState // optional, one of DRAFT, ACTIVE, CLOSED
	StartDateAfter *time.Time // optional, filters greater than event startDate: ISO 8601 timestamp
	EndDateAfter   *time.Time // optional, filters greater than event startDate: ISO 8601 timestamp
	ExternalID     string     // optional, filters equal event externalId
	CreatedAfter   *time.Time // optional, filters greater than event createdAt: ISO 8601 timestamp
	UpdatedAfter   *time.Time // optional, filters greater than event updatedAt: ISO 8601 timestamp
}

func NewEventSearchParameters() EventSearchParameter {
	return EventSearchParameter{}
}

type EventStatisticsSearchParameter struct {
	Id           string     // optional, filters equal event id
	ExternalID   string     // optional, filters equal event externalId
	CreatedAfter *time.Time // optional, filters greater than event createdAt: ISO 8601 timestamp
	UpdatedAfter *time.Time // optional, filters greater than event updatedAt: ISO 8601 timestamp

	MinGuestCount    int // optional, filters guestCount equal or greater than minGuestCount value
	MaxGuestCount    int // optional, filters guestCount equal or lesser than maxGuestCount value
	MinAcceptedCount int // optional, filters acceptedCount equal or greater than minAcceptedCount value
	MaxAcceptedCount int // optional, filters acceptedCount equal or lesser than maxAcceptedCount value
	MinDeclinedCount int // optional, filters declinedCount equal or greater than minDeclinedCount value
	MaxDeclinedCount int // optional, filters declinedCount equal or lesser than maxDeclinedCount value
	MinNoReplyCount  int // optional, filters noReplyCount equal or greater than minNoReplyCount value
	MaxNoReplyCount  int // optional, filters noReplyCount equal or lesser than maxNoReplyCount value
	MinCheckinCount  int // optional, filters checkinCount equal or greater than minCheckinCount value
	MaxCheckinCount  int // optional, filters checkinCount equal or lesser than maxCheckinCount value
}

func NewEventStatisticsSearchParameter() EventStatisticsSearchParameter {
	return EventStatisticsSearchParameter{"", "", nil, nil, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}
}

// SearchEvents will retrieve a list of events filtered by EventSearchParameter
func (api *Client) SearchEvents(params EventSearchParameter) (*Events, error) {
	return api.GetEventsContext(context.Background(), params)
}

// SearchEventStatistics will retrieve a list of statistics for events filtered by EventStatisticsSearchParameter
func (api *Client) SearchEventStatistics(params EventStatisticsSearchParameter) (*EventStatistics, error) {
	return api.GetEventStatisticsContext(context.Background(), params)
}

type InvitationState string

const (
	NONE     InvitationState = "NONE"
	NO_REPLY InvitationState = "NO_REPLY"
	ACCEPTED InvitationState = "ACCEPTED"
	DECLINED InvitationState = "DECLINED"
)

type GuestSearchParameter struct {
	Id                string          // optional, filters equal guest id, result will list of 0 or 1 guest, conditionally required if eventId or invitationId are empty
	FirstName         string          // optional, filters equal case insensitive guest firstName
	FirstNameContains string          // optional, filters contains case insensitive guest firstName
	LastName          string          // optional, filters equal case insensitive guest lastName
	LastNameContains  string          // optional, filters contains case insensitive guest lastName
	Email             string          // optional, filters equal case insensitive guest email
	InvitationState   InvitationState // optional, filters equal guest invitationState: NONE, NO_REPLY, ACCEPTED or DECLINED
	ExternalID        string          // optional, filters equal guest externalId
	CreatedAfter      *time.Time      // optional, filters greater than guest createdAt: ISO 8601 timestamp
	UpdatedAfter      *time.Time      // optional, filters greater than guest updatedAt: ISO 8601 timestamp
	InvitationId      string          // optional, filters equal guest invitationId, result will list of 0 or 1 guest, conditionally required if event_id or id are empty
}

func NewGuestSearchParameters() GuestSearchParameter {
	return GuestSearchParameter{}
}

// SearchGuests will retrieve the complete list of Guests matching the search criteria
func (api *Client) SearchGuests(id string, params GuestSearchParameter) (*Guests, error) {
	return api.GetGuestsContext(context.Background(), id, params)
}

type CategorySearchParameter struct {
	GuestId      string     // optional, filters equal guest id
	Name         string     // optional, filters equal case insensitive category name
	ExternalID   string     // optional, filters equal guest externalId
	CreatedAfter *time.Time // optional, filters greater than guest createdAt: ISO 8601 timestamp
	UpdatedAfter *time.Time // optional, filters greater than guest updatedAt: ISO 8601 timestamp
}

func NewCategorySearchParameter() CategorySearchParameter {
	return CategorySearchParameter{}
}

type GuestBulkImportSearchParameter struct {
	GuestBulkImportId string                // optional, filters equal guest-bulk-import id
	EventId           string                // optional, filters equal event id
	State             GuestBulkImportStatus // optional, filters equal guest-bulk-import state: UPLOAD_STARTED, UPLOAD_FINISHED, IMPORT_STARTED, IMPORT_FINISHED
	ExternalID        string                // optional, filters equal guest-bulk-import externalId
	CreatedAfter      *time.Time            // optional, filters greater than guest-bulk-import createdAt: ISO 8601 timestamp
	UpdatedAfter      *time.Time            // optional, filters greater than guest-bulk-import updatedAt: ISO 8601 timestamp
}

func NewGuestBulkImportSearchParameter() GuestBulkImportSearchParameter {
	return GuestBulkImportSearchParameter{}
}
