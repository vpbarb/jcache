package main

import (
	"flag"

	"github.com/Barberrrry/jcache/server"
)

func main() {
	htpasswdPath := flag.String("htpasswd", "", "Path to .htpasswd file for authentication. Leave blank to disable authentication.")
	host := flag.String("host", "", "Host to listen connection")
	port := flag.String("port", "8081", "Port to listen connection")
	flag.Parse()

	s := server.New(*htpasswdPath)
	s.ListenAndServe(*host + ":" + *port)
}
