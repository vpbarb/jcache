package server

import (
	"fmt"
)

const (
	StorageMemory      = "memory"
	StorageMultiMemory = "multi_memory"
	StorageBolt        = "bolt"
)

type StorageType string

func (s *StorageType) String() string {
	return string(*s)
}

func (s *StorageType) Set(value string) error {
	switch value {
	case StorageMemory:
		*s = StorageType(value)
	case StorageMultiMemory:
		*s = StorageType(value)
	case StorageBolt:
		*s = StorageType(value)
	default:
		return fmt.Errorf("Unknown storage type: %s", value)
	}
	return nil
}
