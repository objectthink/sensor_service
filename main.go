package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
	"github.com/nats-io/nats.go"
)

type SensorStatus struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Temperature int    `json:"temperature"`
	High        int    `json:"high"`
	Low         int    `json:"low"`
}

func natsSetup(writer api.WriteAPIBlocking) {
	// Connect to a server
	nc, _ := nats.Connect("nats://192.168.86.31:4222")

	// Simple Publisher
	nc.Publish("foo", []byte("hello cruel world"))

	// Simple Async Subscriber
	nc.Subscribe("*.status", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))

		// Or write directly line protocol
		//writeAPI := influx.WriteAPIBlocking("home", "test_bucket")

		var sensorStatus SensorStatus
		err := json.Unmarshal(m.Data, &sensorStatus)
		if err != nil {

			// if error is not nil
			// print error
			fmt.Println(err)
		}

		//fmt.Printf("json: %d %d %d", sensorStatus.Temperature, sensorStatus.Humidity, sensorStatus.Level)
		//"stat,unit=temperature avg=%f,max=%f"
		line := fmt.Sprintf("temperature,name=\"%s\",location=\"%s\" current=%d",
			sensorStatus.Name,
			strings.ReplaceAll(sensorStatus.Location, " ", "_"),
			sensorStatus.Temperature)

		writer.WriteRecord(context.Background(), line)
	})
}

func main() {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient("http://192.168.86.31:8086", "U2Efox5YfRR8c_oCEhwO4xrlJ_Cgjj6DGDG5kqM_GB0f43o80KL-d2oUGPodgv4vXVNfI_BfK--OmpdXOZYnxg==")

	writer := client.WriteAPIBlocking("home", "sensors")

	natsSetup(writer)

	done := make(chan bool)
	go forever()
	<-done // Block forever

	//fmt.Scanln()
}

func forever() {
	for {
		fmt.Printf("%v+\n", time.Now())
		time.Sleep(time.Second)
	}
}

/*
func mainx() {
	fmt.Println("hello cruel world")

	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient("http://192.168.86.31:8086", "U2Efox5YfRR8c_oCEhwO4xrlJ_Cgjj6DGDG5kqM_GB0f43o80KL-d2oUGPodgv4vXVNfI_BfK--OmpdXOZYnxg==")
	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking("home", "test_bucket")
	// Create point using full params constructor
	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": 24.5, "max": 45.0},
		time.Now())
	// write point immediately
	writeAPI.WritePoint(context.Background(), p)
	// Create point using fluent style
	p = influxdb2.NewPointWithMeasurement("stat").
		AddTag("unit", "temperature").
		AddField("avg", 23.2).
		AddField("max", 45.0).
		SetTime(time.Now())
	writeAPI.WritePoint(context.Background(), p)

	// Or write directly line protocol
	line := fmt.Sprintf("stat,unit=temperature avg=%f,max=%f", 23.5, 45.0)
	writeAPI.WriteRecord(context.Background(), line)

	// Get query client
	queryAPI := client.QueryAPI("home")
	// Get parser flux query result
	result, err := queryAPI.Query(context.Background(), `from(bucket:"test_bucket")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "stat")`)
	if err == nil {
		// Use Next() to iterate over query result lines
		for result.Next() {
			// Observe when there is new grouping key producing new table
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// read result
			fmt.Printf("row: %s\n", result.Record().String())
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	}
	// Ensures background processes finishes
	//client.Close()

	natsSetup(client)

	fmt.Scanln()
}
*/

