package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type ApiClient struct {
}

func NewApiClient() *ApiClient {
	return &ApiClient{}
}

func (c *ApiClient) MethodHandler(method string, url string, form url.Values) []byte {
	var client http.Client
	method = strings.ToUpper(method)
	req, err := http.NewRequest(method, url, strings.NewReader(form.Encode()))
	if err != nil {
		log.Printf("[ERROR] Error while creating HTTP request: %s", err)
		return nil
	}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Error occurred sending %s request %s: %s\n", method, url, err)
		return nil
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[ERROR] Error occurred reading response data for %s: %s\n", url, err)
		return nil
	}

	return body
}

func (c *ApiClient) Get(url string) []byte {
	return c.MethodHandler("GET", url, nil)
}

func (c *ApiClient) Post(url string, form url.Values) []byte {
	return c.MethodHandler("POST", url, form)
}

func (c *ApiClient) Put(url string, form url.Values) []byte {
	return c.MethodHandler("PUT", url, form)
}

func (c *ApiClient) Patch(url string, form url.Values) []byte {
	return c.MethodHandler("PATCH", url, form)
}

func (c *ApiClient) Delete(url string) []byte {
	return c.MethodHandler("DELETE", url, nil)
}

func main() {
	var res []byte
	client := NewApiClient()
	URL := "https://example.com/robots.txt"

	form := url.Values{}
	form.Add("test", "true")

	res = client.Get(URL)
	// res = client.Put(URL, form)
	// res = client.Patch(URL, form)
	// res = client.Post(URL, "application/x-www-form-urlencoded", form)
	// res = client.Delete(URL)
	fmt.Println(string(res))
}
