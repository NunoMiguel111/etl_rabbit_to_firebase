package main

import (
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func produce_to_rabbit(channel *amqp091.Channel, exchange string, route string, measurement Measurement) error {

	body, err := json.Marshal(measurement)
	if err != nil {
		return fmt.Errorf("Failed to marshall measurement to json: %v", err)
	}
	// publishing a message
	err = channel.Publish(
		exchange, // exchange
		route,    // key
		false,    // mandatory
		false,    // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	return err
}

/*

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

	produce_to_rabbit(channel, "sensor.data", "sensor.data", measurement)
*/
