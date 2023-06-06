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

type Categories []Category

type Category struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ColorHex  string `json:"colorHex"`
	SortIndex int    `json:"sortIndex"`
	EventID   string `json:"eventId"`
}

// GetCategories will retrieve the complete list of guest categories for a given eventId
func (api *Client) GetCategories(eventId string) (*Categories, error) {
	return api.GetCategoriesContext(context.Background(), eventId, NewCategorySearchParameter())
}

// GetCategories will retrieve the complete list of guest categories for a given eventId with a custom context
func (api *Client) GetCategoriesContext(ctx context.Context, eventId string, params CategorySearchParameter) (*Categories, error) {
	values := url.Values{}

	if params.GuestId != "" {
		values.Add("id", params.GuestId)
	}
	if params.Name != "" {
		values.Add("id", params.Name)
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
	response, err := api.categoriesRequest(ctx, "categories", values)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// GetCategoryById will retrieve the guest with the given guestId
func (api *Client) GetCategoryById(categoryId string) (*Category, error) {
	return api.GetCategoryByIdContext(context.Background(), categoryId)
}

// GetGuestByIdContext will retrieve the guest with the given guestId with a custom context
func (api *Client) GetCategoryByIdContext(ctx context.Context, categoryId string) (*Category, error) {
	values := url.Values{}

	response, err := api.categoryRequest(ctx, "categories/"+categoryId, values)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *Client) categoriesRequest(ctx context.Context, path string, values url.Values) (*Categories, error) {
	response := &Categories{}

	err := api.getMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *Client) categoryRequest(ctx context.Context, path string, values url.Values) (*Category, error) {
	response := &Category{}

	err := api.getMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
