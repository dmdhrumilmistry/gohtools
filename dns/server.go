package dns

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

type DnsServer struct {
	Port       *int
	Ip         *string
	listenAddr *string
}

func NewDnsServer(port int, ip string) *DnsServer {
	listenAddr := fmt.Sprintf("%s:%d", ip, port)

	return &DnsServer{
		Port:       &port,
		Ip:         &ip,
		listenAddr: &listenAddr,
	}
}

func (s *DnsServer) StartDnsServer(network string) {
	if network != "tcp" && network != "udp" {
		log.Printf("[!] Invalid Network option %s. Only tcp and udp are supported\n", network)
		return
	}

	dns.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
		log.Printf("[*] Incoming connection from: %s", w.RemoteAddr().String())

		var resp dns.Msg
		resp.SetReply(req)
		for _, q := range req.Question {
			a := dns.A{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				A: net.ParseIP("127.0.0.1").To4(),
			}

			resp.Answer = append(resp.Answer, &a)
		}

		w.WriteMsg(&resp)
	})

	log.Printf("[*] Starting DNS server on %s://%s\n", network, *s.listenAddr)
	log.Fatalln(dns.ListenAndServe(*s.listenAddr, network, nil))

	// dig @[SERVER_IP] [DOMAIN_NAME]
}
