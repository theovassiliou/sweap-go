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
	"fmt"
	"net/url"
	"time"
)

type Events []Event

type Event struct {
	// Common fields
	ID         string      `json:"id"`
	Version    int         `json:"version"`
	UpdatedAt  *time.Time  `json:"updatedAt"`
	CreatedAt  *time.Time  `json:"createdAt"`
	ExternalID interface{} `json:"externalId"`
	// Specific fields
	Name                   string                   `json:"name"`
	StartDate              time.Time                `json:"startDate"`
	EndDate                time.Time                `json:"endDate"`
	ZoneId                 string                   `json:"zoneId"`
	AttendanceMode         AttendanceMode           `json:"attendanceMode"`
	State                  string                   `json:"state"`
	CustomFieldDefinitions []CustomFieldDefinitions `json:"customFieldDefinitions"`
}

type AttendanceMode string

const (
	OFFLINE AttendanceMode = "OFFLINE"
	ONLINE  AttendanceMode = "ONLINE"
	MIXED   AttendanceMode = "MIXED"
)

type CustomFieldDefinitions struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	SortIndex int         `json:"sortIndex"`
	GroupName interface{} `json:"groupName"`
	Options   interface{} `json:"options"`
}

type EventPages struct {
	Status   string `json:"status"`
	Content  Events `json:"content"`
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
	            "name": "IntegrationTestEvent2",
	...
	        }
	    ],
	    "pageable": {
	        "size": 1,
	        "totalElements": 3,
	        "totalPages": 3,
	        "page": 1
	    }
	}
	*/
}

// GetEvents will retrieve the complete list of events
// GET /events?[PARAMS]
func (api *Client) GetEvents() (*Events, error) {
	return api.GetEventsContext(context.Background(), NewEventSearchParameters())
}

// GetEventsContext will retrieve the complete list of events with a custom context
func (api *Client) GetEventsContext(ctx context.Context, params EventSearchParameter) (*Events, error) {
	values := url.Values{}

	if params.Id != "" {
		values.Add("id", params.Id)
	}

	if params.Name != "" {
		values.Add("name", params.Name)
	}

	if params.NameContains != "" {
		values.Add("nameContains", params.NameContains)
	}

	if params.State != "" {
		values.Add("state", string(params.State))
	}

	if params.StartDateAfter != nil {
		values.Add("startDateAfter", params.StartDateAfter.Format(time.RFC3339))
	}

	if params.EndDateAfter != nil {
		values.Add("endDateAfter", params.EndDateAfter.Format(time.RFC3339))
	}

	if params.CreatedAfter != nil {
		values.Add("createdAfter", params.CreatedAfter.Format(time.RFC3339))
	}

	if params.UpdatedAfter != nil {
		values.Add("updatedAfter", params.UpdatedAfter.Format(time.RFC3339))
	}

	if params.ExternalID != "" {
		values.Add("externalId", string(params.ExternalID))
	}

	response, err := api.eventsRequest(ctx, "events", values)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// GetEvents will retrieve the event with the given ID
// GET /events/{ID}
func (api *Client) GetEventById(id string) (*Event, error) {
	return api.GetEventByIdContext(context.Background(), id)
}

// CheckCredentials will check wether valid credentials have been used
// GET /events/{ID}
func (api *Client) CheckCredentials() (bool, error) {
	return api.CheckCredentialsContext(context.Background())
}

// GetEventsContext will retrieve the complete list of events with a custom context
func (api *Client) CheckCredentialsContext(ctx context.Context) (bool, error) {
	values := url.Values{}

	_, err := api.checkCredentialsRequest(ctx, "management/check-credentials", values)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return true, nil
}

func (api *Client) checkCredentialsRequest(ctx context.Context, path string, values url.Values) (bool, error) {
	err := api.getMethod(ctx, path, values, nil)
	if err != nil {
		return false, err
	}

	return true, nil
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

func (e Event) String() string {
	return e.Name + ":" + e.ID
}
