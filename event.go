package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"github.com/satori/go.uuid"
)

type Event struct {
	EventId       string
	Zone          string
	AccountNumber string
	SiaCode       string

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

func (evt *Event) Prime(c Config) {
	evt.AccountNumber = c.Account
	evt.Zone = c.Zone
	sia, _ := evt.GetSia(c)
	evt.SiaCode = sia
	evt.EventId = uuid.Must(uuid.NewV4()).String()
}

func (evt *Event) GetSia(c Config) (string, error) {
	for _, e := range c.Events {
		if evt.EventType == e.EventType && evt.EventState == e.EventState {
			return e.SiaCode, nil
		}
	}
	return "", errors.New(fmt.Sprintf("Found no SiaCode for %s : %s", evt.EventType, evt.EventState))
}

func (evt *Event) QualifiesForPublish(c Config) bool {
	for _, e := range c.Events {
		if evt.EventType == e.EventType && evt.EventState == e.EventState {
			return true
		}
	}
	return false
}

func (evt *Event) Marshal() ([]byte, error) {
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
