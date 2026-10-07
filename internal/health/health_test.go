package health

import "testing"

func TestStatus(t *testing.T) {
	if Status() != "ok" {
		t.Fatalf("want ok")
	}
}
