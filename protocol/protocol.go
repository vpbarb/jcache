package protocol

import (
	"bufio"
	"io"
	"strings"
)

type Request interface {
	Command() string
}

type Encoder interface {
	Encode() ([]byte, error)
}

type Decoder interface {
	Decode([]byte, io.Reader) error
}

func ReadRequestHeader(r io.Reader) ([]byte, string, error) {
	rb := bufio.NewReader(r)
	line, _, err := rb.ReadLine()
	if err != nil {
		return nil, "", err
	}
	parts := strings.SplitN(string(line), " ", 2)
	return line, parts[0], nil
}

// Requests

func NewAuthRequest() *authRequest {
	return &authRequest{request: newRequest("AUTH")}
}

func NewKeysRequest() *request {
	r := newRequest("KEYS")
	return &r
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

func NewHashCreateRequest() *keyTTLRequest {
	return newKeyTTLRequest("HCREATE")
}

func NewHashGetRequest() *keyFieldRequest {
	return newKeyFieldRequest("HGET")
}

func NewHashSetRequest() *keyFieldValueRequest {
	return newKeyFieldValueRequest("HSET")
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
	return newKeyTTLRequest("LCREATE")
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

func NewAuthResponse() *okResponse {
	return &okResponse{}
}

func NewKeysResponse() *keysResponse {
	return &keysResponse{}
}

func NewGetResponse() *valueResponse {
	return &valueResponse{}
}

func NewSetResponse() *okResponse {
	return &okResponse{}
}

func NewDelResponse() *okResponse {
	return &okResponse{}
}

func NewUpdResponse() *okResponse {
	return &okResponse{}
}

func NewHashCreateResponse() *okResponse {
	return &okResponse{}
}

func NewHashGetResponse() *valueResponse {
	return &valueResponse{}
}

func NewHashSetResponse() *okResponse {
	return &okResponse{}
}

func NewHashDelResponse() *okResponse {
	return &okResponse{}
}

func NewHashGetAllResponse() *fieldsResponse {
	return &fieldsResponse{}
}

func NewHashKeysResponse() *keysResponse {
	return &keysResponse{}
}

func NewHashLenResponse() *lenResponse {
	return &lenResponse{}
}

func NewListCreateResponse() *okResponse {
	return &okResponse{}
}

func NewListRightPushResponse() *okResponse {
	return &okResponse{}
}

func NewListLeftPushResponse() *okResponse {
	return &okResponse{}
}

func NewListRightPopResponse() *valueResponse {
	return &valueResponse{}
}

func NewListLeftPopResponse() *valueResponse {
	return &valueResponse{}
}

func NewListLenResponse() *lenResponse {
	return &lenResponse{}
}

func NewListRangeResponse() *valuesResponse {
	return &valuesResponse{}
}

func NewErrorResponse(err error) *okResponse {
	return &okResponse{response: response{Error: err}}
}
