build:
	@go build -ldflags "-w -s" -o ./bin/gohtools

run: build
	@./bin/gohtools

test:
	@go test -v ./...