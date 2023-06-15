package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/theovassiliou/sweap-go"
)

// worker takes in a read-only channel to recieve batches of rows.
// After it processes each row-batch it sends out the processed output
// on its channel.
func createGuestsWorker(ctx context.Context, rowBatch <-chan []sweap.Guest, api *sweap.Client, workerNum int, eventID string) <-chan processed {
	out := make(chan processed)
	guestToCreatePerBatch := numGuests / numWorkers

	go func() {
		cLog := Log.WithFields(log.Fields{
			"worker": workerNum,
		})
		defer close(out)
		p := processed{
			errors:    make(map[string]int),
			startTime: time.Now(),
		}
		j := 0
		for rowBatch := range rowBatch {
			j++
			currentBatchstartTime := time.Now()
			startGuests := p.numRows
			for _, eachGuest := range rowBatch {
				_, err := api.CreateGuest(eachGuest)
				time.Sleep(interRowDelay * time.Millisecond)

				if err != nil {
					p.numErrors++
					p.errors[err.Error()] = p.errors[err.Error()] + 1
					cLog.Warnf("[%v] %v", workerNum, err)
				}
				firstName, lastName := processRow(eachGuest)
				p.lastNames = append(p.lastNames, lastName)
				p.firstNames = append(p.firstNames, firstName)
				p.numRows++
			}
			cLog.Infof("%v-%v/%v -> AvgGCps[actual/mean]: %v/%v\n", p.numRows-startGuests, j*batchSize, guestToCreatePerBatch, float32(float32(p.numRows-startGuests)/(float32(time.Since(currentBatchstartTime).Milliseconds()/1000))), float32(float32(p.numRows)/(float32(time.Since(p.startTime).Milliseconds()/1000))))

			time.Sleep(intraWorkerDelay * time.Millisecond)

		}
		out <- p
	}()

	return out
}

func createBulkImportWorker(ctx context.Context, rowBatch <-chan []sweap.Guest, api *sweap.Client, workerNum int, eventID string) <-chan processed {
	out := make(chan processed)

	go func() {
		cLog := log.WithFields(log.Fields{
			"context": workerNum,
		})
		defer close(out)

		p := processed{
			errors: make(map[string]int),
		}
		j := 0

		gbiNew := sweap.GuestBulkImport{
			ID:                     "",
			Version:                0,
			ExternalID:             "1234",
			Name:                   fmt.Sprintf("Worker %v created", workerNum),
			Guests:                 nil,
			EventId:                eventID,
			CustomFieldDefinitions: []sweap.CustomFieldDefinitions{},
		}

		gbi, err := api.CreateGuestBulkImportObject(gbiNew)
		if err != nil {
			cLog.Warnf("%s\n", err)
			return
		}
		cLog.Debug(pp(gbi))

		for rowBatch := range rowBatch {
			j++

			gbiID := gbi.ID
			err = api.BulkImportUpdateBatch(gbiID, rowBatch)
			cLog.Infof("Uploading batch %v for worker %v", j, workerNum)
			if err != nil {
				p.numErrors++
				p.errors[err.Error()] = p.errors[err.Error()] + 1
				cLog.Warnf("[%v] %v", workerNum, err)
			}

			for _, eachGuest := range rowBatch {
				firstName, lastName := processRow(eachGuest)
				p.lastNames = append(p.lastNames, lastName)
				p.firstNames = append(p.firstNames, firstName)
				p.numRows++
			}

		}

		api.BulkImportFinishUpload(gbi.ID)
		cLog.Infof("Uploading for %v/%v", j, workerNum)

		if err != nil {
			p.numErrors++
			p.errors[err.Error()] = p.errors[err.Error()] + 1
			cLog.Warnf("[%v] %v", workerNum, err)
		}

		out <- p
	}()

	return out
}

func processRow(guest sweap.Guest) (firstName, fullName string) {
	return guest.FirstName, guest.LastName
}

func pp(intf interface{}) string {

	empJSON, err := json.MarshalIndent(intf, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(empJSON)
}
