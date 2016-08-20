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

	keys, err := client.Keys()
	log.Printf("Keys: %v", keys)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	err = client.Set("key1", "value1", time.Hour)
	log.Printf("Set: %s = %s", "key1", "value1")
	if err != nil {
		log.Printf("Error: %s", err)
	}

	keys, err = client.Keys()
	log.Printf("Keys: %v", keys)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	value, err := client.Get("key1")
	log.Printf("Get: %s = %s", "key1", value)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	ttl, err := client.TTL("key1")
	log.Printf("TTL: %s = %s", "key1", ttl)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	err = client.Update("key1", "new_value1")
	log.Printf("Set: %s = %s", "key1", "new_value1")
	if err != nil {
		log.Printf("Error: %s", err)
	}

	value, err = client.Get("key1")
	log.Printf("Get: %s = %s", "key1", value)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	err = client.Delete("key1")
	log.Printf("Delete: %s", "key1")
	if err != nil {
		log.Printf("Error: %s", err)
	}

	keys, err = client.Keys()
	log.Printf("Keys: %v", keys)
	if err != nil {
		log.Printf("Error: %s", err)
	}

}
