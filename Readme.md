# go-cache
==========

## Requirements to run
* go (tested with version go1.14.6 linux/amd64)
* make (optional)

## How to build
* With make, `make build`
* Without make, `go build cmd/server/main.go`

## How to run
* with go cli, `go run cmd/server/main.go`
* Wihout go cli
```bash
go build cmd/server/main.go
./main
```
* run with profiler web server
```bash
go build cmd/server/main.go
./main -dev
```

## How to run tests
* With make, `make test`
* With coverage, `make test test-coverage`

