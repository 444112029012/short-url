package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewHTTPServerMaxHeaderBytes(t *testing.T) {
	srv := newHTTPServer(":0", http.NotFoundHandler())
	if srv.MaxHeaderBytes <= 0 || srv.MaxHeaderBytes > 16384 {
		t.Fatalf("MaxHeaderBytes %d, want 1..16384", srv.MaxHeaderBytes)
	}
	if srv.MaxHeaderBytes != 16384 {
		t.Fatalf("MaxHeaderBytes %d, want 16384", srv.MaxHeaderBytes)
	}
	// Existing timeouts must stay unchanged.
	if srv.ReadHeaderTimeout != 5*time.Second || srv.ReadTimeout != 15*time.Second ||
		srv.WriteTimeout != 15*time.Second || srv.IdleTimeout != 60*time.Second {
		t.Fatalf("timeouts changed: %+v", srv)
	}
}

func TestNewHTTPServerRejectsOversizedHeaders(t *testing.T) {
	cfg := newHTTPServer("", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	ts := httptest.NewUnstartedServer(cfg.Handler)
	ts.Config = cfg
	ts.Start()
	t.Cleanup(ts.Close)

	status := func(headerBytes int) (int, string) {
		t.Helper()
		conn, err := net.Dial("tcp", ts.Listener.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
		_, err = fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: x\r\nX-Pad: %s\r\nConnection: close\r\n\r\n", strings.Repeat("a", headerBytes))
		if err != nil {
			t.Fatal(err)
		}
		res, err := http.ReadResponse(bufio.NewReader(conn), nil)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		return res.StatusCode, res.Header.Get("Content-Type")
	}
	if got, _ := status(1024); got != http.StatusNoContent {
		t.Fatalf("small header status %d", got)
	}
	// A header just under 16 KiB is still served.
	if got, _ := status(16000); got != http.StatusNoContent {
		t.Fatalf("16000-byte header status %d", got)
	}
	// SCR-025: net/http adds 4096 bytes of slack above MaxHeaderBytes, so the
	// effective cap is about 20 KiB. A ~21 KiB header and a 32 KiB header get
	// 431, written by net/http itself (text/plain), not by the router.
	for _, size := range []int{21 * 1024, 32 * 1024} {
		got, ct := status(size)
		if got != http.StatusRequestHeaderFieldsTooLarge {
			t.Fatalf("%d-byte header status %d, want 431", size, got)
		}
		if !strings.HasPrefix(ct, "text/plain") {
			t.Fatalf("%d-byte header 431 Content-Type %q, want text/plain", size, ct)
		}
	}
}
