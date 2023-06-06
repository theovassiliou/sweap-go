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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEvents(t *testing.T) {
	events, err := api.GetEvents()

	assert.Equal(t, 5, len(*events))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestCGetEvent(t *testing.T) {
	event, err := api.GetEventById("9a96ba92-46b4-4e41-bcc0-fb273dbf22b7")

	assert.NotNil(t, event)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	assert.Equal(t, "9a96ba92-46b4-4e41-bcc0-fb273dbf22b7", event.ID)
}

func TestGetEventCantFind(t *testing.T) {
	event, err := api.GetEventById("CANTFIND")

	assert.Nil(t, event)
	if err == nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestSearchName(t *testing.T) {
	search := EventSearchParameter{
		Name: "retro ownership",
	}

	events, err := api.SearchEvents(search)

	assert.NotNil(t, events)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	assert.Equal(t, 1, len(*events))
}

// ------ Aux functions ------
func getEvent(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response := []byte(`
		{
			"id": "9a96ba92-46b4-4e41-bcc0-fb273dbf22b7",
			"externalId": null,
			"version": 2,
			"createdAt": "2022-07-06T10:23:34.871Z",
			"updatedAt": "2022-07-06T10:27:37.130Z",
			"name": "Retro Ownership",
			"startDate": "2022-07-26T12:00:00.000Z",
			"endDate": "2022-07-26T14:00:00.000Z",
			"state": "DRAFT",
			"customFieldDefinitions": [
			{
			"id": "default_meta_attribute__title",
			"name": "Title",
			"sortIndex": 1,
			"type": "TEXT",
			"groupName": null,
			"options": null
			},
			{
			"id": "default_meta_attribute__salutation",
			"name": "Custom Salutation",
			"sorrtIndex": 3,
			"type": "TEXT",
			"groupName": null,
			"options": null
			}
			]
			}`)
	rw.Write(response)

}

func getEvents(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	val := r.URL.Query()

	if val.Get("name") == "retro ownership" {
		response := []byte(`[
		{
			"id": "9a96ba92-46b4-4e41-bcc0-fb273dbf22b7",
			"externalId": null,
			"version": 2,
			"createdAt": "2022-07-06T10:23:34.871Z",
			"updatedAt": "2022-07-06T10:27:37.130Z",
			"name": "Retro Ownership",
			"startDate": "2022-07-26T12:00:00.000Z",
			"endDate": "2022-07-26T14:00:00.000Z",
			"state": "DRAFT",
			"customFieldDefinitions": [
			{
			"id": "default_meta_attribute__title",
			"name": "Title",
			"sortIndex": 1,
			"type": "TEXT",
			"groupName": null,
			"options": null
			},
			{
			"id": "default_meta_attribute__salutation",
			"name": "Custom Salutation",
			"sorrtIndex": 3,
			"type": "TEXT",
			"groupName": null,
			"options": null
			}
			]
			}]`)
		rw.Write(response)
	} else {

		response := []byte(`[
		{
			"id": "9a96ba92-46b4-4e41-bcc0-fb273dbf22b7",
			"externalId": null,
			"version": 2,
			"createdAt": "2022-07-06T10:23:34.871Z",
			"updatedAt": "2022-07-06T10:27:37.130Z",
			"name": "Retro Ownership",
			"startDate": "2022-07-26T12:00:00.000Z",
			"endDate": "2022-07-26T14:00:00.000Z",
			"state": "DRAFT",
			"customFieldDefinitions": [
				{
					"id": "default_meta_attribute__title",
					"name": "Title",
					"sortIndex": 1,
					"type": "TEXT",
					"groupName": null,
					"options": null
				},
				{
					"id": "default_meta_attribute__salutation",
					"name": "Custom Salutation",
					"sortIndex": 3,
					"type": "TEXT",
					"groupName": null,
					"options": null
				}
			]
		},
		{
			"id": "c00a3c96-c35c-44f7-b5fd-d8526e8fdd9b",
			"externalId": null,
			"version": 2,
			"createdAt": "2022-06-22T12:24:52.274Z",
			"updatedAt": "2022-06-28T07:17:19.908Z",
			"name": "Retro / Workshop Ownership",
			"startDate": "2022-07-05T12:00:00.000Z",
			"endDate": "2022-07-05T14:00:00.000Z",
			"state": "ACTIVE",
			"customFieldDefinitions": [
				{
					"id": "default_meta_attribute__title",
					"name": "Title",
					"sortIndex": 1,
					"type": "TEXT",
					"groupName": null,
					"options": null
				},
				{
					"id": "default_meta_attribute__salutation",
					"name": "Custom Salutation",
					"sortIndex": 3,
					"type": "TEXT",
					"groupName": null,
					"options": null
				}
			]
		},
		{
			"id": "d6785054-d9b1-433c-8810-31cd11d6d936",
			"externalId": null,
			"version": 0,
			"createdAt": "2022-06-03T07:42:45.331Z",
			"updatedAt": "2022-06-03T07:42:45.331Z",
			"name": "Summer Reception 2022",
			"startDate": "2022-07-01T10:00:00.000Z",
			"endDate": "2022-07-01T14:00:00.000Z",
			"state": "DRAFT",
			"customFieldDefinitions": [
				{
					"id": "default_meta_attribute__title",
					"name": "Title",
					"sortIndex": 1,
					"type": "TEXT",
					"groupName": null,
					"options": null
				},
				{
					"id": "default_meta_attribute__salutation",
					"name": "Custom Salutation",
					"sortIndex": 3,
					"type": "TEXT",
					"groupName": null,
					"options": null
				}
			]
		},
		{
			"id": "ad2ef848-de80-4974-b69e-bfb50790a8cb",
			"externalId": null,
			"version": 0,
			"createdAt": "2022-03-14T10:51:41.430Z",
			"updatedAt": "2022-03-14T10:51:41.430Z",
			"name": "Demo Event",
			"startDate": "2022-05-13T14:00:00.000Z",
			"endDate": "2022-05-13T18:00:00.000Z",
			"state": "ACTIVE",
			"customFieldDefinitions": [
				{
					"id": "default_meta_attribute__salutation",
					"name": "Eigene Briefanrede",
					"sortIndex": 3,
					"type": "TEXT",
					"groupName": null,
					"options": null
				},
				{
					"id": "default_meta_attribute__company",
					"name": "Firma",
					"sortIndex": 4,
					"type": "TEXT",
					"groupName": "Firmendetails",
					"options": null
				},
				{
					"id": "default_meta_attribute__title",
					"name": "Title",
					"sortIndex": 1,
					"type": "TEXT",
					"groupName": null,
					"options": null
				}
			]
		},
		{
			"id": "d8c05c6b-869c-4c32-8cdc-701a96a077b7",
			"externalId": null,
			"version": 2,
			"createdAt": "2022-03-14T11:02:04.229Z",
			"updatedAt": "2022-03-14T11:02:58.175Z",
			"name": "short term test",
			"startDate": "2022-03-14T11:20:00.000Z",
			"endDate": "2022-03-14T15:20:00.000Z",
			"state": "ACTIVE",
			"customFieldDefinitions": [
				{
					"id": "default_meta_attribute__title",
					"name": "Titel",
					"sortIndex": 1,
					"type": "TEXT",
					"groupName": null,
					"options": null
				},
				{
					"id": "default_meta_attribute__salutation",
					"name": "Eigene Briefanrede",
					"sortIndex": 3,
					"type": "TEXT",
					"groupName": null,
					"options": null
				}
			]
		}
	]`)
		rw.Write(response)
	}
}
