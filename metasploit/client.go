package metasploit

import (
	"bytes"
	"fmt"
	"net/http"

	"log"

	"gopkg.in/vmihailenco/msgpack.v2"
)

type loginReq struct {
	_msgpack struct{} `msgpack:",asArray"`
	Method   string
	Username string
	Password string
}

type loginRes struct {
	Result       string `msgpack:"result"`
	Token        string `msgpack:"token"`
	Error        bool   `msgpack:"error"`
	ErrorClass   string `msgpack:"error_class"`
	ErrorMessage string `msgpack:"error_message"`
}

type logoutReq struct {
	_msgpack    struct{} `msgpack:",asArray"`
	Method      string
	Token       string
	LogoutToken string
}

type logoutRes struct {
	Result string `msgpack:"result"`
}

type MsfClient struct {
	Host     string
	Username string
	Password string
	Token    string
}

func NewMsfClient(host, user, pass string) *MsfClient {
	return &MsfClient{Host: host, Username: user, Password: pass}
}

func (c *MsfClient) Send(msgReq interface{}, msgRes interface{}) error {
	var buff bytes.Buffer
	msgpack.NewEncoder(&buff).Encode(msgReq)
	url := fmt.Sprintf("http://%s/api", c.Host)
	res, err := http.Post(url, "binary/message-pack", &buff)
	if err != nil {
		log.Printf("[MSF-CLIENT-ERROR] API Request Failed due to error: %s\n", err)
		return err
	}

	defer res.Body.Close()

	if err := msgpack.NewDecoder(res.Body).Decode(&msgRes); err != nil {
		log.Printf("[MSF-CLIENT-ERROR] Unable to decode response body due to error: %s\n", err)
		return err
	}

	return nil
}

func (c *MsfClient) Login() error {
	ctx := &loginReq{
		Method:   "auth.login",
		Username: c.Username,
		Password: c.Password,
	}

	var res loginRes
	err := c.Send(&ctx, &res)
	if err != nil {
		log.Printf("[MSF-CLIENT-ERROR] Login request failed due to error: %s\n", err)
		return err
	}

	c.Token = res.Token
	return nil
}

func (c *MsfClient) Logout() error {
	ctx := &logoutReq{
		Method:      "auth.logout",
		Token:       c.Token,
		LogoutToken: c.Token,
	}

	var res logoutRes
	err := c.Send(&ctx, &res)
	if err != nil {
		log.Printf("[MSF-CLIENT-ERROR] Logout request failed due to error: %s\n", err)
		return err
	}

	c.Token = ""
	return nil
}
