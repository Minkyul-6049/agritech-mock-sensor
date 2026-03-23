package main

import (
	"fmt"
	"math/rand"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	// monitor-node의 IP, 새로 발급받은 토큰, 조직(org), 버킷 이름 입력
	url := "http://192.168.202.132:30086"
	token := "y7tYd8StwJpZ9yk9igkFVXUqH8h7-gyLCBof_E1UixFJA4tjqrRJeld9esXkgRhEKfPA8V8AdjOhgUIjQUw2IQ=="
	org := "dutch-agritech"
	bucket := "smartfarm_sensors"

	client := influxdb2.NewClient(url, token)
	writeAPI := client.WriteAPI(org, bucket)

	fmt.Println("🌱 [Farm-Node] 온실 엣지 센서 가동 시작! (데이터 전송 목적지: Monitor-Node)")

	rand.Seed(time.Now().UnixNano())

	// 무한 루프 시작
	for {
		// 온실 온도, 습도, 토양 수분량 랜덤 생성
		temperature := 20.0 + rand.Float64()*5.0
		humidity := 50.0 + rand.Float64()*15.0
		soilMoisture := 30.0 + rand.Float64()*10.0

		// greenhouse_sensors라는 이름으로 데이터 포인트 생성
		p := influxdb2.NewPointWithMeasurement("greenhouse_sensors").
			AddTag("location", "zone-A"). // A구역 온실 꼬리표
			AddField("temperature", temperature).
			AddField("humidity", humidity).
			AddField("soil_moisture", soilMoisture).
			SetTime(time.Now())

		// monitor-node의 InfluxDB로 쏘기!
		writeAPI.WritePoint(p)

		fmt.Printf("[전송 완료] Temp: %.2f℃, Hum: %.2f%%, Soil: %.2f%%\n", temperature, humidity, soilMoisture)

		// --- 여기서부터 제어 로직 --- 

		// 1. 토양 수분량(Soil Moisture) 체크
		if soilMoisture < 35.0 {
			fmt.Println("   🚨 [제어] 토양 수분 35% 미만! 💧 1구역 스프링클러 밸브 개방!")
		} else {
			fmt.Println("   ✅ [상태] 토양 수분 적정 유지 중.")
		}

		// 2. 온도(Temperature) 체크
		if temperature >= 24.0 {
			fmt.Println("   🚨 [제어] 24도 이상 고온 감지! 🌬️ 온실 천장 환풍기 최대 가동!")
		}

		fmt.Println("--------------------------------------------------")

		// 2초 대기 후 다시 처음으로
		time.Sleep(2 * time.Second)
	} // for 루프 끝
} // main 함수 끝
