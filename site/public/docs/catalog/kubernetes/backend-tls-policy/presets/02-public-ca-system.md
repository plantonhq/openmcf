---
title: "Public CA via System Trust Store"
description: "This preset secures the gateway-to-backend hop for backends serving PUBLICLY-issued certificates (Let's Encrypt, a commercial CA): instead of bringing a CA bundle, the policy trusts the gateway..."
type: "preset"
rank: "02"
presetSlug: "02-public-ca-system"
componentSlug: "backend-tls-policy"
componentTitle: "Backend TLS Policy"
provider: "kubernetes"
icon: "package"
order: 2
---

# Public CA via System Trust Store

This preset secures the gateway-to-backend hop for backends serving
PUBLICLY-issued certificates (Let's Encrypt, a commercial CA): instead of
bringing a CA bundle, the policy trusts the gateway implementation's
system certificate store via `well_known_ca_certificates: System` — the one
upstream-defined well-known set. There is nothing to create, mount, or
rotate on the trust side.

## When to Use

- The backend Service fronts an endpoint with a publicly-trusted
  certificate — a Service of type ExternalName toward a SaaS/upstream API,
  or an in-cluster backend deliberately serving a public certificate
- You want verified TLS to the backend without operating any private PKI

## Key Configuration Choices

- **`well_known_ca_certificates: System`** — mutually exclusive with
  `caCertificateRefs` (the spec enforces the CRD's own exactly-one-of
  rule); upstream support for this arm is Implementation-specific, so
  confirm your gateway implementation honors it
- **Literal `name.value` on the targetRef** — the pattern for backends
  that are not Planton-managed; switch the block to `valueFrom:` (kind
  `KubernetesService`, fieldPath `status.outputs.service_name`) to get the
  dependency edge when they are
- **`group: ""` present-but-empty** — Services live in the core API
  group; the key must be emitted for the CRD to accept the reference
- **`hostname`** — the SNI sent to the backend and the identity its
  publicly-issued certificate must prove; no `sectionName`, so the policy
  covers every port of the Service

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<backend-namespace>` | Namespace of the backend Service (the policy must live there too) | Your cluster namespaces |
| `<backend-service-name>` | Kubernetes Service name the policy attaches to | `kubectl get svc -n <backend-namespace>` |
| `api.example.com` | Public hostname on the backend's certificate (also sent as SNI) | The backend's certificate / the upstream API's hostname |

## Related Presets

- **01-internal-ca-configmap** — backends signed by your own internal CA
- **03-spiffe-mtls-backend** — mesh backends whose certificate identity is
  a SPIFFE URI rather than the SNI hostname
