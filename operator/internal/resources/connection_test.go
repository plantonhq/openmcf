package resources

import (
	"testing"
)

// The contract every consumer's env is built from: the CloudNativePG
// cluster's read-write Service, the superuser, and the CloudNativePG-generated
// superuser Secret with basic-auth keys.
func TestPostgreSQLConnection(t *testing.T) {
	info := PostgreSQLConnection("my-planton", "planton-system")

	expectedHost := "my-planton-postgres-rw.planton-system.svc.cluster.local"
	if info.Host != expectedHost {
		t.Errorf("expected host %s, got %s", expectedHost, info.Host)
	}
	if info.Port != PostgreSQLPort {
		t.Errorf("expected port %d, got %d", PostgreSQLPort, info.Port)
	}
	if info.User != PostgreSQLSuperuser {
		t.Errorf("expected user %s, got %s", PostgreSQLSuperuser, info.User)
	}
	if info.SecretName != "my-planton-postgres-superuser" {
		t.Errorf("expected secret my-planton-postgres-superuser, got %s", info.SecretName)
	}
	if info.UserKey != "username" {
		t.Errorf("expected userKey 'username', got '%s'", info.UserKey)
	}
	if info.PassKey != "password" {
		t.Errorf("expected passKey 'password', got '%s'", info.PassKey)
	}
}

// ---------------------------------------------------------------------------
// Redis
// ---------------------------------------------------------------------------

func TestRedisConnection(t *testing.T) {
	info := RedisConnection("my-planton", "ns")
	if info.Host != RedisServiceHost("my-planton", "ns") {
		t.Errorf("host = %s, want %s", info.Host, RedisServiceHost("my-planton", "ns"))
	}
	if info.Port != int32(RedisPort) {
		t.Errorf("port = %d, want %d", info.Port, RedisPort)
	}
	if info.SecretName != RedisSecretName("my-planton") {
		t.Errorf("secretName = %s, want %s", info.SecretName, RedisSecretName("my-planton"))
	}
	if info.PassKey != RedisSecretKey {
		t.Errorf("passKey = %s, want %s", info.PassKey, RedisSecretKey)
	}
}

// ---------------------------------------------------------------------------
// OpenFGA
// ---------------------------------------------------------------------------

func TestOpenFGAConnection(t *testing.T) {
	info := OpenFGAConnection("my-planton", "ns")
	expectedURL := OpenFGAHTTPURL("my-planton", "ns")
	if info.HTTPURL != expectedURL {
		t.Errorf("httpURL = %s, want %s", info.HTTPURL, expectedURL)
	}
	if info.BootstrapConfigMapName != "my-planton-fga-bootstrap" {
		t.Errorf("bootstrapCM = %s, want my-planton-fga-bootstrap", info.BootstrapConfigMapName)
	}
}

// ---------------------------------------------------------------------------
// Temporal
// ---------------------------------------------------------------------------

func TestTemporalConnection(t *testing.T) {
	info := TemporalConnection("my-planton", "ns")
	expectedEndpoint := TemporalFrontendEndpoint("my-planton", "ns")
	if info.FrontendEndpoint != expectedEndpoint {
		t.Errorf("frontendEndpoint = %s, want %s", info.FrontendEndpoint, expectedEndpoint)
	}
	if info.Namespace != "default" {
		t.Errorf("namespace = %s, want default", info.Namespace)
	}
}
