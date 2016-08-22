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

// Requests

func NewAuthRequest() *authRequest {
	return &authRequest{request: newRequest("AUTH")}
}

func NewGetRequest() *keyRequest {
	return newKeyRequest("GET")
}

func NewSetRequest() *setRequest {
	return &setRequest{keyValueRequest: newKeyValueRequest("SET")}
}

func NewDelRequest() *keyRequest {
	return newKeyRequest("DEL")
}

func NewUpdRequest() *keyValueRequest {
	return newKeyValueRequest("UPD")
}

func NewKeysRequest() *request {
	r := newRequest("KEYS")
	return &r
}

func NewHashCreateRequest() *keyTTLRequest {
	return &keyTTLRequest{keyRequest: newKeyRequest("HCREATE")}
}

func NewHashGetRequest() *keyFieldRequest {
	return newKeyFieldRequest("HGET")
}

func NewHashSetRequest() *keyFieldValueRequest {
	return &keyFieldValueRequest{keyFieldRequest: newKeyFieldRequest("HSET")}
}

func NewHashDelRequest() *keyFieldRequest {
	return newKeyFieldRequest("HDEL")
}

func NewHashKeysRequest() *keyRequest {
	return newKeyRequest("HKEYS")
}

func NewHashGetAllRequest() *keyRequest {
	return newKeyRequest("HGETALL")
}

func NewHashLenRequest() *keyRequest {
	return newKeyRequest("HLEN")
}

func NewListCreateRequest() *keyTTLRequest {
	return &keyTTLRequest{keyRequest: newKeyRequest("LCREATE")}
}

func NewListLenRequest() *keyRequest {
	return newKeyRequest("LLEN")
}

func NewListLeftPopRequest() *keyRequest {
	return newKeyRequest("LLPOP")
}

func NewListRightPopRequest() *keyRequest {
	return newKeyRequest("LRPOP")
}

func NewListLeftPushRequest() *keyValueRequest {
	return newKeyValueRequest("LLPUSH")
}

func NewListRightPushRequest() *keyValueRequest {
	return newKeyValueRequest("LRPUSH")
}

func NewListRangeRequest() *listRangeRequest {
	return &listRangeRequest{keyRequest: newKeyRequest("LRANGE")}
}

// Responses

func NewOkResponse() *okResponse {
	return &okResponse{}
}

func NewValueResponse() *valueResponse {
	return &valueResponse{}
}

func NewValuesResponse() *valuesResponse {
	return &valuesResponse{}
}

func NewKeysResponse() *keysResponse {
	return &keysResponse{}
}

func NewLenResponse() *lenResponse {
	return &lenResponse{}
}

func NewFieldsResponse() *fieldsResponse {
	return &fieldsResponse{}
}
