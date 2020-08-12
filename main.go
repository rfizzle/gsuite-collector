package main

import (
	"context"
	"github.com/rfizzle/collector-helpers/outputs"
	"github.com/rfizzle/collector-helpers/state"
	"github.com/rfizzle/gsuite-collector/client"
	"github.com/spf13/viper"
	adminreports "google.golang.org/api/admin/reports/v1"
	"google.golang.org/api/option"
	"log"
	"os"
	"time"
)

func main() {
	// Setup variables
	var maxMessages = int64(5000)

	// Setup Parameters via CLI or ENV
	if err := setupCliFlags(); err != nil {
		log.Fatalf("missing required paramters %v", err.Error())
	}

	// Setup log writer
	tmpWriter, err := outputs.NewTmpWriter()
	if err != nil {
		log.Fatalf("%v\n", err.Error())
	}

	// Setup the channels for handling async messages
	chnMessages := make(chan string, maxMessages)

	// Setup the Go Routine
	pollTime := viper.GetInt("schedule")

	// Start Poll
	go pollEvery(pollTime, chnMessages, tmpWriter)

	// Handle messages in the channel (this will keep the process running indefinitely)
	for message := range chnMessages {
		handleMessage(message, tmpWriter)
	}
}

func pollEvery(seconds int, resultsChannel chan<- string, tmpWriter *outputs.TmpWriter) {
	var currentState *state.State
	var err error

	// Setup State
	if state.Exists(viper.GetString("state-path")) {
		currentState, err = state.Restore(viper.GetString("state-path"))
		if err != nil {
			log.Fatalf("Error getting state: %v\n", err.Error())
		}
	} else {
		currentState = state.New()
	}

	for {
		log.Println("Getting data...")

		// Get events
		eventCount, lastPollTime := getEvents(currentState.LastPollTimestamp, resultsChannel)

		// Copy tmp file to correct outputs
		if eventCount > 0 {
			// Close and rotate file
			_ = tmpWriter.Rotate()

			if err := outputs.WriteToOutputs(tmpWriter.LastFilePath, lastPollTime.Format(time.RFC3339)); err != nil {
				log.Fatalf("Unable to write to output: %v", err)
			}

			// Remove temp file now
			err := os.Remove(tmpWriter.LastFilePath)
			if err != nil {
				log.Fatalf("Unable to remove tmp file: %v", err)
			}
		}

		// Let know that event has been processes
		log.Printf("%v events processed...\n", eventCount)

		// Update state
		currentState.LastPollTimestamp = lastPollTime.Format(time.RFC3339)
		state.Save(currentState, viper.GetString("state-path"))

		// Wait for x seconds until next poll
		<-time.After(time.Duration(seconds) * time.Second)
	}
}

func getEvents(timestamp string, resultChannel chan<- string) (int, time.Time) {
	// Get current time
	now := time.Now()

	// Build an HTTP client with JWT header
	googleClient, err := client.BuildClient(viper.GetString("gsuite-credentials"), viper.GetString("impersonated-user"))
	if err != nil {
		log.Fatalf("Unable to build client: %v", err)
	}

	// Create a new service client with the built HTTP client
	srv, err := adminreports.NewService(context.Background(), option.WithHTTPClient(googleClient))
	if err != nil {
		log.Fatalf("Unable to retrieve reports client %v", err)
	}

	// Define static event types
	var eventTypes = []string{"admin", "calendar", "drive", "login", "mobile", "token", "groups", "saml", "chat", "gplus", "rules", "jamboard", "meet", "user_accounts", "access_transparency", "groups_enterprise", "gcp"}

	// Initialize data array
	dataCount := 0

	// Loop through event types
	for _, eventType := range eventTypes {
		if viper.GetBool("verbose") {
			log.Printf("Getting event type %s\n", eventType)
		}
		resultSize, err := client.ActivitiesList(srv, eventType, timestamp, resultChannel)
		if err != nil {
			log.Fatalf("Unable to retrieve activities list for %s: %v", eventType, err)
		}

		dataCount += resultSize
	}

	return dataCount, now
}

// Handle message in a channel
func handleMessage(message string, tmpWriter *outputs.TmpWriter) {
	if err := tmpWriter.WriteLog(message); err != nil {
		log.Fatalf("Unable to write to temp file: %v", err)
	}
}
