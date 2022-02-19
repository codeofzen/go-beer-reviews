package main

import (
	"errors"
	"log"

	"github.com/kelseyhightower/envconfig"
)

var (
	ErrPortMissing = errors.New("requires valid port to start server")
)


type Specification struct {
	Port	string	`envconfig:"SERVICE_PORT" required:"true"`
	Debug 	bool 	`envconfig:"DEBUG" default:"false"`
	
}


func run(config Specification) error {

	log.Printf("Debug mode: %t\n", config.Debug )
	log.Printf("Starting beer-review service on port: %s\n", config.Port)
	return nil
}


func main() {

	var s Specification

	
	
	if err := envconfig.Process("beer-review", &s); err != nil {
		log.Fatal(err.Error())
	}

	if err := run(s); err != nil {
		log.Fatalf("Starting beer-review service failed: %s", err)
	}
}