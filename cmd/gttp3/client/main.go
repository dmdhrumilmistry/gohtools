package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func main() {
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}, // set a TLS client config, if desired
		QUICConfig: &quic.Config{}, // QUIC connection options
	}
	defer roundTripper.Close()
	client := &http.Client{
		Transport: roundTripper,
	}

	req, err := http.NewRequest(http.MethodGet, "https://localhost:8443", nil)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))
}
