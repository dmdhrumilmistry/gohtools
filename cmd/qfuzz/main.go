package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	fhc "github.com/dmdhrumilmistry/fasthttpclient/client"
	"github.com/valyala/fasthttp"
)

type intSlice []int

func (s *intSlice) String() string {
	return fmt.Sprintf("%v", *s)
}
func (s *intSlice) Set(value string) error {
	for _, item := range strings.Split(value, ",") {
		intValue, err := strconv.Atoi(item)
		if err != nil {
			return err
		}
		*s = append(*s, intValue)
	}

	return nil
}

func matchesStatusCode(targetStatusCode int, statusCodes intSlice) bool {
	for _, statusCode := range statusCodes {
		if targetStatusCode == statusCode {
			return true
		}
	}
	return false
}

func main() {
	// parse args
	var skipTlsVerification bool
	var requestsPerSecond, max, min int
	var targetUrl string
	var matchCodes intSlice

	flag.BoolVar(&skipTlsVerification, "k", false, "Skips TLS verification if provided")
	flag.IntVar(&requestsPerSecond, "r", 100, "Requests per second")
	flag.IntVar(&max, "max", 100, "Max id to fuzz")
	flag.IntVar(&min, "min", 0, "Min id to fuzz")
	flag.StringVar(&targetUrl, "target", "", "target url. example: https://example.com")
	flag.Var(&matchCodes, "mc", "match codes to be printed. example: 200,301,401,403")

	flag.Parse()

	// validate args
	if targetUrl == "" {
		flag.PrintDefaults()
		return
	}

	if len(matchCodes) == 0 {
		matchCodes = append(matchCodes, 200, 201, 400, 401, 403)
	}

	// create fasthttp client with custom configurations
	fhclient := &fasthttp.Client{
		Name:               "qfuzz",
		ReadTimeout:        5 * time.Second,
		WriteTimeout:       5 * time.Second,
		MaxConnWaitTimeout: 60 * time.Second,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: skipTlsVerification,
		},
	}
	// create default header values
	headers := map[string]string{
		"User-Agent": "hfuzz",
	}

	// create rate limited client
	client := fhc.NewRateLimitedClient(requestsPerSecond, 1, fhclient)

	// create requests
	requests := make([]*fhc.Request, 0)
	for i := min; i <= max; i++ {
		queryParams := map[string]string{
			"id": strconv.Itoa(i),
		}
		requests = append(
			requests,
			fhc.NewRequest(
				targetUrl,
				"GET",
				queryParams,
				headers,
				nil,
			),
		)
	}

	responses := fhc.MakeConcurrentRequests(requests, client)
	for _, concurrentResponse := range responses {
		if concurrentResponse.Error != nil {
			log.Fatalln(concurrentResponse.Error)
		}
		if matchesStatusCode(concurrentResponse.Response.StatusCode, matchCodes) {
			log.Println(concurrentResponse.Response.StatusCode)
			log.Println(concurrentResponse.Response.CurlCommand)
			log.Println("============")
		}

	}
}
