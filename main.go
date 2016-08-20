package main

import (
	"flag"

	"github.com/Barberrrry/jcache/server"
)

func main() {
	htpasswdPath := flag.String("htpasswd", "", "Path to .htpasswd file for authentication. Leave blank to disable authentication.")
	listen := flag.String("listen", ":9999", "Host and port to listen connection")
	flag.Parse()

	s := server.New(*htpasswdPath)
	s.ListenAndServe(*listen)
}
