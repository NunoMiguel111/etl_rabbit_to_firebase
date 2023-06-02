package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func initializeFirebaseApp(ctx context.Context, serviceAccountKeyPath string) (*firebase.App, error) {
	opt := option.WithCredentialsFile(serviceAccountKeyPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize Firebase app: %v", err)
	}

	return app, nil
}

// Currently there is no authentication in place
func insertMeasurement(ctx context.Context, app *firebase.App, measurement Measurement, wg *sync.WaitGroup, url string) {
	defer wg.Done()

	payload, err := json.Marshal(measurement)
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to store measurement data in Firebase. Status code: %v\n", resp.StatusCode)
		return
	}

	fmt.Println("Measurement insert successful")
}

func batchInsertMeasurements(ctx context.Context, app *firebase.App, measurements []Measurement, url string) {
	var wg sync.WaitGroup
	wg.Add(len(measurements))

	for _, measurement := range measurements {
		go insertMeasurement(ctx, app, measurement, &wg, url)
	}

	wg.Wait()
	fmt.Println("Batch insert successful")
}
