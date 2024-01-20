package dns

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
)

type Empty struct{}

type DnsFuzzer struct {
	Server       *string
	Workers      *int
	WordlistPath *string
	dnsClient    *DnsClient
}

func NewDnsFuzzer(dnsServer string, workers int, wordlistPath string) *DnsFuzzer {
	if !strings.Contains(dnsServer, ":53") {
		dnsServer = dnsServer + ":53"
	}

	if workers <= 0 {
		log.Printf("[!] Invalid worker %d count. Setting workers count to 100.\n", workers)
		workers = 100
	}

	return &DnsFuzzer{
		Server:       &dnsServer,
		Workers:      &workers,
		WordlistPath: &wordlistPath,
		dnsClient:    NewDnsClient(&dnsServer),
	}
}

func (c *DnsFuzzer) Worker(tracker chan Empty, fqdns chan string, gather chan []Result) {
	for fqdn := range fqdns {
		results := c.dnsClient.LookupCName(fqdn)
		if len(results) > 0 {
			for _, result := range results {
				log.Println(result.Text())
			}
			gather <- results
		}
	}

	var e Empty
	tracker <- e // to avoid race conditions
}

func (c *DnsFuzzer) StartSubDomainFuzzer(domain string) ([]Result, error) {
	var results []Result
	fqdns := make(chan string, *c.Workers)
	gather := make(chan []Result)
	tracker := make(chan Empty)

	// read wordlist file
	fh, err := os.Open(*c.WordlistPath)
	if err != nil {
		log.Printf("[!] Error occurred while opening file: %s\n", *c.WordlistPath)
		return results, err
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)

	// start fuzzing dns entries
	for i := 0; i < *c.Workers; i++ {
		go c.Worker(tracker, fqdns, gather)
	}

	for scanner.Scan() {
		fqdns <- fmt.Sprintf("%s.%s", scanner.Text(), domain)
	}

	err = scanner.Err()
	if err != nil {
		log.Printf("[!] Error Occurred while reading file: %s", err)
	}

	// gather results
	go func(tracker *chan Empty, gather *chan []Result) {
		for r := range *gather {
			results = append(results, r...)
		}

		var e Empty
		*tracker <- e
	}(&tracker, &gather)

	// close channels
	close(fqdns)
	for i := 0; i < *c.Workers; i++ {
		<-tracker
	}

	close(gather)
	<-tracker
	return results, nil
}

func (c *DnsFuzzer) PrintResults(results []Result) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 4, ' ', 0)

	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\n", r.Hostname, r.IpAddress)
	}

	w.Flush()
}
