package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

// handler function to echo received data
func EchoHandler(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 512)
	for {
		size, err := conn.Read(buffer)
		if err == io.EOF {
			log.Println("[*] Connection closed by client")
			break
		}
		if err != nil {
			log.Printf("[!] Error occurred while reading data: %s\n", err)
			break
		}
		// log.Printf("[*] Received %d bytes of data from client\n", size)

		// log.Println("[*] Writing Data")
		if _, err := conn.Write(buffer[0:size]); err != nil {
			log.Fatalf("[!] Error occurred while writing data: %s\n", err)
		}
	}
}

func BuffioEchoHandler(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	data, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("[!] Unable to read data due to error: %s\n", err)
	}

	writer := bufio.NewWriter(conn)
	writer.WriteString(data)
	if err != nil {
		log.Fatalf("[!] Unable to write data due to error: %s\n", err)
	}

	writer.Flush()
}

func IoCopyHandler(conn net.Conn) {
	defer conn.Close()

	if _, err := io.Copy(conn, conn); err != nil {
		log.Fatalf("[!] Unable to read/write data due to error: %s\n", err)
	}
}

func StartEchoServer(port int, handler string) {
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

		switch handler {
		case "buffioEcho":
			go BuffioEchoHandler(conn)
		case "ioCopyEcho":
			go IoCopyHandler(conn)
		default:
			go EchoHandler(conn)
		}
	}
}
