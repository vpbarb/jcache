package main

import (
	"flag"
	"log"

	"github.com/Barberrrry/jcache/server"
	"github.com/Barberrrry/jcache/server/storage/boltdb"
	"github.com/Barberrrry/jcache/server/storage/memory"
)

func main() {
	storageType := server.StorageType(server.StorageMemory)

	htpasswdPath := flag.String("htpasswd", "", "Path to .htpasswd file for authentication. Leave blank to disable authentication.")
	listen := flag.String("listen", ":9999", "Host and port to listen connection")
	flag.Var(&storageType, "storage_type", "Type of storage (memory, multi_memory)")
	storageMultiMemoryCount := flag.Uint("storage_multi_memory_count", 1, "Number of storages inside multi memory storage")
	storageBoltDbPath := flag.String("storage_boltdb_path", "", "Path to BoltDB file")
	flag.Parse()

	var storage server.Storage

	log.Printf(`storage type is "%s"`, storageType)

	switch storageType {
	case server.StorageMemory:
		storage = memory.NewStorage()
	case server.StorageMultiMemory:
		storage = memory.NewMultiStorage(*storageMultiMemoryCount)
	case server.StorageBoltDB:
		var err error
		storage, err = boltdb.NewStorage(*storageBoltDbPath)
		if err != nil {
			log.Fatalln(err)
		}
	}

	s := server.New(storage, *htpasswdPath)
	s.ListenAndServe(*listen)
}
