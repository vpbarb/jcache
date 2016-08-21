package protocol

import (
	"io"
	"regexp"
)

const (
	keyTemplate = "[a-zA-Z0-9_]+"
	intTemplate = "[0-9_]+"
)

var (
	keyRegexp = regexp.MustCompile(keyTemplate)
	ttlRegexp = regexp.MustCompile(intTemplate)
	lenRegexp = regexp.MustCompile(intTemplate)
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
