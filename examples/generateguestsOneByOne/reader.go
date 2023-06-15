package main

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	rn "github.com/random-names/go"
	"github.com/theovassiliou/sweap-go"
)

// reader creates and returns a channel that recieves
// batches of rows (of length batchSize) from the file
func randomGuestGeneratorReader(ctx context.Context, rowsBatch *[]sweap.Guest, eventID string) <-chan []sweap.Guest {
	out := make(chan []sweap.Guest)
	guestList := genRandomGuests(eventID, numGuests)
	nextGuest := sweap.NewGuestIterator(guestList)

	go func() {
		defer close(out) // close channel when we are done sending all rows

		for {
			scanned := nextGuest()

			select {
			case <-ctx.Done():
				return
			default:
				row := scanned
				// if batch size is complete or end of file, send batch out
				if len(*rowsBatch) == batchSize || scanned == nil {
					out <- *rowsBatch
					*rowsBatch = []sweap.Guest{} // clear batch
				}
				if row != nil {
					*rowsBatch = append(*rowsBatch, *row) // add row to current batch
				}
			}

			// if nothing else to scan return
			if scanned == nil {
				return
			}
		}
	}()

	return out
}

func genRandomGuests(eventID string, i int) []sweap.Guest {
	fns, lns, err := getRandomNames(i)

	if err != nil {
		return []sweap.Guest{}
	}

	gl := make([]sweap.Guest, 0)

	for i := range fns {
		gl = append(gl, sweap.Guest{
			ID:              fmt.Sprintf("%v", i),
			Version:         0,
			UpdatedAt:       nil,
			CreatedAt:       nil,
			ExternalID:      nil,
			EventID:         eventID,
			FirstName:       fns[i],
			LastName:        lns[i],
			EntourageCount:  0,
			Comment:         nil,
			Email:           "",
			CustomFields:    sweap.CustomFields{},
			InvitationID:    "",
			InvitationState: sweap.NONE,
			TicketID:        "",
			ParentGuestID:   "",
			CategoryID:      "",
			AttendanceState: sweap.NONEATTENDANCE,
		})
	}
	return gl

}

func getRandomNames(i int) (first, last []string, err error) {
	lastNames, err := rn.GetRandomNames("census-90/all.last", &rn.Options{Number: i})

	if err != nil {
		fmt.Println(err)
	}

	firstNames, err := rn.GetRandomNames("census-90/male.first", &rn.Options{Number: i})

	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(lastNames); i++ {
		lastNames[i] = capitalize(lastNames[i])
		firstNames[i] = capitalize(firstNames[i])
	}

	return firstNames, lastNames, err
}

func capitalize(str string) string {
	str = strings.ToLower(str)
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
