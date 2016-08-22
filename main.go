package main

import (
	"flag"
	"log"
	"time"

	"github.com/Barberrrry/jcache/server"
	"github.com/Barberrrry/jcache/server/storage"
	"github.com/Barberrrry/jcache/server/storage/boltdb"
	"github.com/Barberrrry/jcache/server/storage/memory"
	"github.com/Barberrrry/jcache/server/storage/multi"
)

func main() {
	storageType := server.StorageType(server.StorageMemory)

	htpasswdPath := flag.String("htpasswd", "", "Path to .htpasswd file for authentication. Leave blank to disable authentication.")
	listen := flag.String("listen", ":9999", "Host and port to listen connection")
	flag.Var(&storageType, "storage_type", "Type of storage (memory, multi_memory)")
	storageMemorySize := flag.Uint("storage_memory_size", 10000, "Max number of stored elements")
	storageMultiMemoryCount := flag.Uint("storage_multi_memory_count", 1, "Number of storages inside multi memory storage")
	storageBoltDbPath := flag.String("storage_boltdb_path", "", "Path to BoltDB file")
	storageGCInterval := flag.Duration("storage_gc_interval", time.Minute, "Storage GC interval")
	flag.Parse()

	var storage storage.Storage

	log.Printf(`storage type is "%s"`, storageType)

	switch storageType {
	case server.StorageMemory:
		var err error
		storage, err = memory.NewStorage(int(*storageMemorySize), *storageGCInterval)
		if err != nil {
			log.Fatalln(err)
		}
	case server.StorageMultiMemory:
		ms := multi.NewStorage()
		for i := uint(0); i < *storageMultiMemoryCount; i++ {
			s, err := memory.NewStorage(int(*storageMemorySize), *storageGCInterval)
			if err != nil {
				log.Fatalln(err)
			}
			ms.AddStorage(s)
		}
		storage = ms

	case server.StorageBoltDB:
		var err error
		storage, err = boltdb.NewStorage(*storageBoltDbPath, *storageGCInterval)
		if err != nil {
			log.Fatalln(err)
		}
	}

	s := server.New(storage, *htpasswdPath)
	s.ListenAndServe(*listen)
}
