package client

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/tidwall/pretty"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	adminreports "google.golang.org/api/admin/reports/v1"
	"io/ioutil"
	"net/http"
	"os"
)

func BuildClient(credentialFilePath, impersonationEmail string) (*http.Client, error) {
	// Open our jsonFile
	credentialJsonFile, err := os.Open(credentialFilePath)

	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, errors.New("error reading credential file")
	}

	// Read credential file
	byteValue, err := ioutil.ReadAll(credentialJsonFile)

	// Handle credential file read issues
	if err != nil {
		return nil, errors.New("error reading json from credential file")
	}

	// Define creds
	var creds GoogleServiceAccountCredentials

	// unmarshal into object
	err = json.Unmarshal(byteValue, &creds)

	// return if error during unmarshal
	if err != nil {
		return nil, errors.New("error parsing json from credential file to struct")
	}

	conf := &jwt.Config{
		Email: creds.ClientEmail,
		// The contents of your RSA private key or your PEM file
		// that contains a private key.
		// If you have a p12 file instead, you
		// can use `openssl` to export the private key into a pem file.
		//
		//    $ openssl pkcs12 -in key.p12 -passin pass:notasecret -out key.pem -nodes
		//
		// The field only supports PEM containers with no passphrase.
		// The openssl command will convert p12 keys to passphrase-less PEM containers.
		PrivateKey: []byte(creds.PrivateKey),
		Scopes: []string{
			"https://www.googleapis.com/auth/admin.reports.audit.readonly",
			"https://www.googleapis.com/auth/admin.reports.usage.readonly",
		},
		TokenURL: google.JWTTokenURL,
		// If you would like to impersonate a user, you can
		// create a transport with a subject. The following GET
		// request will be made on the behalf of user@example.com.
		// Optional.
		Subject: impersonationEmail,
	}
	// Initiate an http.Client, the following GET request will be
	// authorized and authenticated on the behalf of user@example.com.
	return conf.Client(context.Background()), nil
}

func ActivitiesList(service *adminreports.Service, eventType, startTime string, endTime string, resultsChannel chan<- string) (int, error) {
	count := 0
	response, err := service.Activities.List("all", eventType).StartTime(startTime).EndTime(endTime).MaxResults(1000).Do()
	if err != nil {
		return 0, err
	}

	// Return if there are no new results
	if len(response.Items) == 0 {
		return 0, nil
	} else {
		// Convert to the activity type
		tmpData := convertActivityTypeToInterface(response.Items)
		count += len(tmpData)
		for _, event := range tmpData {
			// Ugly print the json into a single lined string
			resultsChannel <- string(pretty.Ugly([]byte(event)))
		}
	}

	// Handle paged responses
	for response.NextPageToken != "" {
		response, err := service.Activities.List("all", eventType).StartTime(startTime).EndTime(endTime).MaxResults(1000).PageToken(response.NextPageToken).Do()
		if err != nil {
			return 0, err
		}
		tmpData := convertActivityTypeToInterface(response.Items)
		count += len(tmpData)
		for _, event := range tmpData {
			// Ugly print the json into a single lined string
			resultsChannel <- string(pretty.Ugly([]byte(event)))
		}
	}

	return count, nil
}
