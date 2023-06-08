package main

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/theovassiliou/sweap-go"
)

const numWorkers = 200
const batchSize = 100
const numGuests = 1

type processed struct {
	numRows    int
	firstNames []string
	lastNames  []string
	numErrors  int
}

func main() {

	// ClearGuest List
	api, _ := sweap.New("", "", sweap.OptionDebug(false), sweap.OptionUseStagingEnv(), sweap.OptionEnvFile("../../.stageing-env"))
	search := sweap.EventSearchParameter{
		Name: "Demo Event for Bulk 2",
	}

	cLog := log.WithFields(log.Fields{
		"context": "main",
	})

	cLog.Infof("Started. Searching for %v\n", search.Name)

	events, err := api.SearchEvents(search)

	if err != nil {
		cLog.Printf("%s\n", err)
		return
	}

	if len(*events) == 0 {
		cLog.Printf("No event machting %+v found\n", search)
		return
	}

	cLog.Println("Handling event ", (*events)[0].Name)
	eventID := (*events)[0].ID

	type result struct {
		worker          int
		batchSize       int
		guestCreated    int
		numRows         int
		peopleCount     int
		commonName      string
		commonNameCount int
		numErrors       int
	}

	res := result{
		worker:       numWorkers,
		batchSize:    batchSize,
		guestCreated: numGuests,
	}

	// create a main context, and call cancel at the end, to ensure all our
	// goroutines exit without leaving leaks.
	// Particularly, if this function becomes part of a program with
	// a longer lifetime than this function.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// STAGE 1: start reader
	rowsBatch := []sweap.Guest{}
	rowsCh := randomGuestGeneratorReader(ctx, &rowsBatch, eventID)

	// STAGE 2: create a slice of processed output channels with size of numWorkers
	// and assign each slot with the out channel from each worker.
	workersCh := make([]<-chan processed, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workersCh[i] = createGuestsWorker(ctx, rowsCh, api, i)
	}

	firstNameCount := map[string]int{}
	lastNameCount := map[string]bool{}

	// STAGE 3: read from the combined channel and calculate the final result.
	// this will end once all channels from workers are closed!
	for processed := range genericCombiner(ctx, workersCh...) {
		// add number of rows processed by worker
		res.numRows += processed.numRows
		res.numErrors += processed.numErrors

		// use full names to count people
		for i, lastName := range processed.lastNames {
			lastNameCount[processed.firstNames[i]+lastName] = true
		}
		res.peopleCount = len(lastNameCount)

		// update most common first name based on processed results
		for _, firstName := range processed.firstNames {
			firstNameCount[firstName]++

			if firstNameCount[firstName] > res.commonNameCount {
				res.commonName = firstName
				res.commonNameCount = firstNameCount[firstName]
			}
		}
	}

	fmt.Printf("%#v\n", res)
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}
