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

	status := func(headerBytes int) int {
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
		return res.StatusCode
	}
	if got := status(1024); got != http.StatusNoContent {
		t.Fatalf("small header status %d", got)
	}
	// net/http allows 4096 bytes of slack above MaxHeaderBytes; 32 KiB is well past it.
	if got := status(32 * 1024); got != http.StatusRequestHeaderFieldsTooLarge {
		t.Fatalf("oversized header status %d, want 431", got)
	}
}
