package resources

import "fmt"

// The platform's cache-server role is named "redis" (the wire protocol every
// consumer speaks), but the engine that serves it is Valkey -- the
// BSD-3-Clause, Linux-Foundation fork that ships a redis-compatible server.
// Valkey is chosen over Redis deliberately: Redis 8+ is tri-licensed
// (RSALv2/SSPLv1/AGPLv3), a family self-hosted customers' legal teams often
// cannot accept, while Valkey stays permissive. This mirrors the desktop
// instance, which supervises a valkey-server behind REDIS_* configuration.
// The role keeps the "redis" name across the CRD, status, Secret, Service,
// and connection surfaces so the engine choice is invisible to consumers.
const (
	ValkeyHelmChartVersion = "3.0.31"
	RedisPort              = 6379

	// Valkey image coordinates. Chart v3.0.31 (appVersion 8.1.3) defaults to
	// bitnami/valkey, which was removed from Docker Hub when Bitnami
	// deprecated the free registry -- so the image is pinned to the frozen
	// bitnamilegacy mirror at the matching app version.
	valkeyImageRegistry   = "docker.io"
	valkeyImageRepository = "bitnamilegacy/valkey"
	valkeyImageTag        = "8.1.3-debian-12-r3"

	// RedisSecretKey is the data key inside the operator-generated credential
	// Secret that holds the cache-server password. The Valkey chart is
	// configured to read this key via auth.existingSecret +
	// auth.existingSecretPasswordKey.
	RedisSecretKey = "redis-password"
)

// ValkeyHelmValues builds the Helm values map for rendering the Bitnami Valkey
// chart. The chart is configured in standalone architecture with persistence
// enabled and password authentication via an operator-managed Secret:
//   - fullnameOverride keeps the role-named resources ("{crName}-redis-*")
//   - standalone architecture (no sentinel/replica overhead)
//   - Password from existingSecret (operator-generated, never in values)
//   - Persistence enabled with configurable size
//
// The Valkey chart names the data-serving workload "primary" (StatefulSet and
// Service "{fullname}-primary"), where the Redis chart said "master" -- the
// readiness check and RedisServiceHost follow that naming.
//
// storageClass pins the data volume's StorageClass; empty means the key is
// OMITTED so the cluster default provisions (an explicit "" would disable
// dynamic provisioning in the Bitnami convention).
func ValkeyHelmValues(crName, storageSize, storageClass string) map[string]any {
	persistence := map[string]any{
		"enabled": true,
		"size":    storageSize,
	}
	if storageClass != "" {
		persistence["storageClass"] = storageClass
	}
	return map[string]any{
		"fullnameOverride": redisReleaseName(crName),
		"architecture":     "standalone",
		"image": map[string]any{
			"registry":   valkeyImageRegistry,
			"repository": valkeyImageRepository,
			"tag":        valkeyImageTag,
		},
		"auth": map[string]any{
			"existingSecret":            RedisSecretName(crName),
			"existingSecretPasswordKey": RedisSecretKey,
		},
		"primary": map[string]any{
			"persistence": persistence,
		},
	}
}

// redisReleaseName returns the Helm release name: "{crName}-redis".
func redisReleaseName(crName string) string {
	return fmt.Sprintf("%s-redis", crName)
}

// RedisSecretName returns the credential Secret name: "{crName}-redis-credentials".
func RedisSecretName(crName string) string {
	return fmt.Sprintf("%s-redis-credentials", crName)
}

// RedisServiceHost returns the in-cluster DNS hostname for the cache server's
// primary Service created by the Valkey chart:
// "{crName}-redis-primary.{namespace}.svc.cluster.local".
func RedisServiceHost(crName, namespace string) string {
	return fmt.Sprintf("%s-redis-primary.%s.svc.cluster.local", crName, namespace)
}
