package main

import (
	"fmt"
	"math/rand"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	// Monitor-node config
	url := "http://192.168.202.132:30086"
	token := "y7tYd8StwJpZ9yk9igkFVXUqH8h7-gyLCBof_E1UixFJA4tjqrRJeld9esXkgRhEKfPA8V8AdjOhgUIjQUw2IQ=="
	org := "dutch-agritech"
	bucket := "smartfarm_sensors"

	client := influxdb2.NewClient(url, token)
	writeAPI := client.WriteAPI(org, bucket)

	fmt.Println("🌱 [Farm-Node] Greenhouse edge sensor operation start!")

	rand.Seed(time.Now().UnixNano())

	for {
		// Generate mock data
		temperature := 20.0 + rand.Float64()*5.0
		humidity := 50.0 + rand.Float64()*15.0
		soilMoisture := 30.0 + rand.Float64()*10.0

		// Fixed: InfluxDB point creation (Dot chaining)
		p := influxdb2.NewPointWithMeasurement("greenhouse_sensors").
			AddTag("location", "zone-A").
			AddField("temperature", temperature).
			AddField("humidity", humidity).
			AddField("soil_moisture", soilMoisture).
			SetTime(time.Now())

		// Write to InfluxDB
		writeAPI.WritePoint(p)

		fmt.Printf("[Data sent] Temp: %.2f℃, Hum: %.2f%%, Soil: %.2f%%\n", temperature, humidity, soilMoisture)

		// --- Control logic ---

		// 1. Soil Moisture
		if soilMoisture < 35.0 {
			fmt.Println("   🚨 [Control] Soil Moisture low! 💧 Opening sprinkler...")
		} else {
			fmt.Println("   ✅ [Status] Soil Moisture optimal.")
		}

		// 2. Temperature
		if temperature >= 24.0 {
			fmt.Println("   🚨 [Control] High temp detected! 🌬️ Running ventilation...")
		}

		// 3. Humidity (Fixed: fmt.Println instead of PrintIn)
		if humidity >= 70.0 {
			fmt.Println("   🚨 [Control] High humidity! Operating dehumidifier...")
		}

		fmt.Println("--------------------------------------------------")
		time.Sleep(2 * time.Second)
	}
}
