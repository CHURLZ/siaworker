package main

import (
	"io/ioutil"
	"log"

	"github.com/davecgh/go-spew/spew"
	yaml "gopkg.in/yaml.v2"
)

// Events in the .yaml are events to be forwarded, and what SIA code to apply
type Config struct {
	Account string     `yaml:"account"`
	Zone    string     `yaml:"zone"`
	Events  []SiaEvent `yaml:"events"`
}

type SiaEvent struct {
	EventType  string `yaml:"eventType"`
	EventState string `yaml:"eventState"`
	SiaCode    string `yaml:"siaCode"`
}

func LoadConfig() *Config {
	fileName := "config.yaml"
	var c Config
	data, err := ioutil.ReadFile(fileName)
	failOnError(err, "error reading "+fileName)

	err = yaml.Unmarshal(data, &c)
	failOnError(err, "error unmarshalling "+fileName)

	log.Println("Config loaded: " + fileName)
	spew.Dump(c)
	return &c
}
