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
	response := protocol.NewOkResponse()
	response.Error = err
	data, _ := response.Encode()
	return data
}

func newKeysCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewKeysRequest()
		response := protocol.NewKeysResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Keys = storage.Keys()

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newGetCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewGetRequest()
		response := protocol.NewValueResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Value, response.Error = storage.Get(request.Key)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newSetCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewSetRequest()
		response := protocol.NewOkResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Error = storage.Set(request.Key, request.Value, request.TTL)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newUpdCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewUpdRequest()
		response := protocol.NewOkResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Error = storage.Update(request.Key, request.Value)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newDelCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewDelRequest()
		response := protocol.NewOkResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Error = storage.Delete(request.Key)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newHashCreateCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewHashCreateRequest()
		response := protocol.NewOkResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Error = storage.HashCreate(request.Key, request.TTL)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newHashGetAllCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewHashGetAllRequest()
		response := protocol.NewFieldsResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Fields, response.Error = storage.HashGetAll(request.Key)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newHashGetCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewHashGetRequest()
		response := protocol.NewValueResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Value, response.Error = storage.HashGet(request.Key, request.Field)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newHashSetCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewHashSetRequest()
		response := protocol.NewOkResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Error = storage.HashSet(request.Key, request.Field, request.Value)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newHashDelCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewHashDelRequest()
		response := protocol.NewOkResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Error = storage.HashDelete(request.Key, request.Field)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newHashLenCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewHashLenRequest()
		response := protocol.NewLenResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Len, response.Error = storage.HashLen(request.Key)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newHashKeysCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewHashKeysRequest()
		response := protocol.NewKeysResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Keys, response.Error = storage.HashKeys(request.Key)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newListCreateCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewListCreateRequest()
		response := protocol.NewOkResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Error = storage.ListCreate(request.Key, request.TTL)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newListLeftPopCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewListLeftPopRequest()
		response := protocol.NewValueResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Value, response.Error = storage.ListLeftPop(request.Key)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newListRightPopCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewListRightPopRequest()
		response := protocol.NewValueResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Value, response.Error = storage.ListRightPop(request.Key)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newListLeftPushCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewListLeftPushRequest()
		response := protocol.NewOkResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Error = storage.ListLeftPush(request.Key, request.Value)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newListRightPushCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewListRightPushRequest()
		response := protocol.NewOkResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Error = storage.ListRightPush(request.Key, request.Value)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newListLenCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewListLenRequest()
		response := protocol.NewLenResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Len, response.Error = storage.ListLen(request.Key)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newListRangeCommand(storage storage.Storage) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewListRangeRequest()
		response := protocol.NewValuesResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		response.Values, response.Error = storage.ListRange(request.Key, request.Start, request.Stop)

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}

func newAuthCommand(htpasswdFile *htpasswd.HtpasswdFile, session *session) command {
	return func(header []byte, data io.Reader) []byte {
		request := protocol.NewAuthRequest()
		response := protocol.NewOkResponse()

		err := request.Decode(header, data)
		if err != nil {
			return formatError(err)
		}

		if htpasswdFile == nil || htpasswdFile.Validate(request.User, request.Password) {
			session.authorize()
		} else {
			response.Error = fmt.Errorf("Invalid credentials")
		}

		result, err := response.Encode()
		if err != nil {
			return formatError(err)
		}

		return result
	}
}
