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
	"net/url"
	"time"
)

type GuestBulkImportStatus string

type GuestBulkImportState struct {
	ID    string                `json:"id"`
	State GuestBulkImportStatus `json:"state"`
}

const (
	UPLOADSTARTED  GuestBulkImportStatus = "UPLOAD_STARTED"
	UPLOADFINISHED GuestBulkImportStatus = "UPLOAD_FINISHED"
	IMPORTSTARTED  GuestBulkImportStatus = "IMPORT_STARTED"
	IMPORTFINISHED GuestBulkImportStatus = "IMPORT_FINISHED"
)

type GuestBulkImports []GuestBulkImport

type GuestBulkImport struct {
	// Common
	ID         string      `json:"id"`
	Version    int         `json:"version"`
	UpdatedAt  *time.Time  `json:"updatedAt"`
	CreatedAt  *time.Time  `json:"createdAt"`
	ExternalID interface{} `json:"externalId"`

	// Specific
	Name                   string                   `json:"name"`
	State                  *GuestBulkImportStatus   `json:"state"`
	Guests                 *Guests                  `json:"guests"`
	EventId                string                   `json:"eventId"`
	CustomFieldDefinitions []CustomFieldDefinitions `json:"customFieldDefinitions"`
}

// GetAllBulkImports lists all bulk import objects started or completed.
// GET /guest-bulk-imports?[PARAMS]
func (api *Client) GetAllBulkImports(s ...GuestBulkImportSearchParameter) (*GuestBulkImports, error) {
	if len(s) > 1 {
		return api.GetAllBulkImportsContext(context.Background(), s[0])
	}
	return api.GetAllBulkImportsContext(context.Background(), NewGuestBulkImportSearchParameter())
}

func (api *Client) GetAllBulkImportsContext(ctx context.Context, params GuestBulkImportSearchParameter) (*GuestBulkImports, error) {
	values := url.Values{}

	if params.GuestBulkImportId != "" {
		values.Add("id", params.GuestBulkImportId)
	}

	if params.EventId != "" {
		values.Add("eventId", params.EventId)
	}

	if params.State != "" {
		values.Add("state", string(params.State))
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

	response, err := api.gbisRequest(ctx, "guest-bulk-imports", values)
	if err != nil {
		return nil, err
	}
	return response, nil

}

func (api *Client) GetSpecificBulkImport(gbiId string) (*GuestBulkImport, error) {

	return api.GetSpecificBulkImportContext(context.Background(), gbiId)
}

func (api *Client) GetSpecificBulkImportContext(ctx context.Context, gbiId string) (*GuestBulkImport, error) {

	if gbiId == "" {
		return nil, SweapLibraryError{"no guest bulk import ID given"}
	}

	response, err := api.gbiRequest(ctx, "guest-bulk-imports/"+gbiId, nil)
	if err != nil {
		return nil, err
	}
	return response, nil

}

func (api *Client) GetSpecificBulkImportState(gbiId string) (*GuestBulkImportState, error) {

	return api.GetSpecificBulkImportStateContext(context.Background(), gbiId)
}

func (api *Client) GetSpecificBulkImportStateContext(ctx context.Context, gbiId string) (*GuestBulkImportState, error) {

	if gbiId == "" {
		return nil, SweapLibraryError{"no guest bulk import ID given"}
	}

	response, err := api.gbiStateRequest(ctx, "guest-bulk-imports/"+gbiId+"/state", nil)
	if err != nil {
		return nil, err
	}
	return response, nil

}

// Create Bulk Import Object creates a guest bulk object.
// POST /guest-bulk-imports
func (api *Client) CreateGuestBulkImportObject(gbi GuestBulkImport) (*GuestBulkImport, error) {
	return api.CreateGuestBulkImportObjectContext(context.Background(), gbi)
}

// CreateGuestBulkImportObjectContext creates a GBI in an event and returns a GBI containing it's unigue ID with custom context
func (api *Client) CreateGuestBulkImportObjectContext(ctx context.Context, gbi GuestBulkImport) (*GuestBulkImport, error) {

	request, _ := json.Marshal(gbi)

	resp := &GuestBulkImport{}
	err := postJSON(ctx, api.httpclient, api.endpoint+"guest-bulk-imports", request, &resp, api)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DeleteGuestBulkImportObject will delete the gbi with the given Id. If guest is not found error will be returned
func (api *Client) DeleteGuestBulkImportObject(gbiId string) error {
	return api.DeleteGuestBulkImportObjectContext(context.Background(), gbiId)
}

// DeleteGuest will delete the guest with the given guestId wiht a custom context.
// If guest is not found error will be returned
func (api *Client) DeleteGuestBulkImportObjectContext(ctx context.Context, gbiId string) error {
	if gbiId == "" {
		return SweapLibraryError{"no guest bulk import ID given"}
	}

	return deleteResource(ctx, api.httpclient, api.endpoint+"guest-bulk-imports/"+gbiId, nil, api)

}

// FIXME: Need documentation
func (api *Client) BulkImportUpdateBatch(gibID string, guests Guests) error {
	return api.BulkImportUpdateBatchContext(context.Background(), gibID, guests)
}

// FIXME: BulkImportUpdateBatchContext will update the guest with the given guestId with a custom context
func (api *Client) BulkImportUpdateBatchContext(ctx context.Context, gibID string, guests Guests) error {

	if gibID == "" {
		return SweapLibraryError{"no guest bulk id given"}
	}

	request, _ := json.Marshal(guests)

	err := putJSON(ctx, api.httpclient, api.endpoint+"guest-bulk-imports/"+gibID+"/upload-batch", request, nil, api)

	if err != nil {
		return err
	}

	return nil
}

// Create Bulk Import Object creates a guest bulk object.
// POST /guest-bulk-imports
func (api *Client) BulkImportInOneGo(gbi GuestBulkImport) (*GuestBulkImport, error) {
	return api.BulkImportInOneGoContext(context.Background(), gbi)
}

// CreateGuestBulkImportObjectContext creates a GBI in an event and returns a GBI containing it's unigue ID with custom context
func (api *Client) BulkImportInOneGoContext(ctx context.Context, gbi GuestBulkImport) (*GuestBulkImport, error) {

	request, _ := json.Marshal(gbi)

	resp := &GuestBulkImport{}
	err := postJSON(ctx, api.httpclient, api.endpoint+"guest-bulk-imports", request, &resp, api)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (api *Client) BulkImportFinishUpload(gibID string) error {
	return api.BulkImportStateContext(context.Background(), gibID)
}

func (api *Client) BulkImportStateContext(ctx context.Context, gibID string) error {
	if gibID == "" {
		return SweapLibraryError{"no guest bulk id given"}
	}

	state := GuestBulkImportState{
		State: UPLOADFINISHED,
	}

	request, _ := json.Marshal(state)

	err := putJSON(ctx, api.httpclient, api.endpoint+"guest-bulk-imports/"+gibID+"/state", request, nil, api)
	return err

}

func (api *Client) gbisRequest(ctx context.Context, path string, values url.Values) (*GuestBulkImports, error) {
	response := &GuestBulkImports{}
	err := api.getMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *Client) gbiRequest(ctx context.Context, path string, values url.Values) (*GuestBulkImport, error) {
	response := &GuestBulkImport{}
	err := api.getMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *Client) gbiStateRequest(ctx context.Context, path string, values url.Values) (*GuestBulkImportState, error) {
	response := &GuestBulkImportState{}
	err := api.getMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
