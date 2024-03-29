.PHONY: build
build:
	go build -trimpath -ldflags="-w -s" -o bin/server ./cmd/server
	go build -trimpath -ldflags="-w -s" -o bin/client ./cmd/client

.PHONY: run
run:
	go run ./cmd/server

