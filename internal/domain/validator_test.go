package domain

import (
	"strings"
	"testing"
)

func TestValidateLongURL(t *testing.T) {
	v := NewURLValidator()
	valid2048 := "https://example.com/" + strings.Repeat("a", 2048-len("https://example.com/"))
	if len(valid2048) != 2048 {
		t.Fatalf("fixture length %d", len(valid2048))
	}
	tooLong := valid2048 + "a"

	tests := []struct {
		name    string
		id      string
		raw     string
		wantErr bool
	}{
		{name: "https", id: "UT-VAL-01", raw: "https://example.com/path?q=1"},
		{name: "http", id: "UT-VAL-01", raw: "http://example.com"},
		{name: "empty", id: "UT-VAL-02", raw: "", wantErr: true},
		{name: "not url", id: "UT-VAL-02", raw: "not a url", wantErr: true},
		{name: "ftp", id: "UT-VAL-03", raw: "ftp://example.com/file", wantErr: true},
		{name: "javascript", id: "UT-VAL-03", raw: "javascript:alert(1)", wantErr: true},
		{name: "length 2048", id: "UT-VAL-04", raw: valid2048},
		{name: "length 2049", id: "UT-VAL-05", raw: tooLong, wantErr: true},
		{name: "missing host", id: "UT-VAL-02", raw: "http://", wantErr: true},
		{name: "control char", id: "UT-VAL-03", raw: "https://example.com/\n", wantErr: true},
		{name: "space", id: "UT-VAL-02", raw: "https://example.com/a b", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.id+"/"+tt.name, func(t *testing.T) {
			got, err := v.ValidateLongURL(tt.raw)
			if tt.wantErr {
				if err == nil || err.Type != TypeInvalidURL {
					t.Fatalf("got (%v, %v), want invalid_url", got, err)
				}
				if err.Error() != string(TypeInvalidURL) {
					t.Fatalf("Error() leaked: %q", err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.URL != tt.raw {
				t.Fatalf("stored %q want original %q", got.URL, tt.raw)
			}
		})
	}
}

func TestValidateShortCode(t *testing.T) {
	v := NewURLValidator()
	ok, err := v.ValidateShortCode("Ab12Cd39")
	if err != nil || ok != "Ab12Cd39" {
		t.Fatalf("UT-VAL-06: got %q %v", ok, err)
	}
	for _, code := range []string{"", "short", "Ab12Cd390", "Ab12Cd3!", "Ab12Cd3_", "ab12cd3 "} {
		_, err := v.ValidateShortCode(code)
		if err == nil || err.Type != TypeInvalidURL || err.Message != MsgInvalidShortCode {
			t.Fatalf("UT-VAL-07: code %q got %v", code, err)
		}
	}
}
