package ratelimit

import (
	"net"
	"net/http"
	"strings"
)

// ParseTrustedCIDRs parses a comma-separated CIDR list. An empty string
// yields no trusted proxies, which means forwarding headers are never honored.
func ParseTrustedCIDRs(raw string) ([]*net.IPNet, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	var nets []*net.IPNet
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		_, network, err := net.ParseCIDR(part)
		if err != nil {
			return nil, err
		}
		nets = append(nets, network)
	}
	return nets, nil
}

// ClientSource is the rate-limit identity for a request.
//
// The default is the host part of r.RemoteAddr. X-Forwarded-For and
// X-Real-IP are read only when that immediate peer is inside trusted.
// An untrusted peer cannot spoof its identity by setting those headers.
//
// When the peer is trusted and X-Forwarded-For contains a valid address,
// the chosen address is the rightmost entry that is not itself in trusted
// (the client the nearest trusted hop observed). Trusted-proxy entries are
// skipped from the right. If every entry is trusted, the rightmost valid
// address is used. If X-Forwarded-For has no valid address, X-Real-IP is
// used when it is a single valid IP. Otherwise the peer host is used.
//
// The returned value may be empty; Guard maps that onto one shared bucket.
// This function does not log the address (SEC-005).
func ClientSource(r *http.Request, trusted []*net.IPNet) string {
	if r == nil {
		return ""
	}
	peerHost := remoteHost(r.RemoteAddr)
	peerIP := parseIP(peerHost)
	if peerIP == nil || !ipInNets(peerIP, trusted) {
		return peerHost
	}
	if ip := clientFromForwarded(r.Header.Get("X-Forwarded-For"), trusted); ip != nil {
		return ip.String()
	}
	if ip := parseIP(r.Header.Get("X-Real-IP")); ip != nil {
		return ip.String()
	}
	return peerIP.String()
}

func remoteHost(remoteAddr string) string {
	remoteAddr = strings.TrimSpace(remoteAddr)
	if remoteAddr == "" {
		return ""
	}
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}
	return strings.Trim(host, "[]")
}

func clientFromForwarded(header string, trusted []*net.IPNet) net.IP {
	if strings.TrimSpace(header) == "" {
		return nil
	}
	parts := strings.Split(header, ",")
	var rightmost net.IP
	for i := len(parts) - 1; i >= 0; i-- {
		ip := parseIP(parts[i])
		if ip == nil {
			continue
		}
		if rightmost == nil {
			rightmost = ip
		}
		if !ipInNets(ip, trusted) {
			return ip
		}
	}
	return rightmost
}

func parseIP(raw string) net.IP {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	if host, _, err := net.SplitHostPort(raw); err == nil {
		raw = host
	}
	raw = strings.Trim(raw, "[]")
	return net.ParseIP(raw)
}

func ipInNets(ip net.IP, nets []*net.IPNet) bool {
	if ip == nil {
		return false
	}
	for _, n := range nets {
		if n != nil && n.Contains(ip) {
			return true
		}
	}
	return false
}
