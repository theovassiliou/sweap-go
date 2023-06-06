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

// GetEventStatisticsContext returns all events statistics under your account.
func (api *Client) GetEventStatisticsContext(ctx context.Context, params EventStatisticsSearchParameter) (*EventStatistics, error) {
	values := url.Values{}

	if params.Id != "" {
		values.Add("id", params.Id)
	}

	if params.ExternalID != "" {
		values.Add("externalId", string(params.ExternalID))
	}

	if params.CreatedAfter != nil {
		values.Add("createdAfter", params.CreatedAfter.Format(time.RFC3339))
	}

	if params.UpdatedAfter != nil {
		values.Add("updatedAfter", params.UpdatedAfter.Format(time.RFC3339))
	}

	if params.MinGuestCount > -1 {
		values.Add("minGuestCount", fmt.Sprint(params.MinGuestCount))
	}

	if params.MaxGuestCount > -1 {
		values.Add("maxGuestCount", fmt.Sprint(params.MaxGuestCount))
	}

	if params.MinAcceptedCount > -1 {
		values.Add("minAcceptedCount", fmt.Sprint(params.MinAcceptedCount))
	}

	if params.MaxAcceptedCount > -1 {
		values.Add("maxAcceptedCount", fmt.Sprint(params.MaxAcceptedCount))
	}

	if params.MinDeclinedCount > -1 {
		values.Add("minDeclinedCount", fmt.Sprint(params.MinDeclinedCount))
	}

	if params.MaxDeclinedCount > -1 {
		values.Add("maxDeclinedCount", fmt.Sprint(params.MaxDeclinedCount))
	}

	if params.MaxAcceptedCount > -1 {
		values.Add("maxAcceptedCount", fmt.Sprint(params.MaxAcceptedCount))
	}

	if params.MinNoReplyCount > -1 {
		values.Add("maxAcceptedCount", fmt.Sprint(params.MinNoReplyCount))
	}

	if params.MaxNoReplyCount > -1 {
		values.Add("maxNoReplyCount", fmt.Sprint(params.MaxNoReplyCount))
	}

	if params.MinCheckinCount > -1 {
		values.Add("minCheckinCount", fmt.Sprint(params.MinCheckinCount))
	}

	if params.MaxCheckinCount > -1 {
		values.Add("maxCheckinCount", fmt.Sprint(params.MaxCheckinCount))
	}

	response, err := api.eventStatisticsRequest(ctx, "event-statistics", values)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (api *Client) eventStatisticsRequest(ctx context.Context, path string, values url.Values) (*EventStatistics, error) {
	response := &EventStatistics{}
	err := api.getMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetEventStatisticsById will retrieve the statistic for an event with a given ID
func (api *Client) GetEventStatisticsById(id string) (*EventStatistic, error) {
	return api.GetEventStatisticsByIdContext(context.Background(), id)
}

// GetEventStatisticsByIdContext will retrieve the statistic for an event with a given ID with a given context
func (api *Client) GetEventStatisticsByIdContext(ctx context.Context, id string) (*EventStatistic, error) {
	values := url.Values{}

	response, err := api.eventStatisticRequest(ctx, "events/"+id, values)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *Client) eventStatisticRequest(ctx context.Context, path string, values url.Values) (*EventStatistic, error) {
	response := &EventStatistic{}
	err := api.getMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
