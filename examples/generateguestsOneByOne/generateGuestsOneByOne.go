package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"

	"github.com/theovassiliou/sweap-go"
)

const numWorkers = 200
const batchSize = 50
const numGuests = 120000
const intraWorkerDelay = 200
const interRowDelay = 1

var eventName = "Load Testing 4"

type processed struct {
	startTime  time.Time
	numRows    int
	firstNames []string
	lastNames  []string
	numErrors  int
	errors     map[string]int
}

type result struct {
	worker           int
	batchSize        int
	intraWorkerDelay int
	interRowDelay    int
	guestCreated     int
	eventName        string
	executionTime    int
	numRows          int
	peopleCount      int
	commonName       string
	commonNameCount  int
	numErrors        int
	avgCps           float32
	errors           map[string]int
}

func (r *result) String() string {
	var sb strings.Builder

	sb.WriteString("Result:\n")
	sb.WriteString(fmt.Sprintf("  Worker: %v\n", r.worker))
	sb.WriteString(fmt.Sprintf("  Batch Size: %d\n", r.batchSize))
	sb.WriteString(fmt.Sprintf("  Delay btw Creates (in ms): %d\n", r.interRowDelay))
	sb.WriteString(fmt.Sprintf("  Delay btw Batches (in ms): %d\n", r.intraWorkerDelay))
	sb.WriteString(fmt.Sprintf("  Guest Created: %d\n", r.guestCreated))
	sb.WriteString(fmt.Sprintf("  Event Name: %s\n", r.eventName))
	sb.WriteString(fmt.Sprintf("  Num Rows: %d\n", r.numRows))
	sb.WriteString(fmt.Sprintf("  People Count: %d\n", r.peopleCount))
	sb.WriteString(fmt.Sprintf("  Common Name: %s\n", r.commonName))
	sb.WriteString(fmt.Sprintf("  Common Name Count: %d\n", r.commonNameCount))
	sb.WriteString(fmt.Sprintf("  Execution Time (in s): %d\n", r.executionTime/1000))
	sb.WriteString(fmt.Sprintf("  AvgGCps: %f\n", r.avgCps))

	sb.WriteString(fmt.Sprintf("  Num Errors: %d\n", r.numErrors))
	sb.WriteString("  Errors:\n")
	for key, value := range r.errors {
		sb.WriteString(fmt.Sprintf("    %d: %s\n", value, key))
	}

	return sb.String()
}

var Log *log.Logger

func NewLogger() *log.Logger {
	if Log != nil {
		return Log
	}

	pathMap := lfshook.PathMap{
		log.InfoLevel:  "./info.log",
		log.ErrorLevel: "./info.log",
		log.WarnLevel:  "./info.log",
	}

	Log = log.New()
	Log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&log.TextFormatter{},
	))
	return Log
}

func main() {

	// ClearGuest List
	api, _ := sweap.New("", "", sweap.OptionDebug(false), sweap.OptionUseStagingEnv(), sweap.OptionEnvFile("../../.stageing-env"))
	search := sweap.EventSearchParameter{
		Name: eventName,
	}

	Log = NewLogger()
	cLog := Log.WithFields(log.Fields{
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

	cLog.Println("Will populate event ", (*events)[0].Name)
	eventID := (*events)[0].ID

	res := result{
		worker:           numWorkers,
		batchSize:        batchSize,
		interRowDelay:    interRowDelay,
		intraWorkerDelay: intraWorkerDelay,
		guestCreated:     numGuests,
		eventName:        eventName,
		errors:           make(map[string]int),
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

	startTime := time.Now()
	// STAGE 2: create a slice of processed output channels with size of numWorkers
	// and assign each slot with the out channel from each worker.
	workersCh := make([]<-chan processed, numWorkers)
	for i := 0; i < numWorkers; i++ {
		// workersCh[i] = createGuestsWorker(ctx, rowsCh, api, i, "")
		workersCh[i] = createGuestsWorker(ctx, rowsCh, api, i, eventID)
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

		for e, p := range processed.errors {
			res.errors[e] = res.errors[e] + p

		}
	}
	res.executionTime = int(time.Since(startTime).Milliseconds())
	res.avgCps = float32(res.numRows / (res.executionTime / 1000))
	fmt.Printf("%v\n", res.String())
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

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}
