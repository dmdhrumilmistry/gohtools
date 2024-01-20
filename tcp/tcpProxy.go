package tcp

import (
	"fmt"
	"io"
	"log"
	"net"
)

func HandleProxyRequest(src net.Conn, addr string) {
	dst, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("[!] Unable to connect to host: %s", addr)
		return
	}

	defer dst.Close()

	go func(dst net.Conn, src net.Conn, addr string) {
		if _, err := io.Copy(dst, src); err != nil {
			log.Printf("[!] Error Occur while proxying data: %s\n", addr)
		}
	}(dst, src, addr)

	if _, err := io.Copy(src, dst); err != nil {
		log.Printf("[!] Error Occur while reverting data: %s\n", addr)
	}
}

func StartTcpProxy(listenPort int, remoteAddr string) {
	listenAddr := fmt.Sprintf("0.0.0.0:%d", listenPort)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Printf("[!] Unable to start server due to error: %s\n", remoteAddr)
	}
	log.Printf("[*] Listening on %s", listenAddr)

	for {
		conn, err := listener.Accept()
		log.Printf("[*] Connection from %s", conn.RemoteAddr().String())
		if err != nil {
			log.Printf("[!] Unable to connect to remote server due to error: %s\n", remoteAddr)
		}

		go HandleProxyRequest(conn, remoteAddr)
	}
}
