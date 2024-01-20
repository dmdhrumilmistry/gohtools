package dns

import (
	"errors"
	"fmt"

	"github.com/miekg/dns"
)

type Result struct {
	IpAddress string
	Hostname  string
}

func (r *Result) Text() string {
	return fmt.Sprintf("%-50s\t%s", r.Hostname, r.IpAddress)
}

type DnsClient struct {
	Server *string
}

func NewDnsClient(server *string) *DnsClient {
	return &DnsClient{
		Server: server,
	}
}

func (c *DnsClient) GetARecord(domain string) ([]string, error) {
	var msg dns.Msg
	var ips []string

	fqdn := dns.Fqdn(domain)
	msg.SetQuestion(fqdn, dns.TypeA)
	in, err := dns.Exchange(&msg, *c.Server)
	if err != nil {
		return ips, err
	}

	if len(in.Answer) < 1 {
		return ips, errors.New("no A records found")
	}

	for _, answer := range in.Answer {
		if a, ok := answer.(*dns.A); ok {
			ips = append(ips, a.A.String())
		}
	}

	return ips, nil
}

func (c *DnsClient) GetCNameRecord(domain string) ([]string, error) {
	var msg dns.Msg
	var cnames []string

	fqdn := dns.Fqdn(domain)
	msg.SetQuestion(fqdn, dns.TypeCNAME)
	in, err := dns.Exchange(&msg, *c.Server)
	if err != nil {
		return cnames, err
	}

	if len(in.Answer) < 1 {
		return cnames, errors.New("no CName records found")
	}

	for _, answer := range in.Answer {
		if c, ok := answer.(*dns.CNAME); ok {
			cnames = append(cnames, c.Target)
		}
	}

	return cnames, nil
}

// Follows CName trail till ips are resolved
func (c *DnsClient) LookupCName(fqdn string) []Result {
	var results []Result
	// var remainingCNames []string
	cfqdn := fqdn

	for {
		cnames, err := c.GetCNameRecord(cfqdn)
		if err != nil && len(cnames) > 0 {
			cfqdn = cnames[0]
			continue
		}

		ips, err := c.GetARecord(cfqdn)
		if err != nil {
			break
		}

		for _, ip := range ips {
			results = append(results, Result{
				Hostname:  cfqdn,
				IpAddress: ip,
			})
		}

		break
	}

	return results
}
