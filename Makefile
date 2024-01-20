# build:
# 	@go build -ldflags "-w -s" -o ./bin/gohtools -v cmd/gohtools/main.go

install: 
	@go install -v ./...

run: test build
	@./bin/gohtools

test:
	@go test -v ./...