/*
 Copyright (c) 2023 Theofanis Vassiliou-Gioles

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
	"net/url"
	"time"
)

type EventStatistics []EventStatistic

type EventStatistic struct {
	// Common fields
	ID         string      `json:"id"`
	Version    int         `json:"version"`
	UpdatedAt  time.Time   `json:"updatedAt"`
	CreatedAt  time.Time   `json:"createdAt"`
	ExternalID interface{} `json:"externalId"`
	// Specific fields
	GuestCount     int `json:"guestCount"`    // Guests overall count (incl. companions)
	AcceptedCount  int `json:"acceptedCount"` // Accepted guest count (incl. companions)
	DeclindedCount int `json:"declinedCount"` // Declined guest count (incl. companions)
	NoReplyCount   int `json:"noReplyCount"`  // No Reply guest count (incl. companions)
	CheckinCount   int `json:"checkinCount"`  // Check-in count (incl. companions)

}

// GetEventStatistics returns all events statistics under your account.
func (api *Client) GetEventStatistics() (*EventStatistics, error) {
	return api.GetEventStatisticsContext(context.Background(), NewEventStatisticsSearchParameter())
}

// GetEventStatisticsContext returns all event statistics under your account based on the provided search parameters.
// It takes a context.Context object and an EventStatisticsSearchParameter object as input.
// If successful, it returns a pointer to EventStatistics and nil error.
// If an error occurs during the request, it returns nil and the specific error encountered.
// The function can be called on a Client object.
func (api *Client) GetEventStatisticsContext(ctx context.Context, params EventStatisticsSearchParameter) (*EventStatistics, error) {
	// Initialize URL values to be used for query parameters.
	values := url.Values{}

	// Helper function to set non-negative integer parameters
	setNonNegativeParam := func(key string, value int) {
		if value >= 0 {
			values.Add(key, fmt.Sprint(value))
		}
	}

	// Set query parameters based on the provided search parameters
	if params.Id != "" {
		values.Set("id", params.Id)
	}
	if params.ExternalID != "" {
		values.Set("externalId", string(params.ExternalID))
	}
	if params.CreatedAfter != nil {
		values.Set("createdAfter", params.CreatedAfter.Format(time.RFC3339))
	}
	if params.UpdatedAfter != nil {
		values.Set("updatedAfter", params.UpdatedAfter.Format(time.RFC3339))
	}

	setNonNegativeParam("minGuestCount", params.MinGuestCount)
	setNonNegativeParam("maxGuestCount", params.MaxGuestCount)
	setNonNegativeParam("minAcceptedCount", params.MinAcceptedCount)
	setNonNegativeParam("maxAcceptedCount", params.MaxAcceptedCount)
	setNonNegativeParam("minDeclinedCount", params.MinDeclinedCount)
	setNonNegativeParam("maxDeclinedCount", params.MaxDeclinedCount)
	setNonNegativeParam("minNoReplyCount", params.MinNoReplyCount)
	setNonNegativeParam("maxNoReplyCount", params.MaxNoReplyCount)
	setNonNegativeParam("minCheckinCount", params.MinCheckinCount)
	setNonNegativeParam("maxCheckinCount", params.MaxCheckinCount)

	response, err := api.eventStatisticsRequest(ctx, "event-statistics", values)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (api *Client) eventStatisticsRequest(ctx context.Context, path string, values url.Values) (*EventStatistics, error) {
	response := new(EventStatistics)
	err := api.getMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetEventStatisticsByID will retrieve the statistic for an event with a given ID
func (api *Client) GetEventStatisticsByID(id string) (*EventStatistic, error) {
	return api.GetEventStatisticsByIDContext(context.Background(), id)
}

// GetEventStatisticsByIDContext retrieves the statistics for an event with the given ID using a custom context.
// It takes a context.Context object and an event ID as input.
// If successful, it returns a pointer to the EventStatistic and nil error.
// If an error occurs during the request, it returns nil and the specific error encountered.
// The function can be called on a Client object.
func (api *Client) GetEventStatisticsByIDContext(ctx context.Context, eventID string) (*EventStatistic, error) {
	// Construct the endpoint URL by formatting the eventID into the string.
	endpoint := fmt.Sprintf("events/%s", eventID)
	response, err := api.eventStatisticRequest(ctx, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *Client) eventStatisticRequest(ctx context.Context, path string, values url.Values) (*EventStatistic, error) {
	response := new(EventStatistic)
	err := api.getMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
