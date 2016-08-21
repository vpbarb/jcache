package protocol

import (
	"io"
)

type Request interface {
	Encode() ([]byte, error)
	Decode([]byte, io.Reader) error
	Command() string
}

type Response interface {
	Encode() ([]byte, error)
	Decode([]byte, io.Reader) error
	Error() error
}

func NewGetRequest(key string) *keyRequest {
	return newKeyRequest("GET", key)
}

func NewSetRequest(key, value string, ttl uint64) *setRequest {
	return &setRequest{
		keyRequest: newKeyRequest("SET", key),
		Value:      value,
		TTL:        ttl,
	}
}

func NewDelRequest(key string) *keyRequest {
	return newKeyRequest("DEL", key)
}

func NewEmptyResponse(err error) *response {
	return &response{err: err}
}

func NewValueResponse(value string, err error) *valueResponse {
	return &valueResponse{response: response{err: err}, Value: value}
}
