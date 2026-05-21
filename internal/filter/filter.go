package filter

import (
	"fmt"
	"net"
)

// IP filters a list of resolved IPs and returns the first match
// for the requested address family (IPv4 or IPv6).
func IP(ips []net.IP, ipv6 bool) (string, error) {
	if ipv6 {
		for _, ip := range ips {
			if ip.To16() != nil && ip.To4() == nil {
				return ip.String(), nil
			}
		}
	} else {
		for _, ip := range ips {
			if ip.To4() != nil && ip.To16() != nil {
				return ip.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no suitable IP address found")
}
