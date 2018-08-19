package main

// A worker listens for events and sends these on to a rabbitMQ instance.

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"
)

var config Config

func ConnectStream() *http.Response {
	url := os.Getenv("IP_URL")
	username := os.Getenv("IP_USER")
	passwd := os.Getenv("IP_PASSWD")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	failOnError(err, "Failed to connect to URL.")

	if resp.StatusCode != 200 {
		failOnError(nil, "Connection returned: "+resp.Status)
	}

	return resp
}

func main() {
	log.Println("Worker started.")
	config := LoadConfig()

	response := ConnectStream()
	defer response.Body.Close()

	rabbit := Connect()
	defer rabbit.Close()

	reader := bufio.NewReader(response.Body)
	var buffer []byte
	for {
		line, _ := reader.ReadBytes('\n')
		if strings.Contains(string(line), "--boundary") {
			if len(buffer) == 0 {
				continue
			}
			evt, err := UnmarshalBuffer(buffer)
			warnOnError(err, "error Unmarshaling buffer.")

			if evt.QualifiesForPublish(*config) {
				evt.Prime(*config)
				data, err := evt.Marshal()
				warnOnError(err, "error Marshaling buffer.")

				rabbit.Publish(data)
				log.Printf("Published evt to message queue: %s - %s\n", evt.EventType, evt.EventState)
			} else {
				log.Printf("Event triggered but not forwarded: %s - %s\n", evt.EventType, evt.EventState)
			}
			// clear the buffer
			buffer = []byte{}
		}
		buffer = append(buffer[:], line[:]...)
	}
}
