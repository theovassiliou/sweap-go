package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/theovassiliou/sweap-go"
)

// worker takes in a read-only channel to recieve batches of rows.
// After it processes each row-batch it sends out the processed output
// on its channel.
func createGuestsWorker(ctx context.Context, rowBatch <-chan []sweap.Guest, api *sweap.Client, workerNum int) <-chan processed {
	out := make(chan processed)

	go func() {
		cLog := log.WithFields(log.Fields{
			"context": workerNum,
		})
		cLog.Infoln("Created")
		defer close(out)
		p := processed{}
		j := 0
		for rowBatch := range rowBatch {
			j++
			for _, eachGuest := range rowBatch {
				_, err := api.CreateGuest(eachGuest)
				if err != nil {
					p.numErrors++
					cLog.Warnf("[%v] %v", workerNum, err)
				}
				firstName, lastName := processRow(eachGuest)
				// fmt.Printf("Worker %v processes %v of batch: %v %v %v\n", i, k, j, firstName, lastName)
				p.lastNames = append(p.lastNames, lastName)
				p.firstNames = append(p.firstNames, firstName)
				p.numRows++
			}
		}
		out <- p
	}()

	return out
}

func processRow(guest sweap.Guest) (firstName, fullName string) {
	return guest.FirstName, guest.LastName
}
