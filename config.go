package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type SiaEvent struct {
	EventType  string `yaml:"eventType"`
	EventState string `yaml:"eventState"`
	SiaCode    int    `yaml:"siaCode"`
}

type Config struct {
	Account string     `yaml:"account"`
	Zone    string     `yaml:"zone"`
	Events  []SiaEvent `yaml:"events"`
}

func LoadConfig() *Config {
	fileName := "config.yaml"
	var c Config
	data, err := ioutil.ReadFile(fileName)
	failOnError(err, "error reading "+fileName)

	err = yaml.Unmarshal(data, &c)
	failOnError(err, "error unmarshalling "+fileName)
	return &c
}
