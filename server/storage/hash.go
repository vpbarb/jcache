package storage

type Hash map[string]string

func (h Hash) GetValue(field string) (string, error) {
	if value, found := h[field]; found {
		return value, nil
	} else {
		return "", FieldNotExistError
	}
}
