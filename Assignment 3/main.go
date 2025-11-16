package main

import (
	"log"
)

func main() {
	err := initMongoConnections()
	if err != nil {
		log.Fatalf("failed mongo connection: %v", err)
	}

	startCli()

	err = closeMongoConnections()
	if err != nil {
		log.Fatalf("failed mongo disconnect: %v", err)
	}
}
