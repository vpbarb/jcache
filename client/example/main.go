package main

import (
	"log"

	"github.com/Barberrrry/jcache/client"
)

func main() {
	client, err := client.NewClient("127.0.0.1:9999", "admin", "admin", 1)

	if err != nil {
		log.Printf("Client creation error: %s", err)
	}

	value, err := client.Get("test")
	if err != nil {
		log.Printf("Error: %s", err)
	}

	log.Printf("Value: %s", value)
}
