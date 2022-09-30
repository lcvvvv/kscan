package dns

import (
	"context"
	"github.com/miekg/dns"
	"kscan/lib/misc"
	"net"
	"time"
)

var domainServers = []string{
	"114.114.114.114:53",
	//"8.8.8.8:53",
	"223.6.6.6:53",
}

var resolvers = generateResolver()

func LookupCNAME(domain string) ([]string, error) {
	var lastErr error
	for _, domainServer := range domainServers {
		CNAMES, err := LookupCNAMEWithServer(domain, domainServer)
		if err != nil {
			lastErr = err
		}
		return CNAMES, nil
	}
	return nil, lastErr
}

func LookupCNAMEWithServer(domain, domainServer string) ([]string, error) {
	c := dns.Client{
		Timeout: 5 * time.Second,
	}
	var CNAMES []string
	m := dns.Msg{}
	// 最终都会指向一个ip 也就是typeA, 这样就可以返回所有层的cname.
	m.SetQuestion(domain+".", dns.TypeA)
	r, _, err := c.Exchange(&m, domainServer)
	if err != nil {
		return nil, err
	}
	for _, ans := range r.Answer {
		record, isType := ans.(*dns.CNAME)
		if isType {
			CNAMES = append(CNAMES, record.Target)
		}
	}
	return CNAMES, nil
}

func LookupIP(domain string) ([]string, error) {
	var IPs []string
	var lastErr error
	for _, resolver := range resolvers {
		ips, err := resolver.LookupIPAddr(context.Background(), domain)
		if err != nil {
			lastErr = err
		}
		for _, v := range ips {
			IPs = append(IPs, v.IP.String())
		}
	}
	IPs = misc.RemoveDuplicateElement(IPs)
	return IPs, lastErr
}

func generateResolver() []*net.Resolver {
	var resolvers []*net.Resolver
	for _, server := range domainServers {
		resolver := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: 3 * time.Second,
				}
				return d.DialContext(ctx, "udp", server)
			},
		}
		resolvers = append(resolvers, resolver)
	}
	return resolvers
}
