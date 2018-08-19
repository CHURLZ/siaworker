package main

import (
	"encoding/xml"
	"fmt"
	"io"
)

type Event struct {
	IpAddress        string `xml:"ipAddress"`
	PortNo           string `xml:"portNo"`
	MacAddress       string `xml:"macAddress"`
	ChannelId        string `xml:"channelID"`
	DateTime         string `xml:"dateTime"`
	ActivePostCount  string `xml:"activePostCount"`
	EventType        string `xml:"eventType"`
	EventState       string `xml:"eventState"`
	EventDescription string `xml:"eventDescription"`
}

func (evt *Event) MarshalBuffer() ([]byte, error) {
	return xml.Marshal(evt)
}

func UnmarshalBuffer(data []byte) (*Event, error) {
	var evt Event
	err := xml.Unmarshal(data, &evt)
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("xml.Unmarshal failed with '%s'\n", err)
		}
	}
	return &evt, nil
}
