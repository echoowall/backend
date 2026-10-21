package epay

import (
	"net/url"
	"strings"
	"testing"
)

func TestResolveEndpoint(t *testing.T) {
	for _, tc := range []struct {
		raw, domain, want string
		relative, fail    bool
	}{
		{raw: "https://pay.example", want: "https://pay.example"},
		{raw: "//pay.example", want: "//pay.example"},
		{raw: "/pay", domain: "https://api.example/", want: "https://api.example/pay", relative: true},
		{raw: " /pay ", domain: "https://api.example", want: "https://api.example/pay", relative: true},
		{raw: "/pay", fail: true},
	} {
		got, relative, err := ResolveEndpoint(tc.raw, tc.domain)
		if (err != nil) != tc.fail || got != tc.want || relative != tc.relative {
			t.Errorf("ResolveEndpoint(%q, %q) = %q, %v, %v", tc.raw, tc.domain, got, relative, err)
		}
	}
	endpoint, _, _ := ResolveEndpoint("/pay", "https://api.example")
	payURL, err := NewClient("1", endpoint, "k", "alipay").CreatePayUrl(Order{Name: "n", OrderNo: "O1", Amount: 1, NotifyUrl: "https://api.example/n"})
	if err != nil {
		t.Fatal(err)
	}
	parsed, _ := url.Parse(payURL)
	if rel := parsed.RequestURI(); !strings.HasPrefix(rel, "/pay/submit.php?") {
		t.Fatalf("relative pay URL %q", rel)
	}
}
