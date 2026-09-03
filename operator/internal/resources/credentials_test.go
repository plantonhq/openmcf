package resources

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGeneratePassword_Length(t *testing.T) {
	pw, err := GeneratePassword()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != 32 {
		t.Errorf("expected password length 32, got %d", len(pw))
	}
}

func TestGeneratePassword_Unique(t *testing.T) {
	pw1, err := GeneratePassword()
	if err != nil {
		t.Fatalf("unexpected error generating password 1: %v", err)
	}
	pw2, err := GeneratePassword()
	if err != nil {
		t.Fatalf("unexpected error generating password 2: %v", err)
	}
	if pw1 == pw2 {
		t.Error("two consecutive passwords should not be identical")
	}
}

func TestGeneratePassword_URLSafe(t *testing.T) {
	pw, err := GeneratePassword()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, c := range pw {
		if !isURLSafeBase64Char(c) {
			t.Errorf("password contains non-URL-safe character: %c", c)
		}
	}
}

func isURLSafeBase64Char(c rune) bool {
	return (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') ||
		c == '-' || c == '_'
}

func TestNewCredentialSecret_Fields(t *testing.T) {
	ownerRef := &metav1.OwnerReference{
		APIVersion: "planton.ai/v1",
		Kind:       "PlantonPlatform",
		Name:       "test",
		UID:        "abc-123",
	}

	secret := NewCredentialSecret("test-db-credentials", "default", "root-password", "s3cret", ownerRef)

	if secret.Name != "test-db-credentials" {
		t.Errorf("expected name test-db-credentials, got %s", secret.Name)
	}
	if secret.Namespace != "default" {
		t.Errorf("expected namespace default, got %s", secret.Namespace)
	}
	if secret.Labels["app.kubernetes.io/managed-by"] != ManagedByLabel {
		t.Errorf("expected managed-by label %s, got %s",
			ManagedByLabel, secret.Labels["app.kubernetes.io/managed-by"])
	}
	if string(secret.Data["root-password"]) != "s3cret" {
		t.Errorf("expected password s3cret, got %s", string(secret.Data["root-password"]))
	}
	if len(secret.OwnerReferences) != 1 {
		t.Fatalf("expected 1 owner reference, got %d", len(secret.OwnerReferences))
	}
	if secret.OwnerReferences[0].Name != "test" {
		t.Errorf("expected owner ref name test, got %s", secret.OwnerReferences[0].Name)
	}
}

func TestNewCredentialSecret_NilOwnerRef(t *testing.T) {
	secret := NewCredentialSecret("test-redis-credentials", "ns", "password", "pw", nil)

	if len(secret.OwnerReferences) != 0 {
		t.Errorf("expected no owner references, got %d", len(secret.OwnerReferences))
	}
}

func TestNewCredentialSecret_TypeMeta(t *testing.T) {
	secret := NewCredentialSecret("test", "ns", "key", "val", nil)

	if secret.APIVersion != "v1" {
		t.Errorf("expected apiVersion v1, got %s", secret.APIVersion)
	}
	if secret.Kind != "Secret" {
		t.Errorf("expected kind Secret, got %s", secret.Kind)
	}
}

func TestNewCredentialSecret_OpaqueType(t *testing.T) {
	secret := NewCredentialSecret("test", "ns", "key", "val", nil)

	if secret.Type != "Opaque" {
		t.Errorf("expected Secret type Opaque, got %s", secret.Type)
	}
}
