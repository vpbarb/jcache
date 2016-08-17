package main

import (
	"github.com/Barberrrry/jcache/server"
)

func main() {
	s := server.New()
	s.ListenAndServe(":8081")
}
