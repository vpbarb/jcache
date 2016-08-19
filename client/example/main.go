package main

import (
	"log"
	"sync"
	"time"

	"github.com/Barberrrry/jcache/client"
)

func main() {
	client, err := client.NewClient("127.0.0.1:9999", "admin", "admin", 10*time.Second, 5)

	if err != nil {
		log.Fatalf("Client creation error: %s", err)
	}

	for {
		value, err := client.Get("test1")
		log.Printf("Value: %s", value)
		log.Printf("Error: %s", err)

		time.Sleep(time.Second)
	}

	//value, err = client.Get("test2")
	//log.Printf("Value: %s", value)
	//log.Printf("Error: %s", err)
	//
	//value, err = client.Get("test3")
	//log.Printf("Value: %s", value)
	//log.Printf("Error: %s", err)

	wg := sync.WaitGroup{}
	//wg.Add(3)
	//
	//go func() {
	//	value, err := client.Get("test1")
	//	log.Printf("Value: %s", value)
	//	log.Printf("Error: %s", err)
	//	wg.Done()
	//}()
	//
	//go func() {
	//	value, err := client.Get("test2")
	//	log.Printf("Value: %s", value)
	//	log.Printf("Error: %s", err)
	//	wg.Done()
	//}()
	//
	//go func() {
	//	value, err := client.Get("test3")
	//	log.Printf("Value: %s", value)
	//	log.Printf("Error: %s", err)
	//	wg.Done()
	//}()

	wg.Wait()
}
