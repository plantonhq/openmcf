package component

import (
	"testing"

	v1 "github.com/plantonhq/planton/operator/api/v1"
)

// The bundled secrets manager is deployed by default: it backs the credential
// store, the envelope-encryption KEK, and the OIDC signing key. Opting out is
// the explicit act. The three predicates that answer "is the vault on?" must
// agree, so all three are pinned here against the same arms.
func TestOpenBAO_IsEnabledByDefault(t *testing.T) {
	o := &OpenBAO{}

	p := ingressPlatform(false)
	if !o.IsEnabled(p) || !isVaultEnabled(p) {
		t.Error("vault must be enabled with no spec.vault at all")
	}

	p.Spec.Vault = &v1.OpenBAOSpec{}
	if !o.IsEnabled(p) || !isVaultEnabled(p) {
		t.Error("vault must be enabled with an empty spec.vault")
	}

	off := false
	p.Spec.Vault = &v1.OpenBAOSpec{Enabled: &off}
	if o.IsEnabled(p) || isVaultEnabled(p) {
		t.Error("explicit enabled=false must disable the vault")
	}

	on := true
	p.Spec.Vault = &v1.OpenBAOSpec{Enabled: &on}
	if !o.IsEnabled(p) || !isVaultEnabled(p) {
		t.Error("explicit enabled=true must enable the vault")
	}
}

// With nothing declared, the default install seeds the platform secret
// backend (the vault runs by default and exists to be the secret store); an
// explicit vault opt-out with nothing declared seeds nothing (console
// funnels); a declared backend always wins.
func TestEffectiveSecretBackend_FollowsVaultDefault(t *testing.T) {
	p := ingressPlatform(false)
	binding := effectiveSecretBackend(p)
	if binding == nil || binding.Type != "platform" {
		t.Fatalf("default install must seed the platform backend, got %+v", binding)
	}

	off := false
	p.Spec.Vault = &v1.OpenBAOSpec{Enabled: &off}
	if binding := effectiveSecretBackend(p); binding != nil {
		t.Errorf("vault opt-out with nothing declared must seed no backend, got %+v", binding)
	}

	p.Spec.Bootstrap = &v1.BootstrapSpec{
		SecretBackend: &v1.BootstrapSecretBackendSpec{
			Type: "awsSecretsManager",
			AwsSecretsManager: &v1.BootstrapAwsSecretsManagerSpec{
				Region: "ap-south-1",
			},
		},
	}
	binding = effectiveSecretBackend(p)
	if binding == nil || binding.Type != "aws-secrets-manager" {
		t.Errorf("declared backend must win regardless of the vault toggle, got %+v", binding)
	}
}
