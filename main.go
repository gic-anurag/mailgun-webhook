package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	APIKey    string `json:"apiKey"`
	Domain    string `json:"domain"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func SendSimpleMessage(config *Config) (string, error) {
	// Send the message using your preferred email provider (Mailgun, SendGrid, etc.)
	// Replace this code with your implementation
	return "MESSAGE_ID", nil
}

func main() {
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err.Error())
	}

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		handleWebhookEvent(w, r, config)
	})

	http.HandleFunc("/send-email", func(w http.ResponseWriter, r *http.Request) {
		sendEmail(w, r, config)
	})

	log.Fatal(http.ListenAndServe(":8281", nil))
}

func sendEmail(w http.ResponseWriter, r *http.Request, config *Config) {
	id, err := SendSimpleMessage(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Email sent! Message ID: %s", id)
}

func handleWebhookEvent(w http.ResponseWriter, r *http.Request, config *Config) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var payload map[string]interface{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Printf("Failed to parse webhook event: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Received payload: %v", payload)

	eventData, ok := payload["event-data"].(map[string]interface{})
	if !ok {
		log.Println("Invalid event data received")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	eventType, ok := eventData["event"].(string)
	if !ok {
		log.Println("Invalid event type received")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch eventType {
	case "delivered":
		handleDeliveredEvent(eventData)

	case "opened":
		handleOpenedEvent(eventData)

	case "clicked":
		handleClickedEvent(eventData)

	case "failed":
		handleFailedEvent(eventData)

	default:
		log.Printf("Unknown webhook event received: %s", eventType)
	}

	w.WriteHeader(http.StatusOK)
}

func handleDeliveredEvent(eventData map[string]interface{}) {
	message, ok := getMessageFromEventData(eventData)
	if !ok {
		log.Println("Invalid message data received")
		return
	}

	messageID := getMessageIDFromMessage(message)
	if messageID != "" {
		log.Printf("Email delivered: Message ID - %s", messageID)
	} else {
		log.Println("Invalid message ID received")
	}
}

func getMessageFromEventData(eventData map[string]interface{}) (map[string]interface{}, bool) {
	messageData, ok := eventData["message"].(map[string]interface{})
	if !ok {
		log.Println("Invalid message data received")
		return nil, false
	}

	message, ok := messageData["headers"].(map[string]interface{})
	if !ok {
		log.Println("Invalid message headers received")
		return nil, false
	}

	return message, true
}

func getMessageIDFromMessage(message map[string]interface{}) string {
	messageID, ok := message["message-id"].(string)
	if !ok {
		log.Println("Invalid message ID received")
		return ""
	}

	return messageID
}

func handleOpenedEvent(eventData map[string]interface{}) {
	// Handle opened event
	log.Println("Email opened")
}

func handleClickedEvent(eventData map[string]interface{}) {
	// Handle clicked event
	log.Println("Email clicked")
}

func handleFailedEvent(eventData map[string]interface{}) {
	// Handle failed event
	log.Println("Email failed")
}
