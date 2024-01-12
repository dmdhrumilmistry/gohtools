package main

import (
	"fmt"
	"log"
	"net"
)

type connHandler func(conn net.Conn)

func StartListener(port int, handler connHandler) {
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[!] Cannot bind port due to error: %s\n", err)
	}
	log.Printf("[*] Listening on %s\n", addr)

	for {
		conn, err := listener.Accept()
		log.Printf("[*] Received connection from %s", conn.RemoteAddr().String())
		if err != nil {
			log.Printf("[!] Unable to accept connection due to error: %s", err)
		}

		go handler(conn)
	}
}
