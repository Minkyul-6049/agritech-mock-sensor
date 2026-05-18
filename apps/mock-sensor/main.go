package main

import (
        "fmt"
        "math/rand"
        "time"

        influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
        // 1. Connection config for the Control Tower (Monitor-Node)
        url := "http://192.168.202.132:30086"
        token := "y7tYd8StwJpZ9yk9igkFVXUqH8h7-gyLCBof_E1UixFJA4tjqrRJeld9esXkgRhEKfPA8V8AdjOhgUIjQUw2IQ=="
        org := "dutch-agritech"
        bucket := "smartfarm_sensors"

        client := influxdb2.NewClient(url, token)
        writeAPI := client.WriteAPI(org, bucket)

        fmt.Println("🌱 [Farm-Node] Greenhouse edge sensor operation started!")

        rand.Seed(time.Now().UnixNano())

        for {
                // 2. Generate mock sensor data
                temperature := 20.0 + rand.Float64()*5.0
                humidity := 50.0 + rand.Float64()*15.0
                soilMoisture := 30.0 + rand.Float64()*10.0

                // Pack the data with location tag
                p := influxdb2.NewPointWithMeasurement("greenhouse_sensors").
                        AddTag("location", "zone-A").
                        AddField("temperature", temperature).
                        AddField("humidity", humidity).
                        AddField("soil_moisture", soilMoisture).
                        SetTime(time.Now())

                // Send to InfluxDB
                writeAPI.WritePoint(p)

                fmt.Printf("[Data sent] Temp: %.2f℃, Hum: %.2f%%, Soil: %.2f%%\n", temperature, humidity, soilMoisture)

                // --------------------------------------------------
                // 3. Local Control Logic (Edge Level)
                // This hardware runs automatically even if the cloud connection is lost.
                // --------------------------------------------------

                // Soil Moisture Check
                if soilMoisture < 35.0 {
                        fmt.Println("   🚨 [Action] Soil Moisture low! 💧 Turning on sprinkler...")
                } else {
                        fmt.Println("   ✅ [Status] Soil Moisture optimal.")
                }

                // Temperature Check
                if temperature >= 24.0 {
                        fmt.Println("   🚨 [Action] High temp detected! 🌬️ Turning on ventilation fan...")
                }

                // Humidity Check
                if humidity >= 70.0 {
                        fmt.Println("   🚨 [Action] High humidity! 💨 Turning on dehumidifier...")
                }

                fmt.Println("--------------------------------------------------")
                time.Sleep(2 * time.Second)
        }
}
