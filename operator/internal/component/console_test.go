package component

import (
	"testing"

	v1 "github.com/plantonhq/planton/operator/api/v1"
)

// Every URL-bearing surface derives from ONE resolved front-door URL. Without
// ingress the gateway's localhost URL is deterministic (nothing to wait for);
// with ingress, the only legitimate wait is the auto-hostname window before
// the controller publishes an address.
func TestFrontDoorURL(t *testing.T) {
	gatewayPort := int32(9000)
	tests := []struct {
		name         string
		ingress      *v1.IngressSpec
		gateway      *v1.GatewaySpec
		statusURL    string
		wantURL      string
		wantResolved bool
	}{
		{
			name:         "no ingress -> the gateway's deterministic localhost URL",
			wantURL:      "http://localhost:8080",
			wantResolved: true,
		},
		{
			name:         "no ingress, custom local port -> the port sign-in is pinned to",
			gateway:      &v1.GatewaySpec{LocalPort: &gatewayPort},
			wantURL:      "http://localhost:9000",
			wantResolved: true,
		},
		{
			name:         "explicit hostname, plain HTTP",
			ingress:      &v1.IngressSpec{Enabled: true, Hostname: "planton.example.com"},
			wantURL:      "http://planton.example.com",
			wantResolved: true,
		},
		{
			name: "explicit hostname with TLS",
			ingress: &v1.IngressSpec{
				Enabled: true, Hostname: "planton.example.com",
				TLS: &v1.IngressTLSSpec{SecretName: "corp-cert"},
			},
			wantURL:      "https://planton.example.com",
			wantResolved: true,
		},
		{
			name:         "auto hostname not yet resolved -> wait, do not deploy twice",
			ingress:      &v1.IngressSpec{Enabled: true},
			wantURL:      "",
			wantResolved: false,
		},
		{
			name:         "auto hostname resolved via status",
			ingress:      &v1.IngressSpec{Enabled: true},
			statusURL:    "http://planton.203-0-113-7.sslip.io",
			wantURL:      "http://planton.203-0-113-7.sslip.io",
			wantResolved: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &v1.PlantonPlatform{}
			p.Spec.Ingress = tt.ingress
			p.Spec.Gateway = tt.gateway
			p.Status.ConsoleURL = tt.statusURL

			url, resolved := frontDoorURL(p)
			if url != tt.wantURL || resolved != tt.wantResolved {
				t.Errorf("frontDoorURL() = (%q, %v), want (%q, %v)",
					url, resolved, tt.wantURL, tt.wantResolved)
			}
		})
	}
}
