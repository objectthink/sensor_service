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
	Humidity    int    `json:"humidity"`
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

		line = fmt.Sprintf("humidity,name=\"%s\",location=\"%s\" current=%d",
			sensorStatus.Name,
			strings.ReplaceAll(sensorStatus.Location, " ", "_"),
			sensorStatus.Humidity)

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
		time.Sleep(time.Second)
	}
}
