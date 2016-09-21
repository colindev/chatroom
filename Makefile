all: build

test:
	go test -v ./server

build: test
	go build -o chatroom ./server

