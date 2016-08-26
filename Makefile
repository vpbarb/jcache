all: build

build: deps
	go build -o jcache main.go

test: deps
	go test ./server/... ./protocol

bench: deps
	go test -v ./server/... ./protocol -check.b -check.bmem

deps:
	go get github.com/Masterminds/glide
	glide install
