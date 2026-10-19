package ratelimit

import (
	"net/http"
	"testing"
)

func TestClientSourceIgnoresSpoofedHeadersWhenUntrusted(t *testing.T) {
	req := httptestRequest("203.0.113.9:443")
	req.Header.Set("X-Forwarded-For", "198.51.100.8, 198.51.100.9")
	req.Header.Set("X-Real-IP", "198.51.100.7")
	if got := ClientSource(req, nil); got != "203.0.113.9" {
		t.Fatalf("untrusted peer source %q", got)
	}
	trusted, err := ParseTrustedCIDRs("10.0.0.0/8")
	if err != nil {
		t.Fatal(err)
	}
	if got := ClientSource(req, trusted); got != "203.0.113.9" {
		t.Fatalf("peer outside trusted list source %q", got)
	}
}

func TestClientSourceTrustedProxyUsesRightmostUntrustedForwarded(t *testing.T) {
	trusted, err := ParseTrustedCIDRs("10.0.0.0/8, 192.0.2.0/24")
	if err != nil {
		t.Fatal(err)
	}
	req := httptestRequest("10.0.0.5:5000")
	req.Header.Set("X-Forwarded-For", "198.51.100.8, 10.0.0.9")
	req.Header.Set("X-Real-IP", "203.0.113.1")
	if got := ClientSource(req, trusted); got != "198.51.100.8" {
		t.Fatalf("forwarded source %q", got)
	}

	req.Header.Set("X-Forwarded-For", "198.51.100.8, 198.51.100.50")
	if got := ClientSource(req, trusted); got != "198.51.100.50" {
		t.Fatalf("rightmost source %q", got)
	}
}

func TestClientSourceTrustedProxyFallsBackToXRealIP(t *testing.T) {
	trusted, err := ParseTrustedCIDRs("127.0.0.1/32")
	if err != nil {
		t.Fatal(err)
	}
	req := httptestRequest("127.0.0.1:1234")
	req.Header.Set("X-Real-IP", "198.51.100.40")
	if got := ClientSource(req, trusted); got != "198.51.100.40" {
		t.Fatalf("x-real-ip source %q", got)
	}

	req.Header.Set("X-Forwarded-For", "not-an-ip")
	if got := ClientSource(req, trusted); got != "198.51.100.40" {
		t.Fatalf("invalid xff should fall back, got %q", got)
	}
}

func TestClientSourceIPv6AndEmpty(t *testing.T) {
	req := httptestRequest("[2001:db8::1]:443")
	if got := ClientSource(req, nil); got != "2001:db8::/64" {
		t.Fatalf("ipv6 %q", got)
	}
	if got := ClientSource(nil, nil); got != "" {
		t.Fatalf("nil request %q", got)
	}
	blank := httptestRequest("")
	if got := ClientSource(blank, nil); got != "" {
		t.Fatalf("blank remote %q", got)
	}
}

func TestIPv6PrefixSharesQuota(t *testing.T) {
	sameA := ClientSource(httptestRequest("[2001:db8:1:2::1]:1000"), nil)
	sameB := ClientSource(httptestRequest("[2001:db8:1:2:abcd::9]:2000"), nil)
	other := ClientSource(httptestRequest("[2001:db8:1:3::1]:1000"), nil)
	if sameA != "2001:db8:1:2::/64" || sameA != sameB {
		t.Fatalf("same prefix keys %q %q", sameA, sameB)
	}
	if other != "2001:db8:1:3::/64" || other == sameA {
		t.Fatalf("other prefix %q", other)
	}
	g := NewGuard(1, 1)
	if err := g.CheckCreate(sameA); err != nil {
		t.Fatal(err)
	}
	assertRateLimited(t, g.CheckCreate(sameB))
	if err := g.CheckCreate(other); err != nil {
		t.Fatal("different /64 was limited")
	}
	if err := g.CheckRedirect(sameA); err != nil {
		t.Fatal(err)
	}
	assertRateLimited(t, g.CheckRedirect(sameB))
	if err := g.CheckRedirect(other); err != nil {
		t.Fatal("different /64 redirect was limited")
	}
}

func TestIPv4MappedEqualsPlainIPv4(t *testing.T) {
	plain := ClientSource(httptestRequest("203.0.113.5:9"), nil)
	mapped := ClientSource(httptestRequest("[::ffff:203.0.113.5]:9"), nil)
	if plain != "203.0.113.5" || mapped != plain {
		t.Fatalf("plain %q mapped %q", plain, mapped)
	}
	g := NewGuard(1, 120)
	if err := g.CheckCreate(plain); err != nil {
		t.Fatal(err)
	}
	assertRateLimited(t, g.CheckCreate(mapped))
}

func TestTrustedProxyIPv6IsAggregated(t *testing.T) {
	trusted, err := ParseTrustedCIDRs("10.0.0.0/8")
	if err != nil {
		t.Fatal(err)
	}
	first := httptestRequest("10.1.0.1:4000")
	first.Header.Set("X-Forwarded-For", "2001:db8:9:8::1")
	second := httptestRequest("10.1.0.2:4001")
	second.Header.Set("X-Forwarded-For", "2001:db8:9:8:ffff::2")
	third := httptestRequest("10.1.0.3:4002")
	third.Header.Set("X-Forwarded-For", "2001:db8:9:9::1")
	a := ClientSource(first, trusted)
	b := ClientSource(second, trusted)
	c := ClientSource(third, trusted)
	if a != "2001:db8:9:8::/64" || a != b || c != "2001:db8:9:9::/64" {
		t.Fatalf("trusted ipv6 keys %q %q %q", a, b, c)
	}
	g := NewGuard(1, 1)
	if err := g.CheckRedirect(a); err != nil {
		t.Fatal(err)
	}
	assertRateLimited(t, g.CheckRedirect(b))
	if err := g.CheckRedirect(c); err != nil {
		t.Fatal(err)
	}
}

func TestForwardedForSeparateHeaderLineDoesNotLetClientWin(t *testing.T) {
	trusted, err := ParseTrustedCIDRs("10.0.0.0/8")
	if err != nil {
		t.Fatal(err)
	}
	req := httptestRequest("10.0.0.8:5000")
	// Client line first; the proxy appends the observed client as a second line.
	req.Header["X-Forwarded-For"] = []string{"198.51.100.8", "203.0.113.50"}
	req.Header.Set("X-Real-IP", "198.51.100.8")
	if got := ClientSource(req, trusted); got != "203.0.113.50" {
		t.Fatalf("appended XFF line lost: %q", got)
	}

	req.Header["X-Forwarded-For"] = []string{"198.51.100.8", "10.0.0.9"}
	if got := ClientSource(req, trusted); got != "198.51.100.8" {
		t.Fatalf("trusted appended hop should be skipped, got %q", got)
	}
}

func TestParseTrustedCIDRsRejectsBareIP(t *testing.T) {
	if _, err := ParseTrustedCIDRs("10.0.0.1"); err == nil {
		t.Fatal("bare IP accepted")
	}
	nets, err := ParseTrustedCIDRs(" , ")
	if err != nil || nets != nil {
		t.Fatalf("blank list: %v %v", nets, err)
	}
}

func httptestRequest(remote string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "http://example.test/", nil)
	if err != nil {
		panic(err)
	}
	req.RemoteAddr = remote
	return req
}
