package main

import (
	"log"
	"time"

	"github.com/Barberrrry/jcache/client"
)

func main() {
	client, err := client.New("127.0.0.1:9999", "admin", "admin", 5*time.Second, 5)

	if err != nil {
		log.Fatalf("Client creation error: %s", err)
	}

	for {
		//value, err := client.Get("test1")
		keys, err := client.Keys()
		log.Printf("Value: %v", keys)
		log.Printf("Error: %s", err)

		time.Sleep(time.Second)
	}
}
