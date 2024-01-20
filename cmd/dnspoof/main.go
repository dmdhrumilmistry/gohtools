package main

import (
	"flag"

	"github.com/dmdhrumilmistry/gohtools/dns"
)

func main() {
	/*
		kill any dns server such dnsmasq before running this tool:
		ps -ef | grep dnsmasq
		sudo kill <id>
	*/

	ip := flag.String("ip", "0.0.0.0", "ip address of DNS server to listen for incoming dns connections")
	port := flag.Int("port", 53, "DNS server port")
	network := flag.String("network", "udp", "network type: tcp or udp")
	flag.Parse()

	dnsServer := dns.NewDnsServer(*port, *ip)
	dnsServer.StartDnsServer(*network)

}
