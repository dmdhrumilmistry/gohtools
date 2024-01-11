package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Scanner struct {
	startPort int
	endPort   int
	host      string
	timeout   uint
	wg        sync.WaitGroup
}

func (s *Scanner) TestConnection(network string, port int) {
	defer s.wg.Done()
	if !(network == "tcp" || network == "udp") {
		log.Printf("Invalid Network Type: %s", network)
		return
	}

	address := fmt.Sprintf("%s:%d", s.host, port)
	dialer := net.Dialer{Timeout: time.Duration(s.timeout) * time.Second}
	conn, err := dialer.Dial(network, address)

	if err != nil {
		// fmt.Printf("[ERROR] %v\n", err)
		return
	}
	conn.Close()
	fmt.Printf("[%s OPEN] %d\n", strings.ToUpper(network), port)
}

func (s *Scanner) TestAllConnections(network string) {
	for port := s.startPort; port <= s.endPort; port++ {
		s.wg.Add(1)
		go s.TestConnection(network, port)
	}
}

func StartScanner(startPort int, endPort int, host string, timeout int) {
	scanner := Scanner{
		startPort: startPort,
		endPort:   endPort,
		host:      host,
		timeout:   uint(timeout),
	}
	scanner.TestAllConnections("tcp")
	// scanner.TestAllConnections("udp")

	scanner.wg.Wait()
}
