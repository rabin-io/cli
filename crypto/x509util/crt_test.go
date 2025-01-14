package x509util

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"net"
	"net/url"
	"testing"

	"github.com/smallstep/assert"
)

func TestFingerprint(t *testing.T) {
	tests := []struct {
		name string
		fn   string
		want string
	}{
		{"ok", "test_files/ca.crt", "6908751f68290d4573ae0be39a98c8b9b7b7d4e8b2a6694b7509946626adfe98"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cert := mustParseCertificate(t, tt.fn)
			if got := Fingerprint(cert); got != tt.want {
				t.Errorf("Fingerprint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mustParseCertificate(t *testing.T, filename string) *x509.Certificate {
	pemData, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read %s: %v", filename, err)
	}
	block, rest := pem.Decode([]byte(pemData))
	if block == nil || len(rest) > 0 {
		t.Fatalf("failed to decode PEM in %s", filename)
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("failed to parse certificate in %s: %v", filename, err)
	}
	return cert
}

func TestSplitSANs(t *testing.T) {
	u1, err := url.Parse("https://ca.smallstep.com")
	assert.FatalError(t, err)
	u2, err := url.Parse("https://google.com/index.html")
	assert.FatalError(t, err)
	u3, err := url.Parse("urn:uuid:ddfe62ba-7e99-4bc1-83b3-8f57fe3e9959")
	assert.FatalError(t, err)
	tests := []struct {
		name              string
		sans, dns, emails []string
		ips               []net.IP
		uris              []*url.URL
	}{
		{name: "empty", sans: []string{}, dns: []string{}, ips: []net.IP{}, emails: []string{}, uris: []*url.URL{}},
		{
			name:   "all-dns",
			sans:   []string{"foo.internal", "bar.internal"},
			dns:    []string{"foo.internal", "bar.internal"},
			ips:    []net.IP{},
			emails: []string{},
			uris:   []*url.URL{},
		},
		{
			name:   "all-ip",
			sans:   []string{"0.0.0.0", "127.0.0.1"},
			dns:    []string{},
			ips:    []net.IP{net.ParseIP("0.0.0.0"), net.ParseIP("127.0.0.1")},
			emails: []string{},
			uris:   []*url.URL{},
		},
		{
			name:   "all-email",
			sans:   []string{"max@smallstep.com", "mariano@smallstep.com"},
			dns:    []string{},
			ips:    []net.IP{},
			emails: []string{"max@smallstep.com", "mariano@smallstep.com"},
			uris:   []*url.URL{},
		},
		{
			name:   "all-uri",
			sans:   []string{"https://ca.smallstep.com", "urn:uuid:ddfe62ba-7e99-4bc1-83b3-8f57fe3e9959", "https://google.com/index.html"},
			dns:    []string{},
			ips:    []net.IP{},
			emails: []string{},
			uris:   []*url.URL{u1, u3, u2},
		},
		{
			name:   "mix",
			sans:   []string{"foo.internal", "https://ca.smallstep.com", "max@smallstep.com", "urn:uuid:ddfe62ba-7e99-4bc1-83b3-8f57fe3e9959", "mariano@smallstep.com", "1.1.1.1", "bar.internal", "https://google.com/index.html"},
			dns:    []string{"foo.internal", "bar.internal"},
			ips:    []net.IP{net.ParseIP("1.1.1.1")},
			emails: []string{"max@smallstep.com", "mariano@smallstep.com"},
			uris:   []*url.URL{u1, u3, u2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dns, ips, emails, uris := SplitSANs(tt.sans)
			assert.Equals(t, dns, tt.dns)
			assert.Equals(t, ips, tt.ips)
			assert.Equals(t, emails, tt.emails)
			assert.Equals(t, uris, tt.uris)
		})
	}
}
