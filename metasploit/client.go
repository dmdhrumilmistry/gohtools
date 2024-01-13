package metasploit

import (
	"bytes"
	"fmt"
	"net/http"

	"log"

	"gopkg.in/vmihailenco/msgpack.v2"
)

type MsfClient struct {
	Host     string
	Username string
	Password string
	Token    string
}

func NewMsfClient(host, user, pass string) (*MsfClient, error) {
	msf := &MsfClient{Host: host, Username: user, Password: pass}

	if err := msf.Login(); err != nil {
		log.Fatalln("[!] Failed to log in!")
		return nil, err
	}

	return msf, nil
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

func (c *MsfClient) ListSessions() (map[uint32]sessionListRes, error) {
	var res map[uint32]sessionListRes
	ctx := &sessionListReq{
		Method: "session.list",
		Token:  c.Token,
	}

	err := c.Send(&ctx, &res)
	if err != nil {
		log.Printf("[MSF-CLIENT-ERROR] List sessions request failed due to error: %s\n", err)
		return nil, err
	}

	return res, nil
}
