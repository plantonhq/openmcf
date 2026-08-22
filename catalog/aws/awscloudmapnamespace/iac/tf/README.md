# AwsCloudMapNamespace — Terraform/OpenTofu module

Manages one Cloud Map namespace (`aws_service_discovery_http_namespace` XOR `aws_service_discovery_private_dns_namespace` XOR `aws_service_discovery_public_dns_namespace`, selected by the spec's type) with its services (`aws_service_discovery_service`, keyed by name) and statically registered instances (`aws_service_discovery_instance`, keyed `service//instance_id`).

Module facts worth knowing before editing:

- **The three types are three provider resources**; exactly one exists, all exposing the same downstream surface (the module's shared locals).
- **The HTTP namespace has no update path** — its description change replaces it; the DNS namespaces update description in place.
- **The private namespace's `vpc` is never read back** (imports carry `{namespace_id}:{vpc_id}`).
- **A service binds its namespace in one of two places** — inside dns_config when it publishes records, at the top level when API-only (the provider's documented split for its legacy duplicated pointer); the module picks per service.
- **Instance registration is an AWS upsert** — create and update are the same RegisterInstance call; the module composes the attribute map from the typed spec fields plus the custom passthrough; deregistering an already-gone instance errors (no NotFound tolerance upstream).
- **`force_destroy` deregisters EVERY instance in the service** — including runtime-registered ones this manifest never declared.

Outputs mirror the Pulumi module key-for-key: `namespace_id` (import ID), `namespace_arn`, `hosted_zone_id`, `http_name`, `service_ids`, `service_arns`, `instance_service_ids`.
