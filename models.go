package main

import (
	"time"
)

type MeasurementType string

const (
	WaterTemperature MeasurementType = "water_temperature"
	AirTemperature   MeasurementType = "air_temperature"
	TDS              MeasurementType = "tds"
	pH               MeasurementType = "ph"
)

type MeasurementValue struct {
	Type  MeasurementType `json:"type"`
	Value float64         `json:"value"`
}

type Measurement struct {
	Timestamp    time.Time          `json:"timestamp"`
	Location     [2]float64         `json:"location"`
	Measurements []MeasurementValue `json:"measurements"`
}
