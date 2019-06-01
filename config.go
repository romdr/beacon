package main

import (
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

// Config defines the behavior of the beacon
type Config struct {
	Interval time.Duration
	Targets  []struct {
		Type string
		Arg  string
	}
}

// Load the yaml config file
func (config *Config) load(path string) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("ERROR: Reading config file: %s", err)
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Fatalf("ERROR: Parsing configuration: %s", err)
	}
}

// Print the config
func (config *Config) print() {
	log.Printf("%+v\n", *config)
}
