package memory

import (
	"fmt"
)

type hash map[string]string

func (h hash) getValue(field string) (string, error) {
	if value, found := h[field]; found {
		return value, nil
	} else {
		return "", fmt.Errorf(`Field "%s" does not exist`, field)
	}
}
