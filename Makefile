
clean:
	rm server

run: server
	./server

run-dev:
	go run cmd/server/main.go -- -dev

build:
	go build -o server cmd/server/main.go

server:
	go build -o server cmd/server/main.go

test:
	go test -coverprofile=coverage.out ./...

test-coverage: coverage.out
	go tool cover -html=coverage.out

lint:
	golangci-lint run
