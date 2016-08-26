all: build

build: deps
	GO15VENDOREXPERIMENT=1 go build -o jcache main.go

test: deps
	GO15VENDOREXPERIMENT=1 go test -v ./server/... -check.v

bench: deps
	GO15VENDOREXPERIMENT=1 go test -v ./server/... -check.b -check.bmem

deps:
	go get github.com/Masterminds/glide
	glide install
