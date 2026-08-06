package module

var vars = struct {
	// KarapaceImage is the module's pinned upstream image, used when
	// spec.image is empty. The tag pins the aiven-open/karapace 6.2.1
	// release (upstream publishes ghcr.io/aiven-open/karapace); bump it
	// deliberately, in lockstep with the Terraform module's
	// default_image local.
	KarapaceImage string

	// RegistryCommand and RestCommand are the container entrypoints from
	// upstream's own deployment reference (container/compose.yml): the
	// production image declares NO ENTRYPOINT/CMD of its own, and compose
	// starts the schema-registry role with `python3 -m karapace` and the
	// REST-proxy role with `python3 -m karapace.kafka_rest_apis`. The
	// KARAPACE_KARAPACE_REGISTRY / KARAPACE_KARAPACE_REST flags select
	// which API surface the started process serves.
	RegistryCommand []string
	RestCommand     []string

	// Mount points for Secret-sourced files. Configuration reaches the
	// engine as file PATHS (ssl_cafile, server_tls_certfile,
	// registry_authfile, ...), so every Secret mounts at a fixed,
	// role-independent directory and the env vars point inside it. Must
	// stay byte-identical with the Terraform module's locals.
	KafkaCaMountPath         string
	KafkaClientCertMountPath string
	ServerTlsMountPath       string
	AuthfileMountPath        string

	// HealthCheckPath is the engine's unauthenticated health endpoint —
	// config.py ships "/_health" in sasl_oauthbearer_skip_auth_paths and
	// the upstream image's own HEALTHCHECK curls it, so probes keep
	// working with HTTP authentication enabled.
	HealthCheckPath string
}{
	KarapaceImage:            "ghcr.io/aiven-open/karapace:6.2.1",
	RegistryCommand:          []string{"python3", "-m", "karapace"},
	RestCommand:              []string{"python3", "-m", "karapace.kafka_rest_apis"},
	KafkaCaMountPath:         "/etc/karapace/kafka-ca",
	KafkaClientCertMountPath: "/etc/karapace/kafka-cert",
	ServerTlsMountPath:       "/etc/karapace/server-tls",
	AuthfileMountPath:        "/etc/karapace/auth",
	HealthCheckPath:          "/_health",
}
