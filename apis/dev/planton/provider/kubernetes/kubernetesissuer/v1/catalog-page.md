# KubernetesIssuer

Creates one cert-manager Issuer — a NAMESPACE-scoped certificate signing authority. Identical signing capabilities to KubernetesClusterIssuer (ACME, CA, self-signed, Vault); the namespace scope keeps a team's CA keypair and DNS credentials readable only inside the team's namespace.

## What Gets Created

- **Issuer** — named after the resource, in `spec.namespace`
- **Credential Secrets** — credentials declared in the spec, materialized in the Issuer's own namespace

## Prerequisites

- cert-manager on the cluster (**KubernetesCertManager**)

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesIssuer
metadata:
  name: team-a-selfsigned
spec:
  namespace:
    value: team-a
  config:
    selfSigned: {}
```

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | The Issuer's namespace |
| `issuer_name` | The issuer handle same-namespace Certificates reference |
| `acme_account_key_secret_name` | ACME account key Secret (empty for non-ACME) |

## Next Steps

The internal-PKI bootstrap: this self-signed Issuer signs a root CA **KubernetesCertificate** (`isCa: true`), whose Secret then powers a `ca`-backed Issuer for leaf certificates.
