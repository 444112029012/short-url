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
// When the peer is trusted, every X-Forwarded-For header line is joined in
// order and parsed from the right. The chosen address is the rightmost entry
// that is not itself in trusted (the client the nearest trusted hop observed).
// Trusted-proxy entries are skipped from the right. If every entry is trusted,
// the rightmost valid address is used. A client-supplied line therefore cannot
// win over a later line the proxy appended.
//
// X-Real-IP is consulted only when the peer is trusted and X-Forwarded-For
// has no valid address. Operators must configure that proxy to overwrite
// X-Real-IP; a value the proxy forwards unchanged would be treated as the client.
//
// IP bucket keys are aggregated (SCR-007): IPv4 and IPv4-mapped IPv6 become
// the IPv4 address; IPv6 becomes its /64 prefix. The same normalization
// applies to trusted-proxy-derived addresses.
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
		if peerIP == nil {
			return peerHost
		}
		return bucketKey(peerIP)
	}
	if ip := clientFromForwarded(joinedForwardedFor(r), trusted); ip != nil {
		return bucketKey(ip)
	}
	if ip := parseIP(r.Header.Get("X-Real-IP")); ip != nil {
		return bucketKey(ip)
	}
	return bucketKey(peerIP)
}

// joinedForwardedFor keeps every X-Forwarded-For line. Header.Get would
// return only the first, which lets a client-supplied line hide the hop a
// proxy appended as a separate header.
func joinedForwardedFor(r *http.Request) string {
	if r == nil {
		return ""
	}
	return strings.Join(r.Header.Values("X-Forwarded-For"), ",")
}

// bucketKey collapses addresses that should share one quota.
// IPv4-mapped IPv6 is reduced with To4 first. IPv6 is masked to /64.
func bucketKey(ip net.IP) string {
	if ip == nil {
		return ""
	}
	if v4 := ip.To4(); v4 != nil {
		return v4.String()
	}
	v6 := ip.To16()
	if v6 == nil {
		return ""
	}
	prefix := v6.Mask(net.CIDRMask(64, 128))
	return prefix.String() + "/64"
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
