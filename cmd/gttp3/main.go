package main

import "github.com/dmdhrumilmistry/gohtools/http3"

func main() {
	// mkdir -p /tmp/certs
	// cd /tmp/certs
	// openssl genrsa -out ca.key 2048
	// openssl req -new -x509 -nodes -key ca.key -sha256 -days 365 -out ca.crt
	srv := http3.NewServer("0.0.0.0:8443", "/tmp/certs/ca.crt", "/tmp/certs/ca.key")
	srv.Listen()
}
