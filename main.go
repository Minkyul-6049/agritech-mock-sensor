package main

import (
	"fmt"
	"math/rand"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	// monitor-node의 IP, new generated token, organization, bucket name
	url := "http://192.168.202.132:30086"
	token := "y7tYd8StwJpZ9yk9igkFVXUqH8h7-gyLCBof_E1UixFJA4tjqrRJeld9esXkgRhEKfPA8V8AdjOhgUIjQUw2IQ=="
	org := "dutch-agritech"
	bucket := "smartfarm_sensors"

	client := influxdb2.NewClient(url, token)
	writeAPI := client.WriteAPI(org, bucket)

	fmt.Println("🌱 [Farm-Node]Greenhouse edge sensor operation start! (data transmission target: Monitor-Node)")

	rand.Seed(time.Now().UnixNano())

	// infinite loop for data transmission
	for {
		// Generate random sensor data(temperature, humidity & soil moisture)
		temperature := 20.0 + rand.Float64()*5.0
		humidity := 50.0 + rand.Float64()*15.0
		soilMoisture := 30.0 + rand.Float64()*10.0

		// greenhouse_sensors Create a new data point for InfluxDB
		p := influxdb2.NewPointWithMeasurement("greenhouse_sensors").
			AddTag("location", "zone-A"). // Zone-A Location
			AddField("temperature", temperature).
			AddField("humidity", humidity).
			AddField("soil_moisture", soilMoisture).
			SetTime(time.Now())

		// Write data to Monitor-Node
		writeAPI.WritePoint(p)

		fmt.Printf("[Data sent] Temp: %.2f℃, Hum: %.2f%%, Soil: %.2f%%\n", temperature, humidity, soilMoisture)

		// --- Control logic --- 

		// 1. Soil Moisture Check
		if soilMoisture < 35.0 {
			fmt.Println("   🚨 [Control] Soil Moisture under 35 % detected! 💧Opening Zone-A sprinkler valve !")
		} else {
			fmt.Println("   ✅ [Status] Soil Moisture is optimal.")
		}

		// 2. Temperature check
		if temperature >= 24.0 {
			fmt.Println("   🚨 [Control] Temperature over 24 detected ! 🌬️ Greenhouse ceiling ventilation operating!")
		}

		// 3. humidity check
		if humidity >= 70.0 {
			fmt.Println(" [Control] humidity over 70% detected ! operating dehumidifier!")
		}
	
		fmt.Println("--------------------------------------------------")

		// time.Sleep(2* time.Second)
		time.Sleep(2 * time.Second)
	} // End of for loop
} // End of main function
