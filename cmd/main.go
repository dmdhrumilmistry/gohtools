package main

import (
	"log"
	"os"

	"github.com/dmdhrumilmistry/gohtools/metasploit"
)

func main() {

	// TCP Servers
	// StartEchoServer(4444, "ioCopyEcho")
	// StartTcpProxy(4444, "localhost:8000")
	// StartListener(4444, HandleNcConn) // connect with server using net.Dial or nc IP 4444

	// HTTP request
	// ip := http.GetHostMachinePublicIP() // package to import: "github.com/dmdhrumilmistry/gohtools/http"
	// fmt.Printf("[*] Host Machine Public Ip: %s\n", ip)

	// msf login
	msfHost := os.Getenv("MSF_HOST")
	msfUsername := os.Getenv("MSF_USERNAME")
	msfPassword := os.Getenv("MSF_PASSWORD")

	msf := metasploit.NewMsfClient(msfHost, msfUsername, msfPassword)
	if err := msf.Login(); err != nil {
		log.Fatalln("[!] Failed to log in!")
	}
	log.Printf("[INFO] MSF RPC TOKEN: %s\n", msf.Token)
}
