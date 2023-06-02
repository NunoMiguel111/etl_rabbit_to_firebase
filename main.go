package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	// Load .env file variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Coudlnt load .env variable: %v", err)
	}
	host := os.Getenv("RABBIT_HOST")
	port := os.Getenv("RABBIT_PORT")
	username := os.Getenv("RABBIT_USERNAME")
	password := os.Getenv("RABBIT_PASSWORD")
	connection_string := "amqps://" + username + ":" + password + "@" + host + ":" + port + "/" + username

	fmt.Println(connection_string)

	// connect to RabbitMQ brooker
	conn, err := amqp091.Dial(connection_string)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	defer conn.Close()
	// Open a channel with remote brooker
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	defer channel.Close()
	// Declare queue info
	queue, err := channel.QueueDeclare(
		"sensor.data.queue", // Queue name
		true,                //Durable
		false,               //Delete when unused
		false,               //Exclusive
		false,               // No wait
		nil,
	)

	if err != nil {
		log.Fatalf("Failed to declare a Queue: %v", err)
	}
	/*
		fmt.Println(queue.Name)
		// Create a new measurement
		measurement := Measurement{
			Timestamp: time.Now(),
			Location:  [2]float64{51.5074, -0.1278}, // Example location coordinates (London)
			Measurements: []MeasurementValue{
				{
					Type:  WaterTemperature,
					Value: 25.5,
				},
				{
					Type:  pH,
					Value: 7.2,
				},
			},
		}

		produce_to_rabbit(channel, "teste_1", "microderivador1", measurement)
	*/

	ctx := context.Background()
	cred := os.Getenv("FIREBASE_JSON_CRED_PATH")
	url := os.Getenv("FIREBASE_RTDB_URL")
	firebaseApp, err := initializeFirebaseApp(ctx, cred)
	if err != nil {
		log.Fatalf("Coudln't initialize firebaseApp")
	}

	messages, err := consume_from_rabbit(channel, queue.Name)

	measurements := make([]Measurement, len(messages))

	// Define time to wait for messages before sending
	timeToWaitForMessages := 5 * time.Second

	// Flag that symbolizes if we should terminate to loop over the messages
	var loop bool = true

	// Process the messages
	for loop {
		select {
		case message, ok := <-messages:
			if !ok {
				// Channel closed, no more messages
				fmt.Println("Channel cloce, no more messages")
				loop = false
				break
			}

			// Unmarshal the message body into a struct
			var measurement Measurement
			err := json.Unmarshal(message.Body, &measurement)
			fmt.Println(measurement)
			if err != nil {
				log.Printf("Failed to unmarshal message body: %v", err)
				continue
			}

			measurements = append(measurements, measurement)

		case <-time.After(timeToWaitForMessages):
			// Timeout reached, no more messages within the specified time
			fmt.Println("Timeout reached.. No more messages processed")

			loop = false
		}

	}
	batchInsertMeasurements(ctx, firebaseApp, measurements, url+"measurements.json")
	if err != nil {
		log.Fatalf("Error while inserting measurements in database: %v", err)
	}

	fmt.Println("Program terminated.")
}
