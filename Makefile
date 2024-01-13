build:
	@go build -ldflags "-w -s" -o ./bin/gohtools -v cmd/main.go

run: test build
	@./bin/gohtools

test:
	@go test -v ./...