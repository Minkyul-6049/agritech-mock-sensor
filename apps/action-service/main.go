package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Grafana Webhook Payload structure (extract essential requirement)
type GrafanaAlert struct {
	Status string `json:"status"` // "firing" or "resolved"
	Alerts []struct {
		Labels      map[string]string `json:"labels"`
		Annotations map[string]string `json:"annotations"`
	} `json:"alerts"`
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var alert GrafanaAlert
	if err := json.Unmarshal(body, &alert); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	// Alert status "firing" (warning occurrence)
	if alert.Status == "firing" {
		log.Println("🚨 [ALERT RECEIVED] Condition met! Triggering action...")
		//Pump simulation working by log, but actually GPIO or MQTT message would be sent.
		triggerPump()
	} else if alert.Status == "resolved" {
		log.Println("✅ [ALERT RESOLVED] Situation is back to normal.")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook processed successfully"))
}

// Idempotency for function to work water pump
func triggerPump() {
	fmt.Println("=====================================")
	fmt.Println("⚙️ ACTION: Triggering Virtual Water Pump")
	fmt.Println("💧 Status: Water pump is ON for 5 minutes")
	fmt.Println("=====================================")
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	
	port := ":8080"
	log.Printf("Starting Action Service Webhook receiver on port %s...\n", port)
	
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %s\n", err)
	}
}
