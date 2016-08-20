package storage

import (
	"fmt"
)

type Hash map[string]string

func (h Hash) GetValue(field string) (string, error) {
	if value, found := h[field]; found {
		return value, nil
	} else {
		return "", fmt.Errorf(`Field "%s" does not exist`, field)
	}
}
