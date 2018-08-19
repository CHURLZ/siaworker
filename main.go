package main

// A worker listens for events and sends these on to a rabbitMQ instance.

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	time.Sleep(5 * time.Second)
	fmt.Println("Done waiting, let's go.")
}

func main() {
	// c := LoadConfig()
	// fmt.Println(c.Account)
	// fmt.Println(c.Events)

	username := os.Getenv("IP_USER")
	passwd := os.Getenv("IP_PASSWD")
	url := os.Getenv("IP_URL")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	failOnError(err, "Failed to connect to URL.")

	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("Connection returned: %s", resp.Status))
	}

	rabbit := Connect()
	defer rabbit.Close()

	var buffer []byte
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	for {
		line, _ := reader.ReadBytes('\n')
		if strings.Contains(string(line), "--boundary") {
			if len(buffer) == 0 {
				continue
			}
			evt, err := UnmarshalBuffer(buffer)
			warnOnError(err, "error Unmarshaling buffer.")

			// write the event to log
			// fmt.Printf("%#v\n", evt)

			data, err := evt.MarshalBuffer()
			warnOnError(err, "error Marshaling buffer.")
			rabbit.Publish(data)

			// clear the buffer
			buffer = []byte{}
		}
		buffer = append(buffer[:], line[:]...)
	}
}
