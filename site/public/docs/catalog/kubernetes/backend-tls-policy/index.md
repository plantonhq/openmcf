---
title: "Backend TLS Policy"
description: "Backend TLS Policy deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesbackendtlspolicy"
---

# Kubernetes Backend TLS Policy

Provision a Kubernetes Gateway API `BackendTLSPolicy` -- the standard way to
tell a Gateway implementation to originate TLS to the backends BEHIND the
gateway and verify the certificate they present. Routes decide WHERE traffic
goes; this policy decides HOW the gateway-to-backend hop is secured. Use it
for end-to-end encryption: TLS terminates at the gateway and is
re-originated -- verified -- to the backend.

BackendTLSPolicy is a standard-channel resource served as
`gateway.networking.k8s.io/v1` (Gateway API v1.6.1; the `v1alpha3` version
is deprecated upstream and no longer served); the default standard-channel
CRD install includes it.

## What Gets Created

- A namespaced `gateway.networking.k8s.io/v1` `BackendTLSPolicy` custom
  resource.
- 1-16 target references to same-namespace Services (one is the safest
  portable shape), each optionally narrowed to a named Service port via
  `sectionName`.
- A validation block: the trust anchor (a CA-bundle ConfigMap XOR the
  system trust store), the SNI/authentication hostname, and optional
  Subject Alternative Names.

## Prerequisites

- Gateway API standard-channel CRDs installed (`KubernetesGatewayApiCrds`).
- A Gateway implementation that honors BackendTLSPolicy (support varies).
- The backend Service the policy targets, in the same namespace (upstream
  forbids cross-namespace targetRefs for this policy).
- For the bring-your-own-CA arm: a same-namespace ConfigMap carrying the
  PEM CA bundle in a key named `ca.crt`.

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesBackendTlsPolicy
metadata:
  name: my-backend-tls
spec:
  namespace:
    value: app-ns
  targetRefs:
    - group: "" # Services live in the core API group — empty but present
      kind: Service
      name:
        value: my-backend
  validation:
    wellKnownCACertificates: System
    hostname: api.example.com
```

```bash
planton apply -f backendtlspolicy.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `namespace` | reference | Namespace to create the policy in (the policy can only target Services there). |
| `targetRefs` | list | 1-16 same-namespace Service references; `group` present-but-empty for the core group; `name` is a reference (defaults to `KubernetesService`). |
| `validation` | object | Trust anchor (exactly one of `caCertificateRefs` / `wellKnownCACertificates`), mandatory `hostname` (SNI + certificate identity), optional `subjectAltNames`. |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `targetRefs[].sectionName` | string | Service port name -- narrows the policy to one port; omit to cover all ports. |
| `validation.subjectAltNames` | list | Up to 5 SANs (`Hostname` or `URI` -- SPIFFE IDs are the common URI case); with SANs set, `hostname` is SNI-only. |
| `options` | map | Up to 16 implementation-specific TLS options (domain-prefixed keys). |

## Examples

### Internal CA via ConfigMap

```yaml
spec:
  namespace:
    value: app-ns
  targetRefs:
    - group: ""
      kind: Service
      name:
        value: my-backend
  validation:
    caCertificateRefs:
      - group: ""
        kind: ConfigMap
        name:
          value: internal-ca-bundle
    hostname: backend.internal.example.com
```

### SPIFFE mTLS backend (port-scoped, URI SAN)

```yaml
spec:
  namespace:
    value: app-ns
  targetRefs:
    - group: ""
      kind: Service
      name:
        value: my-backend
      sectionName: https
  validation:
    caCertificateRefs:
      - group: ""
        kind: ConfigMap
        name:
          value: mesh-trust-bundle
    hostname: backend.internal.example.com
    subjectAltNames:
      - type: URI
        uri: spiffe://cluster.example.com/ns/payments/sa/backend
```

## Stack Outputs

| Output | Description |
|--------|-------------|
| `policy_name` | Name of the created BackendTLSPolicy (equals `metadata.name`). |
| `namespace` | Namespace the BackendTLSPolicy was created in. |

## Related Components

- [Kubernetes Gateway](kubernetesgateway)
- [Kubernetes HTTP Route](kuberneteshttproute)
- [Kubernetes Gateway API CRDs](kubernetesgatewayapicrds)
- [Kubernetes Service](kubernetesservice)
- [Kubernetes Config Map](kubernetesconfigmap)
- [Kubernetes Namespace](kubernetesnamespace)
