package main

import (
	"flag"
	"log"
	"os"

	"github.com/dmdhrumilmistry/gohtools/dns"
)

func main() {
	dnsServer := flag.String("server", "1.1.1.1:53", "DNS server to use")
	domain := flag.String("domain", "", "The domain to perform subdomain fuzzing against")
	wordlist := flag.String("wordlist", "", "path of wordlist containing subdomains used for fuzzing")
	workers := flag.Int("workers", 100, "Number of workers to use")
	flag.Parse()

	if *domain == "" || *wordlist == "" {
		log.Printf("wordlist and domain are required. Use -h for more info.")
		os.Exit(-1)
	}

	dnsFuzzer := dns.NewDnsFuzzer(*dnsServer, *workers, *wordlist)
	results, err := dnsFuzzer.StartSubDomainFuzzer(*domain)
	if err != nil {
		log.Printf("[!] Error occurred while fuzzing sub-domains: %s", err)
	}
	dnsFuzzer.PrintResults(results)
}
