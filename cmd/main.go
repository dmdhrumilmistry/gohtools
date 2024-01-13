package main

import (
	"fmt"

	"github.com/dmdhrumilmistry/gohtools/http"
)

func main() {
	// StartEchoServer(4444, "ioCopyEcho")
	// StartTcpProxy(4444, "localhost:8000")
	// StartListener(4444, HandleNcConn) // connect with server using net.Dial or nc IP 4444
	ip := http.GetHostMachinePublicIP()
	fmt.Printf("[*] Host Machine Public Ip: %s\n", ip)
}
