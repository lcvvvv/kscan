package dns

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestDNS(t *testing.T) {
	resolver := &net.Resolver{
		PreferGo: false,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 3 * time.Second,
			}
			return d.DialContext(ctx, "udp", "114.114.114.114:53")
		},
	}
	resolvers = append(resolvers, resolver)
	fmt.Println(resolver.LookupCNAME(context.Background(), "www.t00ls.cc"))
	fmt.Println(resolver.LookupHost(context.Background(), "www.t00ls.cc"))
	fmt.Println(resolver.LookupAddr(context.Background(), "www.t00ls.cc"))
	fmt.Println(resolver.LookupIPAddr(context.Background(), "www.t00ls.cc"))
	fmt.Println(resolver.LookupTXT(context.Background(), "www.t00ls.cc"))
}
