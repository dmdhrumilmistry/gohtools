build:
	@go build -o ./bin/gohtools

run: build
	@./bin/gohtools

test:
	@go test -v ./...