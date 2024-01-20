package dns

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/miekg/dns"
	"gopkg.in/yaml.v3"
)

type ProxyConfig struct {
	Host   string `yaml:"host"`
	Domain string `yaml:"domain"`
}

type ProxiesConfig struct {
	Proxies []ProxyConfig `yaml:"proxies,omitempty"`
}

type DnsProxy struct {
	ProxyConfig          ProxiesConfig
	records              map[string]string
	ListenAddr           string
	proxyConfigFilePath  string
	DefaultDnsServerAddr string
}

func NewDnsProxy(filePath, listenAddr string, defaultDnsServerAddr string) (*DnsProxy, error) {
	if listenAddr == "" {
		listenAddr = "0.0.0.0:53"
	}

	if defaultDnsServerAddr == "" {
		defaultDnsServerAddr = "1.1.1.1:53"
	}

	dnsProxy := &DnsProxy{
		ListenAddr:           listenAddr,
		proxyConfigFilePath:  filePath,
		DefaultDnsServerAddr: defaultDnsServerAddr,
	}
	err := dnsProxy.ParseDnsProxiesFile(filePath)

	return dnsProxy, err
}

func (p *DnsProxy) ParseDnsProxiesFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("[!] Error while reading DNS proxy config: %s", err)
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&p.ProxyConfig)
	if err != nil {
		log.Printf("[!] Error while decoding DNS proxy config: %s", err)
		return err
	}

	// set records
	records := make(map[string]string)
	for _, proxyConfig := range p.ProxyConfig.Proxies {
		records[proxyConfig.Domain] = proxyConfig.Host
	}
	p.records = records

	return nil
}

func (p *DnsProxy) StartDnsProxy(network string) {
	if network != "tcp" && network != "udp" {
		log.Printf("[!] Invalid Network option %s. Only tcp and udp are supported\n", network)
		return
	}

	var recordLock sync.RWMutex

	dns.HandleFunc(".", func(w dns.ResponseWriter, m *dns.Msg) {
		if len(m.Question) == 0 {
			log.Printf("[!] 0 questions count found\n")
			dns.HandleFailed(w, m)
		}

		fqdn := m.Question[0].Name
		parts := strings.Split(fqdn, ".")
		if len(parts) >= 2 {
			fqdn = strings.Join(parts[:len(parts)-1], ".") // remove . from the end
		}
		log.Printf("[INCOMING] %s -> %s\n", w.RemoteAddr().String(), fqdn)

		// lock while reading to avoid race conditions
		recordLock.RLock()
		match := p.records[fqdn]
		recordLock.RUnlock()

		if match == "" {
			log.Printf("[!] Error occurred while getting match %s with status using default DNS server\n", match)
			match = p.DefaultDnsServerAddr
		}

		// proxy request to the target DNS server
		res, err := dns.Exchange(m, match)
		if err != nil {
			log.Printf("[!] Error occurred while exchanging message to server %s:%s\n", match, err)
			dns.HandleFailed(w, m)
			return
		}

		// send proxied response back to client
		if err := w.WriteMsg(res); err != nil {
			log.Printf("[!] Error occurred while writing message to client:%s\n", err)
			dns.HandleFailed(w, m)
			return
		}
	})

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGUSR1)

		for sig := range sigs {
			switch sig {
			case syscall.SIGUSR1:
				log.Println("[!] SIGUSR1: reloading records")
				recordLock.Lock()
				p.ParseDnsProxiesFile(p.proxyConfigFilePath)
				recordLock.Unlock()

			}
		}
	}()

	log.Printf("[*] Starting DNS server on %s://%s\n", network, p.ListenAddr)
	log.Fatalln(dns.ListenAndServe(p.ListenAddr, "udp", nil))
}
