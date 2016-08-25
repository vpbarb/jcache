package server

import (
	"errors"
	"io"

	"github.com/Barberrrry/jcache/protocol"
	"github.com/Barberrrry/jcache/server/htpasswd"
	"github.com/Barberrrry/jcache/server/storage"
)

type command func(io.ReadWriter)

var (
	invalidCredentialsError = errors.New("Invalid credentials")
)

func writeError(writer io.Writer, err error) {
	response := protocol.NewErrorResponse(err)
	response.Encode(writer)
}

func run(rw io.ReadWriter, request protocol.Request, response protocol.Response, action func()) {
	err := request.Decode(rw)
	if err != nil {
		writeError(rw, err)
	}

	action()

	err = response.Encode(rw)
	if err != nil {
		writeError(rw, err)
	}
}

func newKeysCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewKeysRequest()
		response := protocol.NewKeysResponse()
		run(rw, request, response, func() {
			response.Keys = storage.Keys()
		})
	}
}

func newGetCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewGetRequest()
		response := protocol.NewGetResponse()
		run(rw, request, response, func() {
			response.Value, response.Error = storage.Get(request.Key)
		})
	}
}

func newSetCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewSetRequest()
		response := protocol.NewSetResponse()
		run(rw, request, response, func() {
			response.Error = storage.Set(request.Key, request.Value, request.TTL)
		})
	}
}

func newUpdCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewUpdRequest()
		response := protocol.NewUpdResponse()
		run(rw, request, response, func() {
			response.Error = storage.Update(request.Key, request.Value)
		})
	}
}

func newDelCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewDelRequest()
		response := protocol.NewDelResponse()
		run(rw, request, response, func() {
			response.Error = storage.Delete(request.Key)
		})
	}
}

func newHashCreateCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewHashCreateRequest()
		response := protocol.NewHashCreateResponse()
		run(rw, request, response, func() {
			response.Error = storage.HashCreate(request.Key, request.TTL)
		})
	}
}

func newHashGetAllCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewHashGetAllRequest()
		response := protocol.NewHashGetAllResponse()
		run(rw, request, response, func() {
			response.Fields, response.Error = storage.HashGetAll(request.Key)
		})
	}
}

func newHashGetCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewHashGetRequest()
		response := protocol.NewHashGetResponse()
		run(rw, request, response, func() {
			response.Value, response.Error = storage.HashGet(request.Key, request.Field)
		})
	}
}

func newHashSetCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewHashSetRequest()
		response := protocol.NewHashSetResponse()
		run(rw, request, response, func() {
			response.Error = storage.HashSet(request.Key, request.Field, request.Value)
		})
	}
}

func newHashDelCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewHashDelRequest()
		response := protocol.NewHashDelResponse()
		run(rw, request, response, func() {
			response.Error = storage.HashDelete(request.Key, request.Field)
		})
	}
}

func newHashLenCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewHashLenRequest()
		response := protocol.NewHashLenResponse()
		run(rw, request, response, func() {
			response.Len, response.Error = storage.HashLen(request.Key)
		})
	}
}

func newHashKeysCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewHashKeysRequest()
		response := protocol.NewHashKeysResponse()
		run(rw, request, response, func() {
			response.Keys, response.Error = storage.HashKeys(request.Key)
		})
	}
}

func newListCreateCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewListCreateRequest()
		response := protocol.NewListCreateResponse()
		run(rw, request, response, func() {
			response.Error = storage.ListCreate(request.Key, request.TTL)
		})
	}
}

func newListLeftPopCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewListLeftPopRequest()
		response := protocol.NewListLeftPopResponse()
		run(rw, request, response, func() {
			response.Value, response.Error = storage.ListLeftPop(request.Key)
		})
	}
}

func newListRightPopCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewListRightPopRequest()
		response := protocol.NewListRightPopResponse()
		run(rw, request, response, func() {
			response.Value, response.Error = storage.ListRightPop(request.Key)
		})
	}
}

func newListLeftPushCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewListLeftPushRequest()
		response := protocol.NewListLeftPushResponse()
		run(rw, request, response, func() {
			response.Error = storage.ListLeftPush(request.Key, request.Value)
		})
	}
}

func newListRightPushCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewListRightPushRequest()
		response := protocol.NewListRightPushResponse()
		run(rw, request, response, func() {
			response.Error = storage.ListRightPush(request.Key, request.Value)
		})
	}
}

func newListLenCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewListLenRequest()
		response := protocol.NewListLenResponse()
		run(rw, request, response, func() {
			response.Len, response.Error = storage.ListLen(request.Key)
		})
	}
}

func newListRangeCommand(storage storage.Storage) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewListRangeRequest()
		response := protocol.NewListRangeResponse()
		run(rw, request, response, func() {
			response.Values, response.Error = storage.ListRange(request.Key, request.Start, request.Stop)
		})
	}
}

func newAuthCommand(htpasswdFile *htpasswd.HtpasswdFile, session *session) command {
	return func(rw io.ReadWriter) {
		request := protocol.NewAuthRequest()
		response := protocol.NewAuthResponse()
		run(rw, request, response, func() {
			if htpasswdFile == nil || htpasswdFile.Validate(request.User, request.Password) {
				session.authorize()
			} else {
				response.Error = invalidCredentialsError
			}
		})
	}
}
