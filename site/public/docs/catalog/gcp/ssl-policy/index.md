---
title: "SSL Policy"
description: "SSL Policy deployment documentation"
icon: "package"
order: 100
componentName: "gcpsslpolicy"
---

# GCP SSL Policy

Creates a Compute Engine SSL policy — the control for which TLS versions and cipher suites a load balancer accepts from clients. Reference its `self_link` from a target HTTPS (or SSL) proxy to enforce modern TLS; without one, GCP's permissive default applies (minimum TLS 1.0, COMPATIBLE ciphers).

## What Gets Created

A single SSL policy — global when `region` is empty, regional when set. Many proxies can share one policy, and tightening it later applies fleet-wide in place.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId` (or the provider's default project)
- **IAM permissions** — any role carrying `compute.sslPolicies.*` on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpSslPolicy
metadata:
  name: prod-tls-floor
spec:
  projectId:
    value: my-gcp-project-123
  profile: MODERN
  minTlsVersion: TLS_1_2
```

```shell
planton apply -f ssl-policy.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `profile` | `string` | `COMPATIBLE` | Cipher profile: `COMPATIBLE`, `MODERN`, `RESTRICTED`, or `CUSTOM`. Mutable. |
| `minTlsVersion` | `string` | `TLS_1_0` | Minimum TLS version: `TLS_1_0`, `TLS_1_1`, or `TLS_1_2`. Mutable. |
| `customFeatures` | `string[]` | `[]` | Exact cipher suites — required with (and only valid with) `CUSTOM`. Mutable. |
| `region` | `string` | `""` (global) | Region for a regional policy; empty means global. Immutable. |
| `sslPolicyName` | `string` | `metadata.name` | Cloud-side name (RFC1035). Immutable. |
| `description` | `string` | `""` | Why this policy exists. Immutable. |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the policy. Immutable. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `self_link` | Self-link URI — the value a target HTTPS (or SSL) proxy references in `ssl_policy` |
| `ssl_policy_name` | Name of the SSL policy in GCP |
| `enabled_features` | Cipher suites the policy actually enables, as computed by GCP |
| `region` | Region of a regional policy; empty for global |

## Related Components

- [GcpTargetHttpsProxy](/docs/catalog/gcp/target-https-proxy) — attaches this policy to harden client handshakes
- [GcpSslCertificate](/docs/catalog/gcp/ssl-certificate-self-managed) — self-managed certificate presented by the same proxy
- [GcpManagedSslCertificate](/docs/catalog/gcp/managed-ssl-certificate) — Google-managed certificate alternative
