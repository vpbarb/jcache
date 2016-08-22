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

	testValueType(client)
	testHashType(client)
	testListType(client)
}

func testValueType(client *client.Client) {
	log.Print("Value type operations")

	key := "key1"

	keys, err := client.Keys()
	log.Printf("Keys: %v", keys)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	err = client.Set(key, "value1", 3600)
	log.Printf("Set: %s", key)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	keys, err = client.Keys()
	log.Printf("Keys: %v", keys)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	value, err := client.Get(key)
	log.Printf("Get: %s = %s", key, value)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	err = client.Update(key, "new_value1")
	log.Printf("Set: %s = %s", key, "new_value1")
	if err != nil {
		log.Printf("Error: %s", err)
	}

	value, err = client.Get(key)
	log.Printf("Get: %s = %s", key, value)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	err = client.Delete(key)
	log.Printf("Delete: %s", key)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	value, err = client.Get(key)
	log.Printf("Get: %s = %s", key, value)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	keys, err = client.Keys()
	log.Printf("Keys: %v", keys)
	if err != nil {
		log.Printf("Error: %s", err)
	}
}

func testHashType(client *client.Client) {
	log.Print("Hash type operations")

	key := "hash1"
	field1 := "field1"
	field2 := "field2"

	keys, err := client.Keys()
	log.Printf("Keys: %v", keys)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	err = client.HashCreate(key, 3600)
	log.Printf("Create: %s", key)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	keys, err = client.Keys()
	log.Printf("Keys: %v", keys)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	err = client.HashSet(key, field1, "value1")
	log.Printf("Hash set: %s[%s]", key, field1)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	value, err := client.HashGet(key, field1)
	log.Printf("Hash get: %s[%s] = %s", key, field1, value)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	err = client.HashSet(key, field2, "value2")
	log.Printf("Hash set: %s[%s]", key, field2)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	values, err := client.HashGetAll(key)
	log.Printf("Hash get all: %s = %+v", key, values)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	keys, err = client.HashKeys(key)
	log.Printf("Hash keys: %s = %+v", key, keys)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	len, err := client.HashLength(key)
	log.Printf("Hash length: %s = %d", key, len)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	err = client.HashDelete(key, field1)
	log.Printf("Hash delete: %s[%s]", key, field1)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	value, err = client.HashGet(key, field1)
	log.Printf("Hash get: %s[%s] = %s", key, field1, value)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	err = client.Delete(key)
	log.Printf("Delete: %s", key)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	keys, err = client.Keys()
	log.Printf("Keys: %v", keys)
	if err != nil {
		log.Printf("Error: %s", err)
	}
}

func testListType(client *client.Client) {
	log.Print("List type operations")

	key := "list1"

	keys, err := client.Keys()
	log.Printf("Keys: %v", keys)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	err = client.ListCreate(key, 3600)
	log.Printf("Create: %s", key)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	keys, err = client.Keys()
	log.Printf("Keys: %v", keys)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	err = client.ListRightPush(key, "value1")
	log.Printf("List right push: %s", key)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	err = client.ListLeftPush(key, "value2")
	log.Printf("List left push: %s", key)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	err = client.ListRightPush(key, "value3")
	log.Printf("List right push: %s", key)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	values, err := client.ListRange(key, 0, 2)
	log.Printf("List range: %s = %+v", key, values)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	len, err := client.ListLength(key)
	log.Printf("List length: %s = %d", key, len)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	value, err := client.ListRightPop(key)
	log.Printf("List right pop: %s = %s", key, value)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	value, err = client.ListLeftPop(key)
	log.Printf("List right pop: %s = %s", key, value)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	err = client.Delete(key)
	log.Printf("Delete: %s", key)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	keys, err = client.Keys()
	log.Printf("Keys: %v", keys)
	if err != nil {
		log.Printf("Error: %s", err)
	}
}
