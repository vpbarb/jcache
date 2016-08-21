package protocol

import (
	"io"
)

type Encoder interface {
	Encode() ([]byte, error)
}

type Decoder interface {
	Decode([]byte, io.Reader) error
}

func NewAuthRequest(user, password string) *authRequest {
	return &authRequest{
		request:  newRequest("AUTH"),
		User:     user,
		Password: password,
	}
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
	return &response{Error: err}
}

func NewValueResponse(value string, err error) *valueResponse {
	return &valueResponse{response: response{Error: err}, Value: value}
}
