package main

import (
	"fmt"
	"math/rand"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	// Monitor-node configuration: IP, token, organization, and bucket name
	url := "http://192.168.202.132:30086"
	token := "y7tYd8StwJpZ9yk9igkFVXUqH8h7-gyLCBof_E1UixFJA4tjqrRJeld9esXkgRhEKfPA8V8AdjOhgUIjQUw2IQ=="
	org := "dutch-agritech"
	bucket := "smartfarm_sensors"

	client := influxdb2.NewClient(url, token)
	writeAPI := client.WriteAPI(org, bucket)

	fmt.Println("🌱 [Farm-Node] Greenhouse edge sensor operation start! (Target: Monitor-Node)")

	rand.Seed(time.Now().UnixNano())

	// Infinite loop for data transmission and edge control
	for {
		// Generate random sensor data (simulating real sensors)
		temperature := 20.0 + rand.Float64()*5.0
		humidity := 50.0 + rand.Float64()*15.0
		soilMoisture := 30.0 + rand.Float64()*10.0

		// Create a new data point for InfluxDB (Fixed the chaining issue here)
		p := influxdb2.NewPointWithMeasurement("greenhouse_sensors").
			AddTag("location", "zone-A").
			AddField("temperature", temperature).
			AddField("humidity", humidity).
			AddField("soil_moisture", soilMoisture).
			SetTime(time.Now())

		// Write data to Monitor-Node
		writeAPI.WritePoint(p)

		fmt.Printf("[Data sent] Temp: %.2f℃, Hum: %.2f%%, Soil: %.2f%%\n", temperature, humidity, soilMoisture)

		// --- Edge Control Logic ---

		// 1. Soil Moisture Check: Triggers sprinkler if too dry
		if soilMoisture < 35.0 {
			fmt.Println("   🚨 [Control] Soil Moisture under 35% detected! 💧 Opening Zone-A sprinkler valve!")
		} else {
			fmt.Println("   ✅ [Status] Soil Moisture is optimal.")
		}

		// 2. Temperature Check: Triggers ventilation if too hot
		if temperature >= 24.0 {
			fmt.Println("   🚨 [Control] Temperature over 24℃ detected! 🌬️ Operating greenhouse ceiling ventilation!")
		}

		// 3. Humidity Check: Triggers dehumidifier if too humid (Fixed Println typo here)
		if humidity >= 70.0 {
			fmt.Println("   🚨 [Control] Humidity over 70% detected! 🌀 Operating dehumidifier!")
		}

		fmt.Println("--------------------------------------------------")

		// Wait for 2 seconds before the next reading
		time.Sleep(2 * time.Second)
	}
}
