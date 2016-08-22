package server

import (
	"fmt"
	"io"

	"github.com/Barberrrry/jcache/protocol"
	"github.com/Barberrrry/jcache/server/htpasswd"
	"github.com/Barberrrry/jcache/server/storage"
)

type command func(header []byte, data io.Reader) []byte

func formatError(err error) []byte {
	response := protocol.NewErrorResponse(err)
	data, _ := response.Encode()
	return data
}

func run(header []byte, data io.Reader, request protocol.Decoder, response protocol.Encoder, action func()) []byte {
	err := request.Decode(header, data)
	if err != nil {
		return formatError(err)
	}

	action()

	result, err := response.Encode()
	if err != nil {
		return formatError(err)
	}

	return result
}

func newKeysCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewKeysRequest()
		response := protocol.NewKeysResponse()
		return run(header, data, request, response, func() {
			response.Keys = storage.Keys()
		})
	}
}

func newGetCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewGetRequest()
		response := protocol.NewGetResponse()
		return run(header, data, request, response, func() {
			response.Value, response.Error = storage.Get(request.Key)
		})
	}
}

func newSetCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewSetRequest()
		response := protocol.NewSetResponse()
		return run(header, data, request, response, func() {
			response.Error = storage.Set(request.Key, request.Value, request.TTL)
		})
	}
}

func newUpdCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewUpdRequest()
		response := protocol.NewUpdResponse()
		return run(header, data, request, response, func() {
			response.Error = storage.Update(request.Key, request.Value)
		})
	}
}

func newDelCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewDelRequest()
		response := protocol.NewDelResponse()
		return run(header, data, request, response, func() {
			response.Error = storage.Delete(request.Key)
		})
	}
}

func newHashCreateCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewHashCreateRequest()
		response := protocol.NewHashCreateResponse()
		return run(header, data, request, response, func() {
			response.Error = storage.HashCreate(request.Key, request.TTL)
		})
	}
}

func newHashGetAllCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewHashGetAllRequest()
		response := protocol.NewHashGetAllResponse()
		return run(header, data, request, response, func() {
			response.Fields, response.Error = storage.HashGetAll(request.Key)
		})
	}
}

func newHashGetCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewHashGetRequest()
		response := protocol.NewHashGetResponse()
		return run(header, data, request, response, func() {
			response.Value, response.Error = storage.HashGet(request.Key, request.Field)
		})
	}
}

func newHashSetCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewHashSetRequest()
		response := protocol.NewHashSetResponse()
		return run(header, data, request, response, func() {
			response.Error = storage.HashSet(request.Key, request.Field, request.Value)
		})
	}
}

func newHashDelCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewHashDelRequest()
		response := protocol.NewHashDelResponse()
		return run(header, data, request, response, func() {
			response.Error = storage.HashDelete(request.Key, request.Field)
		})
	}
}

func newHashLenCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewHashLenRequest()
		response := protocol.NewHashLenResponse()
		return run(header, data, request, response, func() {
			response.Len, response.Error = storage.HashLen(request.Key)
		})
	}
}

func newHashKeysCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewHashKeysRequest()
		response := protocol.NewHashKeysResponse()
		return run(header, data, request, response, func() {
			response.Keys, response.Error = storage.HashKeys(request.Key)
		})
	}
}

func newListCreateCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewListCreateRequest()
		response := protocol.NewListCreateResponse()
		return run(header, data, request, response, func() {
			response.Error = storage.ListCreate(request.Key, request.TTL)
		})
	}
}

func newListLeftPopCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewListLeftPopRequest()
		response := protocol.NewListLeftPopResponse()
		return run(header, data, request, response, func() {
			response.Value, response.Error = storage.ListLeftPop(request.Key)
		})
	}
}

func newListRightPopCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewListRightPopRequest()
		response := protocol.NewListRightPopResponse()
		return run(header, data, request, response, func() {
			response.Value, response.Error = storage.ListRightPop(request.Key)
		})
	}
}

func newListLeftPushCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewListLeftPushRequest()
		response := protocol.NewListLeftPushResponse()
		return run(header, data, request, response, func() {
			response.Error = storage.ListLeftPush(request.Key, request.Value)
		})
	}
}

func newListRightPushCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewListRightPushRequest()
		response := protocol.NewListRightPushResponse()
		return run(header, data, request, response, func() {
			response.Error = storage.ListRightPush(request.Key, request.Value)
		})
	}
}

func newListLenCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewListLenRequest()
		response := protocol.NewListLenResponse()
		return run(header, data, request, response, func() {
			response.Len, response.Error = storage.ListLen(request.Key)
		})
	}
}

func newListRangeCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewListRangeRequest()
		response := protocol.NewListRangeResponse()
		return run(header, data, request, response, func() {
			response.Values, response.Error = storage.ListRange(request.Key, request.Start, request.Stop)
		})
	}
}

func newAuthCommand(htpasswdFile *htpasswd.HtpasswdFile, session *session) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewAuthRequest()
		response := protocol.NewAuthResponse()
		return run(header, data, request, response, func() {
			if htpasswdFile == nil || htpasswdFile.Validate(request.User, request.Password) {
				session.authorize()
			} else {
				response.Error = fmt.Errorf("Invalid credentials")
			}
		})
	}
}
