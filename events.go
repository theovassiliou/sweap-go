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
	"net/url"
	"time"
)

type Events []Event

type Event struct {
	ID                     string                   `json:"id"`
	ExternalID             interface{}              `json:"externalId"`
	Version                int                      `json:"version"`
	CreatedAt              time.Time                `json:"createdAt"`
	UpdatedAt              time.Time                `json:"updatedAt"`
	Name                   string                   `json:"name"`
	StartDate              time.Time                `json:"startDate"`
	EndDate                time.Time                `json:"endDate"`
	State                  string                   `json:"state"`
	CustomFieldDefinitions []CustomFieldDefinitions `json:"customFieldDefinitions"`
}

type CustomFieldDefinitions struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	SortIndex int         `json:"sortIndex"`
	Type      string      `json:"type"`
	GroupName interface{} `json:"groupName"`
	Options   interface{} `json:"options"`
}

// GetEvents will retrieve the complete list of events
func (api *Client) GetEvents() (*Events, error) {
	return api.GetEventsContext(context.Background(), NewEventSearchParameters())
}

// GetEventsContext will retrieve the complete list of events with a custom context
func (api *Client) GetEventsContext(ctx context.Context, params EventSearchParameter) (*Events, error) {
	values := url.Values{}

	if params.Name != "" {
		values.Add("name", params.Name)
	}

	if params.State != "" {
		values.Add("state", string(params.State))
	}

	if params.StartDateAfter != nil {
		values.Add("start_date_after", params.StartDateAfter.Format(time.RFC3339))
	}

	if params.EndDateAfter != nil {
		values.Add("end_date_after", params.EndDateAfter.Format(time.RFC3339))
	}

	if params.CreatedAfter != nil {
		values.Add("created_after", params.CreatedAfter.Format(time.RFC3339))
	}

	if params.UpdatedAfter != nil {
		values.Add("updated_after", params.UpdatedAfter.Format(time.RFC3339))
	}

	if params.ExternalID != "" {
		values.Add("external_id", string(params.ExternalID))
	}

	response, err := api.eventsRequest(ctx, "events", values)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// GetEvents will retrieve the event with the given ID
func (api *Client) GetEventById(id string) (*Event, error) {
	return api.GetEventByIdContext(context.Background(), id)
}

// GetEventsContext will retrieve the complete list of events with a custom context
func (api *Client) GetEventByIdContext(ctx context.Context, id string) (*Event, error) {
	values := url.Values{}

	response, err := api.eventRequest(ctx, "events/"+id, values)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *Client) eventsRequest(ctx context.Context, path string, values url.Values) (*Events, error) {
	response := &Events{}
	err := api.getMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *Client) eventRequest(ctx context.Context, path string, values url.Values) (*Event, error) {
	response := &Event{}
	err := api.getMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
