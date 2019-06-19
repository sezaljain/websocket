runserver: build
	go run main.go

build: install
	go build main.go
	go build client.go


install:
	go get github.com/gorilla/websocket
	go get github.com/gorilla/mux
