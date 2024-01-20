package main

import (
	"flag"
	"log"
	"os"

	"github.com/dmdhrumilmistry/gohtools/dns"
)

func main() {
	/*
		proxies.yaml config file example:
		---------------------------------
		proxies:
			- domain: attacker1.com
			  host: 127.0.0.1:2020
			- domain: attacker2.com
			  host: 127.0.0.1:2021

		# start malicious dns servers
		$ dnspoof -ip 0.0.0.0 -network udp -port 2020
		$ dnspoof -ip 0.0.0.0 -network udp -port 2021

		# start DNS proxy server
		$ dnsproxy -config proxies.yaml

		# make calls to dns server
		$ dig @[ip] attacker1.com
		$ dig @[ip] attacker2.com
		$ dig @[ip] example.com
	*/

	filePath := flag.String("config", "", "proxies config file path in yaml format")
	listenAddr := flag.String("addr", "0.0.0.0:53", "DNS proxy server address")
	defaultDnsServerAddr := flag.String("dnsaddr", "1.1.1.1:53", "Default DNS server address to use when no match is found")
	network := flag.String("network", "udp", "network type: tcp or udp")
	flag.Parse()

	if *filePath == "" {
		log.Println("config file path is required")
		os.Exit(-1)
	}

	proxy, err := dns.NewDnsProxy(*filePath, *listenAddr, *defaultDnsServerAddr)
	if err != nil {
		log.Panic(err)
	}

	proxy.StartDnsProxy(*network)
}
