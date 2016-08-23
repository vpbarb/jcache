package protocol

import (
	"fmt"
	"io"
)

type Request interface {
	Encode() ([]byte, error)
	Decode(io.Reader) error
}

type Response interface {
	Encode() ([]byte, error)
	Decode(io.Reader) error
}

func ReadRequestCommand(r io.Reader) (string, error) {
	var command string
	_, err := fmt.Fscanf(r, "%s", &command)
	if err != nil {
		return "", err
	}
	return command, nil
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
	return newOkResponse()
}

func NewKeysResponse() *keysResponse {
	return &keysResponse{dataResponse: newDataResponse()}
}

func NewGetResponse() *valueResponse {
	return &valueResponse{dataResponse: newDataResponse()}
}

func NewSetResponse() *okResponse {
	return newOkResponse()
}

func NewDelResponse() *okResponse {
	return newOkResponse()
}

func NewUpdResponse() *okResponse {
	return newOkResponse()
}

func NewHashCreateResponse() *okResponse {
	return newOkResponse()
}

func NewHashGetResponse() *valueResponse {
	return &valueResponse{dataResponse: newDataResponse()}
}

func NewHashSetResponse() *okResponse {
	return newOkResponse()
}

func NewHashDelResponse() *okResponse {
	return newOkResponse()
}

func NewHashGetAllResponse() *fieldsResponse {
	return &fieldsResponse{dataResponse: newDataResponse()}
}

func NewHashKeysResponse() *keysResponse {
	return &keysResponse{dataResponse: newDataResponse()}
}

func NewHashLenResponse() *lenResponse {
	return &lenResponse{dataResponse: newDataResponse()}
}

func NewListCreateResponse() *okResponse {
	return newOkResponse()
}

func NewListRightPushResponse() *okResponse {
	return newOkResponse()
}

func NewListLeftPushResponse() *okResponse {
	return newOkResponse()
}

func NewListRightPopResponse() *valueResponse {
	return &valueResponse{dataResponse: newDataResponse()}
}

func NewListLeftPopResponse() *valueResponse {
	return &valueResponse{dataResponse: newDataResponse()}
}

func NewListLenResponse() *lenResponse {
	return &lenResponse{dataResponse: newDataResponse()}
}

func NewListRangeResponse() *valuesResponse {
	return &valuesResponse{dataResponse: newDataResponse()}
}

func NewErrorResponse(err error) *okResponse {
	return &okResponse{response: &response{Error: err}}
}
