all: build-server build-client

test-client:
	go test -v ./clientcli

build-client: test-client
	go build -o chatroom-cli ./clientcli

test-server:
	go test -v ./server

build-server: test-server
	go build -o chatroom-server ./server

