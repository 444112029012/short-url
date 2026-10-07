package domain

import (
	"errors"
	"regexp"
	"testing"
)

func TestGenerateFormat(t *testing.T) {
	if len(alphabet) != 62 {
		t.Fatalf("alphabet size %d", len(alphabet))
	}
	g := NewShortCodeGenerator()
	pattern := regexp.MustCompile(`^[A-Za-z0-9]{8}$`)
	seen := map[string]struct{}{}
	var previous string
	monotonic := true
	for i := 0; i < 100; i++ {
		code, err := g.Generate()
		if err != nil {
			t.Fatalf("UT-GEN-01: %v", err)
		}
		if !pattern.MatchString(code) {
			t.Fatalf("UT-GEN-01: %q", code)
		}
		seen[code] = struct{}{}
		if previous != "" && code <= previous {
			monotonic = false
		}
		previous = code
	}
	if len(seen) < 100 {
		t.Fatalf("UT-GEN-02: expected 100 distinct codes, got %d", len(seen))
	}
	if monotonic {
		t.Fatal("UT-GEN-02: 100 codes were strictly increasing")
	}
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("csprng unavailable at /tmp/secret")
}

func TestGenerateReaderFailure(t *testing.T) {
	g := NewShortCodeGeneratorWithReader(errReader{})
	code, err := g.Generate()
	if code != "" || err == nil || err.Type != TypeInternal {
		t.Fatalf("got %q %v", code, err)
	}
	if err.Error() != string(TypeInternal) {
		t.Fatalf("Error() leaked: %q", err.Error())
	}
}
