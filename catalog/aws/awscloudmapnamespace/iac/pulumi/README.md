# AwsCloudMapNamespace — Pulumi module

Manages one Cloud Map namespace (`servicediscovery.HttpNamespace` XOR `servicediscovery.PrivateDnsNamespace` XOR `servicediscovery.PublicDnsNamespace`, selected by the spec's type) with its services (`servicediscovery.Service`) and statically registered instances (`servicediscovery.Instance`).

Module facts worth knowing before editing:

- **The three types are three provider resources**; exactly one exists, all exposing the same downstream surface.
- **The HTTP namespace has no update path** — its description change replaces it; the DNS namespaces update in place.
- **The private namespace's Vpc is never read back** (imports carry `{namespace_id}:{vpc_id}`).
- **A service binds its namespace in one of two places** — DnsConfig.NamespaceId when it publishes records, the top-level NamespaceId when API-only; the module picks per service.
- **Instance registration is an AWS upsert** — the module composes the Attributes map from the typed spec fields (AWS_INSTANCE_IPV4/IPV6/PORT/CNAME, AWS_ALIAS_DNS_NAME, AWS_EC2_INSTANCE_ID) plus the custom passthrough; deregistering an already-gone instance errors upstream.
- **ForceDestroy deregisters EVERY instance in the service** — including runtime-registered ones.

Outputs mirror the Terraform module key-for-key: `namespace_id` (import ID), `namespace_arn`, `hosted_zone_id`, `http_name`, `service_ids`, `service_arns`, `instance_service_ids`.
