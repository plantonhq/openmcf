package module

import (
	"regexp"
	"testing"
)

// The derivation contract: 10 digits, no leading zero, deterministic, and
// pinned to independently computed values so any change to the algorithm
// (which must stay in lockstep with the Terraform module's locals.tf)
// fails loudly here.
func TestDeriveEndpointName(t *testing.T) {
	numeric := regexp.MustCompile(`^[1-9][0-9]{9}$`)

	cases := []struct {
		org, env, name string
	}{
		{"", "", "my-endpoint"},
		{"acme", "prod", "recommendations"},
		{"planton-oss", "e2e", "planton-oss-e2e-gcpvep-minimal"},
	}
	for _, c := range cases {
		got := deriveEndpointName(c.org, c.env, c.name)
		if !numeric.MatchString(got) {
			t.Errorf("deriveEndpointName(%q,%q,%q) = %q; want 10 digits with no leading zero", c.org, c.env, c.name, got)
		}
		if again := deriveEndpointName(c.org, c.env, c.name); again != got {
			t.Errorf("deriveEndpointName not deterministic: %q then %q", got, again)
		}
	}

	// Pinned values computed independently of this implementation
	// (shell: sha256 of the identity string, first 12 hex chars, mapped
	// into [1000000000, 9999999999]). These are the exact IDs the
	// Terraform module's locals.tf derivation must also produce.
	pinned := []struct {
		org, env, name string
		want           string
	}{
		{"acme", "prod", "recommendations", "5835711114"},
		{"", "", "my-endpoint", "5035811386"},
		// The identity used in the first live cross-engine proof: Pulumi
		// created this exact ID and Terraform independently derived the
		// byte-identical value (colliding on GCP's deleted-ID reservation —
		// the collision itself was the determinism proof).
		{"planton-oss", "e2e", "planton-oss-e2e-gcpvep-minimal", "1853927074"},
	}
	for _, p := range pinned {
		if got := deriveEndpointName(p.org, p.env, p.name); got != p.want {
			t.Errorf("deriveEndpointName(%q,%q,%q) = %q; want pinned %q — the derivation drifted from the Terraform module's algorithm", p.org, p.env, p.name, got, p.want)
		}
	}
}
